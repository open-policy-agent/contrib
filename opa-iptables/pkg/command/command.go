package command

import (
	"fmt"
	"io/ioutil"
	"os/exec"
)

func RunCommand(name string, args ...string) (output []byte, err error) {
	cmd := exec.Command(name, args...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}

	err = cmd.Start()
	if err != nil {
		return nil, err
	}

	stderrorCmd, err := ioutil.ReadAll(stderr)
	if err != nil {
		return nil, err
	}

	stdoutCmd, err := ioutil.ReadAll(stdout)
	if err != nil {
		return nil, err
	}

	if err := cmd.Wait(); err != nil {
		return nil, err
	}

	if string(stderrorCmd) != "" {
		return nil, fmt.Errorf(string(stderrorCmd))
	}

	return stdoutCmd, nil
}