package main

import (
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"time"
)

type watcher struct {
	client *autoscalingClient
	config *Config
}

func NewWatcher(client *autoscalingClient, config *Config) *watcher {
	return &watcher{
		client: client,
		config: config,
	}
}

func (w *watcher) Watch() ([]autoscaling.Instance, error) {
	g, err := w.client.DescribeAutoscalingGroup(w.config.Watcher.AutoscalingGroupName)
	if err != nil {
		return nil, err
	}

	instances := []autoscaling.Instance{}
	for _, i := range g.Instances {
		if *i.LifecycleState == "Terminating" {
			instances = append(instances, *i)
		}
	}

	return instances, nil
}

func (w *watcher) Start(channel chan autoscaling.Instance) error {
	for {
		instances, err := w.Watch()
		if err != nil {
			return err
		}

		for _, i := range instances {
			channel <- i
		}

		time.Sleep(time.Second * 60)
	}
}
