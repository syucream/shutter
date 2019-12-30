package main

import (
	"github.com/mattn/go-shellwords"
	"os/exec"
)

var parser = shellwords.NewParser()

func DoCommand(commands string) error {
	cmds, err := parser.Parse(commands)
	if err != nil {
		return err
	}

	return exec.Command(cmds[0], cmds[1:]...).Start()
}
