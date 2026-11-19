package main

import (
	"errors"
	"fmt"
	"slices"
)

var (
	ErrRepeatTask        = errors.New("cannot add repeat task")
	ErrMaxFailuresExceed = errors.New("engine reached max failures")
)

type TaskSummary struct {
	Running int
	Success int
	Failed  int
	Skipped int
}

type TaskEngine struct {
	Tasks       []*Task
	Handler     TaskHandler
	MaxFailures int
	failedCount int
}

func NewTaskEngine(maxFailures int, handler TaskHandler) *TaskEngine {
	if handler == nil {
		handler = DefaultHandler
	}
	// maxFailures < 0 归零
	// maxFailures == 0 不允许任务失败
	// maxFailures > 0 达到指定失败次数后停止
	if maxFailures < 0 {
		maxFailures = 0
	}
	return &TaskEngine{
		Tasks:       make([]*Task, 0),
		Handler:     handler,
		MaxFailures: maxFailures,
	}
}

func (te *TaskEngine) AddTask(task *Task) error {
	if te == nil {
		return ErrNilReceiver
	}
	if task == nil {
		return ErrTaskNil
	}
	if slices.Contains(te.Tasks, task) {
		return ErrRepeatTask
	}
	te.Tasks = append(te.Tasks, task)
	return nil
}

func (te TaskEngine) GetSummary() TaskSummary {
	ts := TaskSummary{}
	for _, t := range te.Tasks {
		if t == nil {
			continue
		}
		switch t.Status {
		case Pending:
			ts.Skipped++
		case Running:
			ts.Running++
		case Success:
			ts.Success++
		case Failed:
			ts.Failed++
		}
	}
	return ts
}

func (te TaskEngine) PrintSummary() {
	ts := te.GetSummary()
	fmt.Printf("Success: %d\n", ts.Success)
	fmt.Printf("Running: %d\n", ts.Running)
	fmt.Printf("Failed: %d\n", ts.Failed)
	fmt.Printf("Skipped: %d\n", ts.Skipped)
}

func (te *TaskEngine) execute(task *Task) error {
	for {
		if err := te.Handler(task); err == nil {
			return nil
		}
		if err := task.Fail(); err != nil {
			return err
		}
		te.failedCount++
		if te.failedCount >= te.MaxFailures {
			return ErrMaxFailuresExceed
		}
		err := task.RetryTask()
		if err == nil {
			continue
		}
		return err
	}
}

func (te *TaskEngine) Run() error {
	if te == nil {
		return ErrNilReceiver
	}

	if te.Handler == nil {
		te.Handler = DefaultHandler
	}

	for _, t := range te.Tasks {
		if te.failedCount >= te.MaxFailures {
			return ErrMaxFailuresExceed
		}

		if t == nil {
			continue
		}
		if err := t.Start(); err != nil {
			fmt.Printf("engine start task-%d failed: %v\n", t.ID, err)
			continue
		}
		err := te.execute(t)
		if errors.Is(err, ErrMaxFailuresExceed) {
			t.Print()
			fmt.Printf("engine execute task-%d reached max failures: %d, retries: %d\n", t.ID, te.failedCount, t.Retry)
			return ErrMaxFailuresExceed
		}

		if err != nil {
			fmt.Printf("engine execute task-%d failed: %v\n", t.ID, err)
			continue
		}

		if err := t.Complete(); err != nil {
			fmt.Printf("engine complte task-%d failed: %v\n", t.ID, err)
		}
		t.Print()
	}
	return nil
}
