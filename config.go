package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"time"
)

type Config struct {
	AwsRegion string   `yaml:"aws_region"`
	Watcher   Watcher  `yaml:"watcher"`
	Finisher  Finisher `yaml:"finisher"`
}

type Watcher struct {
	AutoscalingGroupName string        `yaml:"autoscaling_group_name"`
	IntervalSec          time.Duration `yaml:"interval_sec"`
}

type Finisher struct {
	LifecycleHookName     string    `yaml:"lifecycle_hook_name"`
	LifecycleActionResult string    `yaml:"lifecycle_action_result"`
	Terminate             Terminate `yaml:"terminate"`
	Wait                  Wait      `yaml:"wait"`
}

type Terminate struct {
	Command string `yaml:"command"`
}

type Wait struct {
	Command     string        `yaml:"command"`
	IntervalSec time.Duration `yaml:"interval_sec"`
	MaxTries    int64         `yaml:"max_tries"`
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
