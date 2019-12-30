package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	Watcher  Watcher  `yaml:"watcher"`
	Finisher Finisher `yaml:"finisher"`
}

type Watcher struct {
	AutoscalingGroupName string `yaml:"autoscaling_group_name"`
}

type Finisher struct {
	LifecycleHookName      string `yaml:"lifecycle_hook_name"`
	StartCompletionCommand string `yaml:"start_command"`
	WaitCompletionCommand  string `yaml:"wait_command"`
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
