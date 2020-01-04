package shutter

import (
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/mattn/go-shellwords"
	"os/exec"
	"strings"
)

var parser = shellwords.NewParser()

func Render(commands string, instance *ec2.Instance) string {
	replaced := commands

	if instance.InstanceId != nil {
		replaced = strings.ReplaceAll(replaced, "${INSTANCE_ID}", *instance.InstanceId)
	}

	if instance.PublicIpAddress != nil {
		replaced = strings.ReplaceAll(replaced, "${PUBLIC_IP_ADDRESS}", *instance.PublicIpAddress)
	}

	if instance.PrivateIpAddress != nil {
		replaced = strings.ReplaceAll(replaced, "${PRIVATE_IP_ADDRESS}", *instance.PrivateIpAddress)
	}

	return replaced
}

func DoCommand(commands string) (int, error) {
	cmds, err := parser.Parse(commands)
	if err != nil {
		return 1, err
	}

	cmd := exec.Command(cmds[0], cmds[1:]...)
	err = cmd.Start()
	if err != nil {
		return 1, err
	}

	return cmd.ProcessState.ExitCode(), nil
}
