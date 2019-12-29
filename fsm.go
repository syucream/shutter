package main

import "log"

type state int
const (
	initState state      = iota
	waitTerminatingState
	terminatingState
	waitCompletionState
	completedState
	finishedState
)

type fsm struct {
	config *Config
	state state
	err error
}

func (f *fsm) initStateHandler() state {
	log.Println(f)
	return waitTerminatingState
}

func (f *fsm) waitTerminatingStateHandler() (state, error) {
	log.Println(f)
	return terminatingState, nil
}

func (f *fsm) terminatingStateHandler() (state, error) {
	log.Println(f)
	return waitCompletionState, nil
}

func (f *fsm) waitCompletionStateHandler() (state, error) {
	log.Println(f)
	return completedState, nil
}

func (f *fsm) completeStateHandler() (state, error) {
	log.Println(f)
	return finishedState, nil
}

func (f *fsm) finishedStateHandler() {
	log.Println(f)
}

func NewFsm(c *Config) *fsm {
	return &fsm{
		config: c,
		state: initState,
		err: nil,
	}
}

func (f *fsm) Do() {
	var next state
	var err error

	switch f.state {
	case initState:
		next = f.initStateHandler()
	case waitTerminatingState:
		next, err = f.waitTerminatingStateHandler()
	case terminatingState:
		next, err = f.terminatingStateHandler()
	case waitCompletionState:
		next, err = f.waitCompletionStateHandler()
	case completedState:
		next, err = f.completeStateHandler()
	case finishedState:
		f.finishedStateHandler()
	}

	if err != nil {
		f.err = err
		next = finishedState
	}
	f.state = next
}

func (f *fsm) IsFinished() bool {
	return f.state == finishedState
}

func (f *fsm) Start() error {
	for {
		f.Do()
		if f.IsFinished() {
			break
		}
	}

	return f.err
}
