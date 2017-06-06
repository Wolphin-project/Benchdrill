package beedrilltasks

import (
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
