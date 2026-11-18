package main

import (
	"fmt"
	"log"
)

const MaxFailed = 2

type TaskEngine struct {
	tasks     []*Task
	maxFailed int
}

func NewTaskEngine(maxFailed int) *TaskEngine {
	if maxFailed <= 0 {
		maxFailed = MaxFailed
	}
	return &TaskEngine{
		tasks:     make([]*Task, 0),
		maxFailed: maxFailed,
	}
}

func (e *TaskEngine) AddTask(task *Task) error {
	if task == nil {
		return ErrNilTask
	}
	e.tasks = append(e.tasks, task)
	return nil
}

func (e *TaskEngine) Run() {
	failCount := 0
	for _, t := range e.tasks {
		if t.ID%9 == 0 {
			err := t.Print()
			if err != nil {
				fmt.Printf("print task info error: %v\n", err)
			}
			continue
		}
		if err := t.Start(); err != nil {
			log.Printf("Task-%d start failed: %v\n", t.ID, err)
			continue
		}
		switch {
		case t.ID%5 == 0:
			if err := t.Fail(); err != nil {
				continue
			}
			failCount++
		case t.ID%3 == 0:
			if err := t.Complete(); err != nil {
				continue
			}
		}
		err := t.Print()
		if err != nil {
			fmt.Printf("print task info error: %v\n", err)
		}
		if failCount >= e.maxFailed {
			fmt.Println("Failure threshold reached")
			break
		}
	}
}

type TaskSummary struct {
	Skipped int
	Running int
	Success int
	Failed  int
}

func (e *TaskEngine) GetSummary() TaskSummary {
	ts := TaskSummary{}
	for _, t := range e.tasks {
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

func (e *TaskEngine) PrintSummary() {
	ts := e.GetSummary()
	fmt.Printf("Success: %d\n", ts.Success)
	fmt.Printf("Running: %d\n", ts.Running)
	fmt.Printf("Failed: %d\n", ts.Failed)
	fmt.Printf("Skipped: %d\n", ts.Skipped)
}
