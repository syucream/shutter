package main

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestNewConfig(t *testing.T) {
	validFile, err := ioutil.TempFile("", "validConfig.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Remove(validFile.Name()) }()
	defer func() { _ = validFile.Close() }()

	_, err = validFile.Write([]byte(`
aws_region: us-east-1
watcher:
  autoscaling_group_name: test
  interval_sec: 60
finisher:
  lifecycle_hook_name: test
  lifecycle_action_result: ABANDON
  terminate:
    command: "ping -c 5 ${PUBLIC_IP_ADDRESS}"
  wait:
    command: "ping -c 5 ${PUBLIC_IP_ADDRESS}"
    interval_sec: 60
    max_tries: 30
`))
	if err != nil {
		t.Fatal(err)
	}

	invalidFile, err := ioutil.TempFile("", "invalidConfig.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Remove(invalidFile.Name()) }()
	defer func() { _ = invalidFile.Close() }()

	_, err = invalidFile.Write([]byte(`
aws_region: us-east-1
watcher: hoge
finisher: fuga
`))
	if err != nil {
		t.Fatal(err)
	}

	_, err = NewConfig(validFile.Name())
	if err != nil {
		t.Fatal("err should not occured", err)
	}

	_, err = NewConfig(invalidFile.Name())
	if err == nil {
		t.Fatal("err should occured")
	}
}
