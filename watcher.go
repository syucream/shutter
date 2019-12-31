package main

import (
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"go.uber.org/zap"
	"time"
)

const terminatingLifecycleState = "Terminating:Wait"

type watcher struct {
	client *awsClient
	config *Config
	logger *zap.Logger
}

func NewWatcher(client *awsClient, config *Config, logger *zap.Logger) *watcher {
	return &watcher{
		client: client,
		config: config,
		logger: logger,
	}
}

func (w *watcher) Watch() ([]autoscaling.Instance, error) {
	g, err := w.client.DescribeAutoscalingGroup(w.config.Watcher.AutoscalingGroupName)
	if err != nil {
		return nil, err
	}
	w.logger.Info("describe instances under the autoscaling group", zap.Reflect("instances", g.Instances))

	instances := []autoscaling.Instance{}
	for _, i := range g.Instances {
		if *i.LifecycleState == terminatingLifecycleState {
			instances = append(instances, *i)
		}
	}
	w.logger.Info("filter terminating instances", zap.Reflect("instances", instances))

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

		time.Sleep(time.Second * w.config.Watcher.IntervalSec)
	}
}
