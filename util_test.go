package shutter

import (
	"github.com/aws/aws-sdk-go/service/ec2"
	"testing"
)

func TestRender(t *testing.T) {
	instanceId := "i-1234abcd"
	publicIpAddr := "8.8.8.8"
	privateIpAddr := "192.168.0.1"

	instance := &ec2.Instance{
		InstanceId:       &instanceId,
		PublicIpAddress:  &publicIpAddr,
		PrivateIpAddress: &privateIpAddr,
	}
	source := "${INSTANCE_ID} ; ${PUBLIC_IP_ADDRESS} ; ${PRIVATE_IP_ADDRESS}"

	actual := Render(source, instance)
	expected := "i-1234abcd ; 8.8.8.8 ; 192.168.0.1"

	if actual != expected {
		t.Fatal("actual doesn't match expected", actual)
	}
}

func TestDoCommand(t *testing.T) {
	validCommand := "echo hoge"
	invalidCommand := "cat /path/to/noexistent.txt"

	_, err := DoCommand(validCommand)
	if err != nil {
		t.Fatal("it should not fail", err)
	}

	s, _ := DoCommand(invalidCommand)
	if s.Success() {
		t.Fatal("it should fail", s)
	}
}
