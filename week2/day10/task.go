package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

var (
	ErrTaskNil    = errors.New("task cannot be nil")
	ErrContextNil = errors.New("context cannot be nil")
)

type Task struct {
	ID      int
	Name    string
	timeout time.Duration
}

type Result struct {
	TaskID int
	Output string
	Err    error
}

func ExecuteTask(parent context.Context, task *Task) Result {
	if task == nil {
		return Result{
			Err: ErrTaskNil,
		}
	}
	if parent == nil {
		return Result{
			TaskID: task.ID,
			Err:    ErrContextNil,
		}
	}
	if err := parent.Err(); err != nil {
		return Result{
			TaskID: task.ID,
			Err:    err,
		}
	}
	currContext := parent
	if task.timeout > 0 {
		var cancel context.CancelFunc
		currContext, cancel = context.WithTimeout(parent, task.timeout)
		defer cancel()
	}

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	count := 0

	for {
		select {
		case <-currContext.Done():
			return Result{
				TaskID: task.ID,
				Err:    currContext.Err(),
			}
		case <-ticker.C:
			if err := currContext.Err(); err != nil {
				return Result{
					TaskID: task.ID,
					Err:    err,
				}
			}
			count++
			if count == 10 {
				return Result{
					TaskID: task.ID,
					Output: fmt.Sprintf("Task-%d: completed\n", task.ID),
				}
			}
			fmt.Printf("Task-%d: running %d%%\n", task.ID, count*10)
		}
	}
}
