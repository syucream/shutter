package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	TerminateCommand      string `yaml:"terminate_command"`
	WaitCompletionCommand string `yaml:"wait_completion_command"`
}

func NewConfig(name string) (*Config, error) {
	d, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, err
	}

	c := &Config{}
	err = yaml.Unmarshal(d, &c)
	if err != nil {
		return nil, err
	}

	return c, nil
}
