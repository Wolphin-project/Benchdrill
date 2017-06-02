package main

import (
	"fmt"
	"os"
	"time"

	"git.rnd.alterway.fr/Wolphin-project/beedrill/pkg"

	machinery "github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/config"
	"github.com/RichardKnop/machinery/v1/log"
	"github.com/RichardKnop/machinery/v1/tasks"
	"github.com/urfave/cli"
)

// Define flags
var (
	app           *cli.App
	configPath    string
	broker        string
	resultBackend string
	defaultQueue  string
	task          string
	arguments     string
	times         int
)

func init() {
	app = cli.NewApp()

	app.Name = "Beedrill"
	app.Usage = "sysbench tasks sent with machinery to workers"
	app.Author = "Alter Way"
	app.Email = "qqch@alterway.fr"
	app.Version = "0.1.0"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "c",
			Value:       "config_beedrill.yml",
			Destination: &configPath,
			Usage:       "Path to a configuration file",
		},
		cli.StringFlag{
			Name:        "b",
			Value:       "redis://127.0.0.1:6379/",
			Destination: &broker,
			Usage:       "Broker URL",
		},
		cli.StringFlag{
			Name:        "r",
			Value:       "redis://127.0.0.1:6379/",
			Destination: &resultBackend,
			Usage:       "Result backend",
		},
		cli.StringFlag{
			Name:        "q",
			Value:       "machinery_tasks",
			Destination: &defaultQueue,
			Usage:       "Ephemeral Redis queue name",
		},
		cli.StringFlag{
			Name:        "task",
			Value:       "",
			Destination: &task,
			Usage:       "Command to be executed by workers",
		},
		cli.StringFlag{
			Name:        "arguments",
			Value:       "",
			Destination: &arguments,
			Usage:       "Arguments to be passed with the task’s flag command to workers",
		},
		cli.IntFlag{
			Name:        "times",
			Value:       1,
			Destination: &times,
			Usage:       "Number of times tasks are sent",
		},
	}
}

func startServer() (*machinery.Server, error) {
	cnf := config.Config{
		Broker:        broker,
		DefaultQueue:  defaultQueue,
		ResultBackend: resultBackend,
	}

	// If present, the config file takes priority over cli flags
	data, err := config.ReadFromFile(configPath)
	if err != nil {
		log.WARNING.Printf("Could not load config from file: %s", err.Error())
	} else {
		if err = config.ParseYAMLConfig(&data, &cnf); err != nil {
			return nil, fmt.Errorf("Could not parse config file: %s", err.Error())
		}
	}

	// Create server instance
	server, err := machinery.NewServer(&cnf)
	if err != nil {
		return nil, fmt.Errorf("Could not initialize server: %s", err.Error())
	}

	// Register tasks
	tasks := map[string]interface{}{
		"task_args": beedrilltasks.TaskArgs,
	}

	server.RegisterTasks(tasks)

	return server, nil
}

func send() error {
	server, err := startServer()
	if err != nil {
		return err
	}

	var task0 tasks.Signature

	var initTasks = func() {
		task0 = tasks.Signature{
			Name: "task_args",
			Args: []tasks.Arg{
				{
					Type:  "string",
					Value: task,
				},
				{
					Type:  "string",
					Value: arguments,
				},
			},
		}
	}

	initTasks()

	log.INFO.Println("Command passed to worker…")

	for i := 0; i < times; i++ {
		asyncResult, err := server.SendTask(&task0)
		if err != nil {
			return fmt.Errorf("Could not send task: %s", err.Error())
		}

		results, err := asyncResult.Get(time.Duration(time.Millisecond * 5))
		if err != nil {
			return fmt.Errorf("Getting task result failed with error: %s", err.Error())
		}

		log.INFO.Printf("%v\n", results[0].Interface())
	}

	return nil
}

func worker() error {
	server, err := startServer()
	if err != nil {
		return err
	}

	// The second argument is a consumer tag
	// Ideally, each worker should have a unique tag (worker1, worker2 etc)
	worker := server.NewWorker("machinery_worker")

	if err := worker.Launch(); err != nil {
		return err
	}

	return nil
}

func main() {
	// Set the CLI app commands
	app.Commands = []cli.Command{
		{
			Name:  "worker",
			Usage: "Launch Beedrill worker",
			Action: func(c *cli.Context) error {
				return worker()
			},
		},
		{
			Name:  "Beedrill",
			Usage: "Send Beedrill tasks",
			Action: func(c *cli.Context) error {
				return send()
			},
		},
	}

	// Run the CLI app
	if err := app.Run(os.Args); err != nil {
		log.FATAL.Print(err)
	}
}
