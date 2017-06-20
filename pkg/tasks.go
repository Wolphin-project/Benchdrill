package beedrilltasks

import (
	"io/ioutil"
	"os/exec"
	"strings"
)

// Command passed to workers
func TaskArgs(cmd string) (string, error) {
	splitted_args := strings.Split(cmd, " ")

	res, err := exec.Command(splitted_args[0], splitted_args[1:]...).Output()

	if err != nil {
		return "Error when executing " + splitted_args[0], err
	}

	return string(res), nil
}

func TaskFile(cmd, file string) (string, error) {
	if err := ioutil.WriteFile("/root/workload.f", []byte(file), 0644); err != nil {
		return "Error when writing workload.f", err
	}

	return TaskArgs(cmd + "/root/workload.f")
}
