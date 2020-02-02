package shutter

import (
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"go.uber.org/zap"
	"time"
)

var (
	inServiceLifecycleState   = "InService"
	terminatingLifecycleState = "Terminating:Wait"
)

// watcher watches terminating EC2 instances under a Lifecycle Hook
type watcher struct {
	client AwsClient
	config *Config
	logger *zap.Logger
}

func NewWatcher(client AwsClient, config *Config, logger *zap.Logger) *watcher {
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

func (w *watcher) Notify(channel chan autoscaling.Instance) error {
	instances, err := w.Watch()
	if err != nil {
		return err
	}

	for _, i := range instances {
		channel <- i
	}

	return nil
}

func (w *watcher) Start(channel chan autoscaling.Instance) error {
	for {
		if err := w.Notify(channel); err != nil {
			return err
		}

		time.Sleep(time.Second * w.config.Watcher.IntervalSec)
	}
}
