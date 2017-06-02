package main

import (
	"fmt"

	"git.rnd.alterway.fr/Wolphin-project/beedrill/pkg"

	machinery "github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/config"
	//"github.com/RichardKnop/machinery/v1/log"
	//"github.com/RichardKnop/machinery/v1/tasks"
	"github.com/urfave/cli"
)

// Define flags
var (
	app           *cli.App
	configPath    string
	broker        string
	resultBackend string
	defaultQueue  string
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
			Value:       "config_beedrill-worker.yml",
			Destination: &configPath,
			Usage:       "Path to a configuration file",
		},
		cli.StringFlag{
			Name:        "b",
			Value:       "redis://redis:6379/",
			Destination: &broker,
			Usage:       "Broker URL",
		},
		cli.StringFlag{
			Name:        "r",
			Value:       "redis://redis:6379/",
			Destination: &resultBackend,
			Usage:       "Result backend",
		},
		cli.StringFlag{
			Name:        "q",
			Value:       "machinery_tasks",
			Destination: &defaultQueue,
			Usage:       "Ephemeral Redis queue name",
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

func main() {
	server, err := startServer()
	if err != nil {
		fmt.Errorf("Could not start the worker server", err.Error())
	}

	worker := server.NewWorker("beedrill_worker")

	if err := worker.Launch(); err != nil {
		fmt.Errorf("Could not launch the worker", err.Error())
	}
}
