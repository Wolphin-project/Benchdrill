package integrationtests

import (
	"os"
	"testing"

	"github.com/RichardKnop/machinery/v1/config"
)

func TestAmqpAmqp(t *testing.T) {
	amqpURL := os.Getenv("AMQP_URL")
	if amqpURL == "" {
		return
	}

	// AMQP broker, AMQP result backend
	server := setup(&config.Config{
		Broker:        amqpURL,
		DefaultQueue:  "test_queue",
		ResultBackend: amqpURL,
		AMQP: &config.AMQPConfig{
			Exchange:      "test_exchange",
			ExchangeType:  "direct",
			BindingKey:    "test_task",
			PrefetchCount: 1,
		},
	})
	worker := server.NewWorker("test_worker")
	go worker.Launch()
	testAll(server, t)
	worker.Quit()
}
