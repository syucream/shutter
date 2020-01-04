package main

import (
	"flag"
	"github.com/syucream/shutter"
	"go.uber.org/zap"
	"log"
)

func main() {
	file := flag.String("file", "", "a config file path")
	daemon := flag.Bool("daemon", false, "do as daemon")
	instanceId := flag.String("instanceid", "", "EC2 instance id (optional, used if daemon = false)")
	flag.Parse()

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}

	config, err := shutter.NewConfig(*file)
	if err != nil {
		log.Fatal(err)
	}

	client, err := shutter.NewAwsClient(config)
	if err != nil {
		log.Fatal(err)
	}

	if *daemon {
		err = shutter.DoForever(client, config, logger)
	} else {
		if *instanceId != "" {
			err = shutter.DoOnceWithInstanceId(client, config, logger, *instanceId)
		} else {
			err = shutter.DoOnce(client, config, logger)
		}
	}

	if err != nil {
		log.Fatal(err)
	}
}
