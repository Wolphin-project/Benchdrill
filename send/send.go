package main

import (
	"flag"
	"fmt"

	machinery "github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/config"
	"github.com/RichardKnop/machinery/v1/errors"
	"github.com/RichardKnop/machinery/v1/signatures"
)

// Define flags
var (
	configPath    = flag.String("c", "config.yml", "Path to a configuration file")
	broker        = flag.String("b", "redis://127.0.0.1:6379/", "Broker URL")
	resultBackend = flag.String("r", "redis://127.0.0.1:6379/", "Result backend")
	// resultBackend = flag.String("r", "redis://127.0.0.1:6379", "Result backend")
	// resultBackend = flag.String("r", "memcache://127.0.0.1:11211", "Result backend")
	// resultBackend = flag.String("r", "mongodb://127.0.0.1:27017", "Result backend")
	exchange     = flag.String("e", "machinery_exchange", "Durable, non-auto-deleted AMQP exchange name")
	exchangeType = flag.String("t", "direct", "Exchange type - direct|fanout|topic|x-custom")
	defaultQueue = flag.String("q", "machinery_tasks", "Ephemeral AMQP queue name")
	bindingKey   = flag.String("k", "machinery_task", "AMQP binding key")

	cnf                                             config.Config
	server                                          *machinery.Server
	task0                                           signatures.TaskSignature
)

func init() {
	// Parse the flags
	flag.Parse()

	cnf = config.Config{
		Broker:        *broker,
		ResultBackend: *resultBackend,
		Exchange:      *exchange,
		ExchangeType:  *exchangeType,
		DefaultQueue:  *defaultQueue,
		BindingKey:    *bindingKey,
	}

	// Parse the config
	// NOTE: If a config file is present, it has priority over flags
	data, err := config.ReadFromFile(*configPath)
	if err == nil {
		err = config.ParseYAMLConfig(&data, &cnf)
		errors.Fail(err, "Could not parse config file")
	}

	server, err = machinery.NewServer(&cnf)
	errors.Fail(err, "Could not initialize server")
}

func initTasks() {
	task0 = signatures.TaskSignature{
		Name: "simple_test",
	}
}

func main() {
	/*
	 * First, let's try sending a single task
	 */
	initTasks()
	fmt.Println("Single simple task:")

	asyncResult, err := server.SendTask(&task0)
	errors.Fail(err, "Could not send task")

	result, err := asyncResult.Get()
	errors.Fail(err, "Getting task state failed with error")
	fmt.Printf("%v\n", result.Interface())
}
