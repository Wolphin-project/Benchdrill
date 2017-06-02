package brokers

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/RichardKnop/machinery/v1/config"
	"github.com/RichardKnop/machinery/v1/log"
	"github.com/RichardKnop/machinery/v1/retry"
	"github.com/RichardKnop/machinery/v1/tasks"
	"github.com/streadway/amqp"
)

// AMQPBroker represents an AMQP broker
type AMQPBroker struct {
	Broker
}

// NewAMQPBroker creates new AMQPBroker instance
func NewAMQPBroker(cnf *config.Config) Interface {
	return &AMQPBroker{Broker: New(cnf)}
}

// StartConsuming enters a loop and waits for incoming messages
func (b *AMQPBroker) StartConsuming(consumerTag string, taskProcessor TaskProcessor) (bool, error) {
	b.startConsuming(consumerTag, taskProcessor)

	conn, channel, queue, _, err := b.connect()
	if err != nil {
		b.retryFunc()
		return b.retry, err
	}
	defer b.close(channel, conn)

	b.retryFunc = retry.Closure()

	if err = channel.Qos(
		b.cnf.AMQP.PrefetchCount,
		0,     // prefetch size
		false, // global
	); err != nil {
		return b.retry, fmt.Errorf("Channel qos error: %s", err)
	}

	deliveries, err := channel.Consume(
		queue.Name,  // queue
		consumerTag, // consumer tag
		false,       // auto-ack
		false,       // exclusive
		false,       // no-local
		false,       // no-wait
		nil,         // arguments
	)
	if err != nil {
		return b.retry, fmt.Errorf("Queue consume error: %s", err)
	}

	log.INFO.Print("[*] Waiting for messages. To exit press CTRL+C")

	if err := b.consume(deliveries, taskProcessor); err != nil {
		return b.retry, err
	}

	return b.retry, nil
}

// StopConsuming quits the loop
func (b *AMQPBroker) StopConsuming() {
	b.stopConsuming()
}

// Publish places a new message on the default queue
func (b *AMQPBroker) Publish(signature *tasks.Signature) error {
	b.AdjustRoutingKey(signature)

	// Check the ETA signature field, if it is set and it is in the future,
	// delay the task
	if signature.ETA != nil {
		now := time.Now().UTC()

		if signature.ETA.After(now) {
			delayMs := int64(signature.ETA.Sub(now) / time.Millisecond)

			return b.delay(signature, delayMs)
		}
	}

	message, err := json.Marshal(signature)
	if err != nil {
		return fmt.Errorf("JSON marshal error: %v", err)
	}

	conn, channel, _, confirmsChan, err := b.connect()
	if err != nil {
		return err
	}
	defer b.close(channel, conn)

	if err := channel.Publish(
		b.cnf.AMQP.Exchange,  // exchange
		signature.RoutingKey, // routing key
		false,                // mandatory
		false,                // immediate
		amqp.Publishing{
			Headers:      amqp.Table(signature.Headers),
			ContentType:  "application/json",
			Body:         message,
			DeliveryMode: amqp.Persistent,
		},
	); err != nil {
		return err
	}

	confirmed := <-confirmsChan

	if confirmed.Ack {
		return nil
	}

	return fmt.Errorf("Failed delivery of delivery tag: %v", confirmed.DeliveryTag)
}

