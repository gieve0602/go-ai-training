package main

import "errors"

var (
	ErrTaskNil    = errors.New("parameter task cannot be nil")
	ErrTaskFailed = errors.New("task execute failed")
)

type TaskHandler func(*Task) error

func DefaultHandler(t *Task) error {
	if t == nil {
		return ErrTaskNil
	}
	if t.ID%2 != 0 {
		return ErrTaskFailed
	}
	return nil
}

func AlternateHandler(t *Task) error {
	if t == nil {
		return ErrTaskNil
	}
	if t.ID%3 != 0 {
		return ErrTaskFailed
	}
	return nil
}
