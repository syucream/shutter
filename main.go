package main

import (
	"flag"
	"log"
)

func main() {
	file := flag.String("file", "", "a config file path")
	flag.Parse()

	config, err := NewConfig(*file)
	if err != nil {
		log.Fatal(err)
	}

	shutter := NewFsm(config)
	err = shutter.Start()
	if err != nil {
		log.Fatal(err)
	}
}
