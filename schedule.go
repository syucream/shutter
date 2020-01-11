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

	mux := sync.Mutex{}
	statuses := map[autoscaling.Instance]bool{} // instance id -> isStarted

	eg.Go(func() error {
		for {
			i := <-ch

			if started, ok := statuses[i]; ok && started {
				// A finisher has already started before so ignore it
				continue
			}
			mux.Lock()
			statuses[i] = true // mark started to prevent reenter
			mux.Unlock()

			finisher := NewFinisher(client, config, logger, *i.InstanceId)
			eg.Go(func() error {
				err := finisher.Process()

				mux.Lock()
				delete(statuses, i) // release the instance id
				mux.Unlock()

				return err
			})
		}
	})

	return eg.Wait()
}
