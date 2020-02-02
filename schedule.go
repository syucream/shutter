package shutter

import (
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"sync"
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
		mux := sync.Mutex{}
		statuses := map[string]bool{} // started instance ids

		for {
			i := <-ch

			func(instanceId string) {
				mux.Lock()
				defer mux.Unlock()

				if started, ok := statuses[instanceId]; ok && started {
					// A finisher has already started before so ignore it
					return
				}
				statuses[instanceId] = true // mark started to prevent reenter

				eg.Go(func() error {
					finisher := NewFinisher(client, config, logger, instanceId)
					err := finisher.Process()

					mux.Lock()
					defer mux.Unlock()
					delete(statuses, instanceId) // release the instance id

					return err
				})
			}(*i.InstanceId)
		}
	})

	return eg.Wait()
}
