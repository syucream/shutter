package shutter

import (
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

const maxChanSize = 16

func DoOnce(client AwsClient, config *Config, logger *zap.Logger) error {
	watcher := NewWatcher(client, config, logger)
	instances, err := watcher.Watch()
	if err != nil {
		return err
	}

	eg := errgroup.Group{}
	for _, i := range instances {
		finisher := NewFinisher(client, config, logger, *i.InstanceId)
		eg.Go(func() error {
			return finisher.Process()
		})
	}

	return eg.Wait()
}

func DoOnceWithInstanceId(client AwsClient, config *Config, logger *zap.Logger, instanceId string) error {
	finisher := NewFinisher(client, config, logger, instanceId)
	return finisher.Process()
}

func DoForever(client AwsClient, config *Config, logger *zap.Logger) error {
	eg := errgroup.Group{}
	ch := make(chan autoscaling.Instance, maxChanSize)

	watcher := NewWatcher(client, config, logger)
	eg.Go(func() error {
		return watcher.Start(ch)
	})

	eg.Go(func() error {
		for {
			i := <-ch
			finisher := NewFinisher(client, config, logger, *i.InstanceId)

			eg.Go(func() error {
				return finisher.Process()
			})
		}
	})

	return eg.Wait()
}