// consume takes delivered messages from the channel and manages a worker pool
// to process tasks concurrently
func (b *AMQPBroker) consume(deliveries <-chan amqp.Delivery, taskProcessor TaskProcessor) error {
	maxWorkers := b.cnf.MaxWorkerInstances
	pool := make(chan struct{}, maxWorkers)

	// initialize worker pool with maxWorkers workers
	go func() {
		for i := 0; i < maxWorkers; i++ {
			pool <- struct{}{}
		}
	}()

	errorsChan := make(chan error)

	// Use wait group to make sure task processing completes on interrupt signal
	var wg sync.WaitGroup
	defer wg.Wait()

	for {
		select {
		case err := <-errorsChan:
			return err
		case d := <-deliveries:
			if maxWorkers != 0 {
				// get worker from pool (blocks until one is available)
				<-pool
			}

			wg.Add(1)

			// Consume the task inside a gotourine so multiple tasks
			// can be processed concurrently
			go func() {
				defer wg.Done()

				if err := b.consumeOne(d, taskProcessor); err != nil {
					errorsChan <- err
				}

				if maxWorkers != 0 {
					// give worker back to pool
					pool <- struct{}{}
				}
			}()
		case <-b.stopChan:
			return nil
		}
	}
}

// consumeOne processes a single message using TaskProcessor
func (b *AMQPBroker) consumeOne(d amqp.Delivery, taskProcessor TaskProcessor) error {
	if len(d.Body) == 0 {
		d.Nack(false, false)                           // multiple, requeue
		return errors.New("Received an empty message") // RabbitMQ down?
	}

	log.INFO.Printf("Received new message: %s", d.Body)

	// Unmarshal message body into signature struct
	signature := new(tasks.Signature)
	if err := json.Unmarshal(d.Body, signature); err != nil {
		d.Nack(false, false) // multiple, requeue
		return err
	}

	// If the task is not registered, we nack it and requeue,
	// there might be different workers for processing specific tasks
	if !b.IsTaskRegistered(signature.Name) {
		d.Nack(false, true) // multiple, requeue
		return nil
	}

	d.Ack(false) // multiple
	return taskProcessor.Process(signature)
}

// delay a task by delayDuration miliseconds, the way it works is a new queue
// is created without any consumers, the message is then published to this queue
// with appropriate ttl expiration headers, after the expiration, it is sent to
// the proper queue with consumers
func (b *AMQPBroker) delay(signature *tasks.Signature, delayMs int64) error {
	if delayMs <= 0 {
		return errors.New("Cannot delay task by 0ms")
	}

	var (
		conn    *amqp.Connection
		channel *amqp.Channel
		queue   amqp.Queue
		err     error
	)

	message, err := json.Marshal(signature)
	if err != nil {
		return fmt.Errorf("JSON marshal error: %v", err)
	}

	// Connect to server
	conn, channel, err = b.open()
	if err != nil {
		return err
	}
	defer b.close(channel, conn)

	// Declare an exchange
	if err = channel.ExchangeDeclare(
		b.cnf.AMQP.Exchange,     // name of the exchange
		b.cnf.AMQP.ExchangeType, // type
		true,  // durable
		false, // delete when complete
		false, // internal
		false, // noWait
		nil,   // arguments
	); err != nil {
		return fmt.Errorf("Exchange declare error: %s", err)
	}

	// It's necessary to redeclare the queue each time (to zero its TTL timer).
	holdQueue := fmt.Sprintf(
		"delay.%d.%s.%s",
		delayMs, // delay duration in mileseconds
		b.cnf.AMQP.Exchange,
		b.cnf.AMQP.BindingKey, // routing key
	)
	holdQueueArgs := amqp.Table{
		// Exchange where to send messages after TTL expiration.
		"x-dead-letter-exchange": b.cnf.AMQP.Exchange,
		// Routing key which use when resending expired messages.
		"x-dead-letter-routing-key": b.cnf.AMQP.BindingKey,
		// Time in milliseconds
		// after that message will expire and be sent to destination.
		"x-message-ttl": delayMs,
		// Time after that the queue will be deleted.
		"x-expires": delayMs * 2,
	}
	queue, err = channel.QueueDeclare(
		holdQueue,     // name
		false,         // durable
		true,          // delete when unused
		false,         // exclusive
		false,         // no-wait
		holdQueueArgs, // arguments
	)
	if err != nil {
		return fmt.Errorf("Queue declare error: %s", err)
	}

	// Bind the queue
	if err := channel.QueueBind(
		queue.Name,          // name of the queue
		queue.Name,          // binding key
		b.cnf.AMQP.Exchange, // source exchange
		false,               // noWait
		amqp.Table(b.cnf.AMQP.QueueBindingArguments), // arguments
	); err != nil {
		return fmt.Errorf("Queue bind error: %s", err)
	}

	if err := channel.Publish(
		b.cnf.AMQP.Exchange, // exchange
		holdQueue,           // routing key
		false,               // mandatory
		false,               // immediate
		amqp.Publishing{
			Headers:      amqp.Table(signature.Headers),
			ContentType:  "application/json",
			Body:         message,
			DeliveryMode: amqp.Persistent,
		},
	); err != nil {
		return err
	}

	return nil
}

