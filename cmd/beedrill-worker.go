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
	configPath    = flag.String("c", "config_beedrill-worker.yml", "Path to a configuration file")
	broker        = flag.String("b", "redis://redis:6379/", "Broker URL")
	resultBackend = flag.String("r", "redis://redis:6379/", "Result backend")
	defaultQueue  = flag.String("q", "machinery_tasks", "Ephemeral Redis queue name")

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
		DefaultQueue:  *defaultQueue,
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
		"add":         exampletasks.Add,
		"multiply":    exampletasks.Multiply,
		"panic_task":  exampletasks.PanicTask,
		"simple_test": exampletasks.SimpleTest,
		"sleep":       exampletasks.RestfulSleep,
		"get_busy":    exampletasks.GetBusy,
		"TCP_socket":  exampletasks.OperateTCP,
		"task_args":   exampletasks.TaskArgs,
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
