package main

import (
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/mattn/go-shellwords"
	"os/exec"
	"strings"
)

var parser = shellwords.NewParser()

func Render(commands string, instance *ec2.Instance) string {
	replaced := commands

	replaced = strings.ReplaceAll(replaced, "${PUBLIC_IP_ADDRESS}", *instance.PublicIpAddress)
	replaced = strings.ReplaceAll(replaced, "${PRIVATE_IP_ADDRESS}", *instance.PrivateIpAddress)

	return replaced
}

func DoCommand(commands string) error {
	cmds, err := parser.Parse(commands)
	if err != nil {
		return err
	}

	return exec.Command(cmds[0], cmds[1:]...).Start()
}
