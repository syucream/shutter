package main

import (
	"flag"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"golang.org/x/sync/errgroup"
	"log"
)

func main() {
	file := flag.String("file", "", "a config file path")
	flag.Parse()

	config, err := NewConfig(*file)
	if err != nil {
		log.Fatal(err)
	}

	client, err := NewAutoscalingClient()
	if err != nil {
		log.Fatal(err)
	}

	eg := errgroup.Group{}
	ch := make(chan autoscaling.Instance, 16)

	watcher := NewWatcher(client, config)
	eg.Go(func() error {
		return watcher.Start(ch)
	})

	eg.Go(func() error {
		for {
			i := <-ch
			finisher := NewFinisher(client, config, *i.InstanceId)

			eg.Go(func() error {
				return finisher.Process()
			})
		}
	})

	err = eg.Wait()
	if err != nil {
		log.Fatal(err)
	}
}
