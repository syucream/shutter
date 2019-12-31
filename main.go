package main

import (
	"flag"
	"go.uber.org/zap"
	"log"
)

func main() {
	file := flag.String("file", "", "a config file path")
	daemon := flag.Bool("daemon", false, "do as daemon")
	flag.Parse()

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}

	config, err := NewConfig(*file)
	if err != nil {
		log.Fatal(err)
	}

	client, err := NewAwsClient(config)
	if err != nil {
		log.Fatal(err)
	}

	if *daemon {
		err = DoForever(client, config, logger)
	} else {
		err = DoOnce(client, config, logger)
	}

	if err != nil {
		log.Fatal(err)
	}
}
