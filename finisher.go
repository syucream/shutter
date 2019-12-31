package main

import (
	"go.uber.org/zap"
	"time"
)

type state int

const (
	initState state = iota
	terminateState
	waitState
	completedState
	finishedState
	abortedState
)

type finisher struct {
	config     *Config
	state      state
	err        error
	instanceId string
	client     *awsClient
	logger     *zap.Logger
}

func (f *finisher) initHandler() state {
	f.logger.Info("init handler")

	return terminateState
}

func (f *finisher) terminateHandler() (state, error) {
	f.logger.Info("terminate handler")

	instance, err := f.client.DescribeInstance(f.instanceId)
	if err != nil {
		return abortedState, err
	}
	cmd := Render(f.config.Finisher.Terminate.Command, instance)
	f.logger.Info("execute terminate command", zap.String("cmd", cmd))

	err = DoCommand(cmd)
	if err != nil {
		return abortedState, err
	}

	return waitState, nil
}

func (f *finisher) waitHandler() (state, error) {
	f.logger.Info("wait handler")

	instance, err := f.client.DescribeInstance(f.instanceId)
	if err != nil {
		return abortedState, err
	}
	cmd := Render(f.config.Finisher.Wait.Command, instance)
	f.logger.Info("execute wait command", zap.String("cmd", cmd))

	for i := int64(0); i < f.config.Finisher.Wait.MaxTries; i++ {
		err := DoCommand(cmd)
		if err == nil {
			break
		}

		f.logger.Info("command failed; will be retried", zap.Error(err))
		time.Sleep(time.Second * f.config.Finisher.Wait.IntervalSec)
	}

	return completedState, nil
}

func (f *finisher) completeHandler() (state, error) {
	f.logger.Info("complete handler")

	err := f.client.CompleteLifecycleAction(f.instanceId, f.config.Finisher.LifecycleActionResult, f.config.Finisher.LifecycleHookName)
	if err != nil {
		return abortedState, err
	}

	return finishedState, nil
}

func (f *finisher) finishedHandler() {
	f.logger.Info("finish handler")
}

func NewFinisher(client *awsClient, c *Config, logger *zap.Logger, instanceId string) *finisher {
	return &finisher{
		config:     c,
		state:      initState,
		err:        nil,
		instanceId: instanceId,
		client:     client,
		logger:     logger,
	}
}

func (f *finisher) Do() {
	var next state
	var err error

	switch f.state {
	case initState:
		next = f.initHandler()
	case terminateState:
		next, err = f.terminateHandler()
	case waitState:
		next, err = f.waitHandler()
	case completedState:
		next, err = f.completeHandler()
	case finishedState:
		f.finishedHandler()
	}

	if err != nil {
		f.logger.Error("error occured in a handler", zap.Error(err))
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
			f.logger.Info("finisher processes are finished")
			break
		}
	}

	return f.err
}
