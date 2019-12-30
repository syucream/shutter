package main

import (
	"log"
)

type state int

const (
	initState state = iota
	startCompletionState
	waitCompletionState
	completedState
	finishedState
	abortedState
)

type finisher struct {
	config *Config
	state  state
	err    error

	instanceId        string
	autoscalingClinet *autoscalingClient
}

func (f *finisher) initStateHandler() state {
	log.Println(f)
	return startCompletionState
}

func (f *finisher) startCompletionStateHandler() (state, error) {
	err := DoCommand(f.config.Finisher.StartCompletionCommand)
	if err != nil {
		return abortedState, err
	}

	return waitCompletionState, nil
}

func (f *finisher) waitCompletionStateHandler() (state, error) {
	err := DoCommand(f.config.Finisher.WaitCompletionCommand)
	if err != nil {
		return abortedState, err
	}

	return completedState, nil
}

func (f *finisher) completeStateHandler() (state, error) {
	err := f.autoscalingClinet.CompleteLifecycleAction(f.instanceId, f.config.Finisher.LifecycleHookName)
	if err != nil {
		return abortedState, err
	}

	return finishedState, nil
}

func (f *finisher) finishedStateHandler() {
	log.Println(f)
}

func NewFinisher(client *autoscalingClient, c *Config, instanceId string) *finisher {
	return &finisher{
		config:            c,
		state:             initState,
		err:               nil,
		instanceId:        instanceId,
		autoscalingClinet: client,
	}
}

func (f *finisher) Do() {
	var next state
	var err error

	switch f.state {
	case initState:
		next = f.initStateHandler()
	case startCompletionState:
		next, err = f.startCompletionStateHandler()
	case waitCompletionState:
		next, err = f.waitCompletionStateHandler()
	case completedState:
		next, err = f.completeStateHandler()
	case finishedState:
		f.finishedStateHandler()
	}

	if err != nil {
		f.err = err
	}
	f.state = next
}

func (f *finisher) IsFinished() bool {
	return f.state == finishedState || f.state == abortedState
}

func (f *finisher) Process() error {
	for {
		f.Do()
		if f.IsFinished() {
			break
		}
	}

	return f.err
}
