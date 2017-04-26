package main

import (
	"flag"

	"git.rnd.alterway.fr/beedrill/pkg"
	machinery "github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/config"
	"github.com/RichardKnop/machinery/v1/errors"
)

// Define flags
var (
	configPath    = flag.String("c", "exampleconfig.yml", "Path to a configuration file")
	broker        = flag.String("b", "redis://redis:6379/", "Broker URL")
	resultBackend = flag.String("r", "redis://redis:6379/", "Result backend")
	// resultBackend = flag.String("r", "redis://127.0.0.1:6379", "Result backend")
	// resultBackend = flag.String("r", "memcache://127.0.0.1:11211", "Result backend")
	// resultBackend = flag.String("r", "mongodb://127.0.0.1:27017", "Result backend")
	exchange     = flag.String("e", "machinery_exchange", "Durable, non-auto-deleted AMQP exchange name")
	exchangeType = flag.String("t", "direct", "Exchange type - direct|fanout|topic|x-custom")
	defaultQueue = flag.String("q", "machinery_tasks", "Ephemeral AMQP queue name")
	bindingKey   = flag.String("k", "machinery_task", "AMQP binding key")

	cnf    config.Config
	server *machinery.Server
	worker *machinery.Worker
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

	// Register tasks
	tasks := map[string]interface{}{
		"add":           exampletasks.Add,
		"multiply":      exampletasks.Multiply,
		"panic_task":    exampletasks.PanicTask,
		"simple_test":   exampletasks.SimpleTest,
		"sleep":         exampletasks.RestfulSleep,
		"get_busy":      exampletasks.GetBusy,
		"TCP_socket":    exampletasks.OperateTCP,
		"sysbench_task": exampletasks.SysbenchTask,
	}
	server.RegisterTasks(tasks)

	// The second argument is a consumer tag
	// Ideally, each worker should have a unique tag (worker1, worker2 etc)
	worker = server.NewWorker("machinery_worker")
}

func main() {
	err := worker.Launch()
	errors.Fail(err, "Could not launch worker")
}
