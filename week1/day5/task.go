package main

import (
	"errors"
	"fmt"
	"io"
	"os"
)

var ErrNilTask = errors.New("nil task cannot call method")

const MaxFailures = 3
const (
	Pending = iota
	Running
	Success
	Failed
)

type Task struct {
	ID     int
	Name   string
	Status int
	Retry  int
}

func (t *Task) Start() error {
	if t == nil {
		return ErrNilTask
	}
	if t.Status != Pending && t.Status != Failed {
		return fmt.Errorf("task %d cannot running from status %s", t.ID, t.GetStatusName())
	}
	if t.Status == Failed && t.Retry >= MaxFailures {
		return fmt.Errorf("task %d has retry reached max retry times", t.ID)
	}
	t.Status = Running
	return nil
}

func (t *Task) Complete() error {
	if t == nil {
		return ErrNilTask
	}
	if t.Status != Running {
		return fmt.Errorf("task %d cannot success from status %s", t.ID, t.GetStatusName())
	}
	t.Status = Success
	return nil
}

func (t *Task) Fail() error {
	if t == nil {
		return ErrNilTask
	}
	if t.Status != Running {
		return fmt.Errorf("task %d cannot fail from status %s", t.ID, t.GetStatusName())
	}
	t.Status = Failed
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

func (t Task) Print() error {
	return t.PrintTo(os.Stdout)
}

func (t Task) PrintTo(w io.Writer) error {
	status := t.GetStatusName()
	if t.Status == Pending {
		status = "Pending (Skipped)"
	}
	_, err := fmt.Fprintf(w, "Task-%d: %s\n", t.ID, status)
	return err
}
