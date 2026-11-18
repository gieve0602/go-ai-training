package main

import (
	"errors"
	"fmt"
)

type TaskStatus int

const (
	Pending TaskStatus = iota
	Running
	Success
	Failed
)

var (
	ErrNilReceiver      = errors.New("method receiver cannot be nil")
	ErrMaxRetriesExceed = errors.New("task retry times reach max retries")
)

type Task struct {
	ID         int
	Name       string
	Status     TaskStatus
	Retry      int
	MaxRetries int
}

func (t *Task) Start() error {
	if t == nil {
		return ErrNilReceiver
	}
	if t.Status != Pending {
		return fmt.Errorf("task-%d cannot transition to Running from %s", t.ID, t.GetStatusName())
	}
	t.Status = Running
	return nil
}

func (t *Task) Complete() error {
	if t == nil {
		return ErrNilReceiver
	}
	if t.Status != Running {
		return fmt.Errorf("task-%d cannot transition to Success from %s", t.ID, t.GetStatusName())
	}
	t.Status = Success
	return nil
}

func (t *Task) Fail() error {
	if t == nil {
		return ErrNilReceiver
	}
	if t.Status != Running {
		return fmt.Errorf("task-%d cannot transition to Failed from %s", t.ID, t.GetStatusName())
	}
	t.Status = Failed
	return nil
}

func (t *Task) RetryTask() error {
	if t == nil {
		return ErrNilReceiver
	}
	if t.Status != Failed {
		return fmt.Errorf("task-%d cannot transition to Running from %s", t.ID, t.GetStatusName())
	}
	if t.Retry >= t.MaxRetries {
		return ErrMaxRetriesExceed
	}
	t.Status = Running
	t.Retry++
	return nil
}

func (t Task) GetStatusName() string {
	switch t.Status {
	case Pending:
		return "Pending"
	case Running:
		return "Running"
	case Success:
		return "Success"
	case Failed:
		return "Failed"
	}
	return "Unknown"
}

func (t Task) Print() {
	fmt.Printf("Task-%d: %s\n", t.ID, t.GetStatusName())
}