// connect opens a connection to RabbitMQ, declares an exchange, opens a channel,
// declares and binds the queue and enables publish notifications
func (b *AMQPBroker) connect() (*amqp.Connection, *amqp.Channel, amqp.Queue, <-chan amqp.Confirmation, error) {
	var (
		conn    *amqp.Connection
		channel *amqp.Channel
		queue   amqp.Queue
		err     error
	)

	// Connect to server
	conn, channel, err = b.open()
	if err != nil {
		return conn, channel, queue, nil, err
	}

	// Declare an exchange
	if err = channel.ExchangeDeclare(
		b.cnf.AMQP.Exchange,     // name of the exchange
		b.cnf.AMQP.ExchangeType, // type
		true,  // durable
		false, // delete when complete
		false, // internal
		false, // noWait
		nil,   // arguments
	); err != nil {
		return conn, channel, queue, nil, fmt.Errorf("Exchange declare error: %s", err)
	}

	// Declare a queue
	queue, err = channel.QueueDeclare(
		b.cnf.DefaultQueue, // name
		true,               // durable
		false,              // delete when unused
		false,              // exclusive
		false,              // no-wait
		nil,                // arguments
	)
	if err != nil {
		return conn, channel, queue, nil, fmt.Errorf("Queue declare error: %s", err)
	}

	// Bind the queue
	if err := channel.QueueBind(
		queue.Name,            // name of the queue
		b.cnf.AMQP.BindingKey, // binding key
		b.cnf.AMQP.Exchange,   // source exchange
		false,                 // noWait
		amqp.Table(b.cnf.AMQP.QueueBindingArguments), // arguments
	); err != nil {
		return conn, channel, queue, nil, fmt.Errorf("Queue bind error: %s", err)
	}

	// Enable publish confirmations
	if err := channel.Confirm(false); err != nil {
		return conn, channel, queue, nil, fmt.Errorf("Channel could not be put into confirm mode: %s", err)
	}

	return conn, channel, queue, channel.NotifyPublish(make(chan amqp.Confirmation, 1)), nil
}

// open new RabbitMQ connection
func (b *AMQPBroker) open() (*amqp.Connection, *amqp.Channel, error) {
	var (
		conn    *amqp.Connection
		channel *amqp.Channel
		err     error
	)

	// Connect
	// From amqp docs: DialTLS will use the provided tls.Config when it encounters an amqps:// scheme
	// and will dial a plain connection when it encounters an amqp:// scheme.
	conn, err = amqp.DialTLS(b.cnf.Broker, b.cnf.TLSConfig)
	if err != nil {
		return conn, channel, fmt.Errorf("Dial error: %s", err)
	}

	// Open a channel
	channel, err = conn.Channel()
	if err != nil {
		return conn, channel, fmt.Errorf("Open channel error: %s", err)
	}

	return conn, channel, nil
}

// close connection
func (b *AMQPBroker) close(channel *amqp.Channel, conn *amqp.Connection) error {
	if channel != nil {
		if err := channel.Close(); err != nil {
			return fmt.Errorf("Close channel error: %s", err)
		}
	}

	if conn != nil {
		if err := conn.Close(); err != nil {
			return fmt.Errorf("Close connection error: %s", err)
		}
	}

	return nil
}
