package beedrilltasks

import (
	"io/ioutil"
	"os/exec"
	"strings"
)

// Command passed to workers
func TaskArgs(args ...string) (string, error) {
	splitted_args := strings.Split(strings.Join(args, " "), " ")

	res, err := exec.Command(splitted_args[0], splitted_args[1:]...).Output()

	if err != nil {
		return "Error when executing " + splitted_args[0], err
	}

	return string(res), nil
}

func TaskFile(cmd string, contents []byte) (string, error) {
	err := ioutil.WriteFile("/root/readfiles.f", contents, 0644)
	if err != nil {
		return "Error when writing readfiles.f", err
	}

	return TaskArgs(cmd, "-f readfiles.f")
}
