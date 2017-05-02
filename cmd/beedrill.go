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
	configPath    = flag.String("c", "config_beedrill.yml", "Path to a configuration file")
	broker        = flag.String("b", "redis://127.0.0.1:6379/", "Broker URL")
	resultBackend = flag.String("r", "redis://127.0.0.1:6379/", "Result backend")
	defaultQueue  = flag.String("q", "machinery_tasks", "Ephemeral Redis queue name")
	times         = flag.Int("times", 1, "Number of times tasks are sent")
	sysbench      = flag.String("sysbench", "--help", "Sysbench benchmark")

	cnf    config.Config
	server *machinery.Server
	task0  signatures.TaskSignature
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
}

func initTasks() {
	task0 = signatures.TaskSignature{
		Name: "sysbench_task",
		Args: []signatures.TaskArg{
			{
				Type:  "string",
				Value: sysbench,
			},
		},
	}

	/*task1 = signatures.TaskSignature{
		Name: "sleep",
	}

	task2 = signatures.TaskSignature{
		Name: "sleep",
	}

	task3 = signatures.TaskSignature{
		Name: "get_busy",
	}*/
}

func main() {
	/*
	 * First, let's try sending a single task
	 */
	initTasks()
	fmt.Println("Single simple task:")

	for i := 0; i < *times; i++ {
		asyncResult, err := server.SendTask(&task0)
		errors.Fail(err, "Could not send task")

		result, err := asyncResult.Get()
		errors.Fail(err, "Getting task state failed with error")
		fmt.Printf("%v\n", result.Interface())
	}
}
