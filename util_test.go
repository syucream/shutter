package main

import (
	"github.com/aws/aws-sdk-go/service/ec2"
	"testing"
)

func TestRender(t *testing.T) {
	publicIpAddr := "8.8.8.8"
	privateIpAddr := "192.168.0.1"

	instance := &ec2.Instance{
		PublicIpAddress:  &publicIpAddr,
		PrivateIpAddress: &privateIpAddr,
	}
	source := "${PUBLIC_IP_ADDRESS} ; ${PRIVATE_IP_ADDRESS}"

	actual := Render(source, instance)
	expected := "8.8.8.8 ; 192.168.0.1"

	if actual != expected {
		t.Fatal("actual doesn't match expected", actual)
	}
}
