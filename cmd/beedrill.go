package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"git.rnd.alterway.fr/Wolphin-project/beedrill/pkg"

	"github.com/RichardKnop/machinery/v1"
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
	times         int
)

func init() {
	app = cli.NewApp()

	app.Name = "Beedrill"
	app.Usage = "Basic benchmark tool"
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
			Value:       "beedrill_tasks",
			Destination: &defaultQueue,
			Usage:       "Ephemeral Redis queue name",
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

	// If present, the config file takes priority over CLI flags
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
		"task_file": beedrilltasks.TaskFile,
	}

	return server, server.RegisterTasks(tasks)
}

func sendCmdArgs(cmd string) error {
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
					Value: cmd,
				},
			},
		}
	}

	initTasks()

	var s []*tasks.Signature

	for i := 0; i < times; i++ {
		s = append(s, &task0)
	}

	groupedTasks := tasks.NewGroup(s...)

	asyncResults, err := server.SendGroup(groupedTasks)
	if err != nil {
		return fmt.Errorf("Could not send task: %s", err.Error())
	}

	log.INFO.Println("Command passed to worker…")

	for _, asyncResult := range asyncResults {
		results, err := asyncResult.Get(time.Duration(time.Millisecond * 5))
		if err != nil {
			return fmt.Errorf("Getting task result failed with error: %s", err.Error())
		}

		for _, result := range results {
			log.INFO.Printf("%v\n", result.Interface())
		}
	}

	return nil
}

func sendCmdFile(cmd, file string) error {
	server, err := startServer()
	if err != nil {
		return err
	}

	var task0 tasks.Signature

	var initTasks = func() {
		task0 = tasks.Signature{
			Name: "task_file",
			Args: []tasks.Arg{
				{
					Type:  "string",
					Value: cmd,
				},
				{
					Type:  "string",
					Value: file,
				},
			},
		}
	}

	initTasks()

	var s []*tasks.Signature

	for i := 0; i < times; i++ {
		s = append(s, &task0)
	}

	groupedTasks := tasks.NewGroup(s...)

	asyncResults, err := server.SendGroup(groupedTasks)
	if err != nil {
		return fmt.Errorf("Could not send task: %s", err.Error())
	}

	log.INFO.Println("Command passed to worker…")

	for _, asyncResult := range asyncResults {
		results, err := asyncResult.Get(time.Duration(time.Millisecond * 5))
		if err != nil {
			return fmt.Errorf("Getting task result failed with error: %s", err.Error())
		}

		for _, result := range results {
			log.INFO.Printf("%v\n", result.Interface())
		}
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

func getStatus(queue string) error {
	server, err := startServer()

	if err != nil {
		return err
	}

	waitingTasks, err := server.GetBroker().GetPendingTasks(queue)

	if err != nil {
		return err
	}

	for _, waitingTask := range waitingTasks {
		log.INFO.Printf("%v", waitingTask.Name)
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
			Name:  "send_cmd_args",
			Usage: "Send command with arguments",
			Action: func(c *cli.Context) error {
				return sendCmdArgs(c.Args().First())
			},
		},
		{
			Name:  "send_cmd_file",
			Usage: "Send command with file",
			Action: func(c *cli.Context) error {
				file, err := ioutil.ReadFile("/dev/stdin")

				if err != nil {
					return fmt.Errorf("%s", err.Error())
				}

				return sendCmdFile(c.Args().First(), string(file))
			},
		},
		{
			Name:  "get_status",
			Usage: "Get tasks which are waiting in the queue",
			Action: func(c *cli.Context) error {
				return getStatus(c.Args().First())
			},
		},
	}

	// Run the CLI app
	if err := app.Run(os.Args); err != nil {
		log.FATAL.Print(err)
	}
}
