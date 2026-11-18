package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddNormalTask(t *testing.T) {
	task := &Task{
		ID:         1,
		Name:       "Task-1",
		MaxRetries: 2,
	}
	e := NewTaskEngine(2, AlternateHandler)
	err := e.AddTask(task)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(e.Tasks))
	assert.Equal(t, task, e.Tasks[0])
}

func TestAddNilTask(t *testing.T) {
	var task *Task
	e := NewTaskEngine(2, AlternateHandler)
	err := e.AddTask(task)
	assert.ErrorIs(t, err, ErrTaskNil)
	assert.Equal(t, 0, len(e.Tasks))
}

func TestRepeatAddSameTask(t *testing.T) {
	task := &Task{
		ID:         1,
		Name:       "Task-1",
		MaxRetries: 2,
	}
	e := NewTaskEngine(2, AlternateHandler)
	err := e.AddTask(task)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(e.Tasks))
	assert.Equal(t, task, e.Tasks[0])
	err = e.AddTask(task)
	assert.ErrorIs(t, err, ErrRepeatTask)
	assert.Equal(t, 1, len(e.Tasks))
}

func TestGetSummary(t *testing.T) {
	task1 := &Task{
		ID:   1,
		Name: "task1",
	}
	task2 := &Task{
		ID:     2,
		Name:   "task2",
		Status: Running,
	}
	task3 := &Task{
		ID:     3,
		Name:   "task3",
		Status: Success,
	}
	task4 := &Task{
		ID:     4,
		Name:   "task4",
		Status: Failed,
	}
	task5 := &Task{
		ID:   5,
		Name: "task5",
	}
	e := NewTaskEngine(2, DefaultHandler)
	err := e.AddTask(task1)
	assert.NoError(t, err)
	err = e.AddTask(task2)
	assert.NoError(t, err)
	err = e.AddTask(task3)
	assert.NoError(t, err)
	err = e.AddTask(task4)
	assert.NoError(t, err)
	err = e.AddTask(task5)
	assert.NoError(t, err)
	assert.Equal(t, 5, len(e.Tasks))

	e.Tasks = append(e.Tasks, nil)

	ts := e.GetSummary()
	assert.Equal(t, 1, ts.Running)
	assert.Equal(t, 1, ts.Success)
	assert.Equal(t, 1, ts.Failed)
	assert.Equal(t, 2, ts.Skipped)
}

func TestRunWithAlternateHandler(t *testing.T) {
	task3 := &Task{
		ID:   3,
		Name: "task3",
	}
	task6 := &Task{
		ID:   6,
		Name: "task6",
	}
	task9 := &Task{
		ID:   9,
		Name: "task9",
	}

	e := NewTaskEngine(2, AlternateHandler)
	err := e.AddTask(task3)
	assert.NoError(t, err)
	err = e.AddTask(task6)
	assert.NoError(t, err)
	err = e.AddTask(task9)
	assert.NoError(t, err)

	err = e.Run()
	assert.NoError(t, err)
	assert.Equal(t, Success, task3.Status)
	assert.Equal(t, 0, task3.Retry)
	assert.Equal(t, Success, task6.Status)
	assert.Equal(t, 0, task6.Retry)
	assert.Equal(t, Success, task9.Status)
	assert.Equal(t, 0, task9.Retry)
	ts := e.GetSummary()
	assert.Equal(t, 3, ts.Success)
	assert.Equal(t, 0, ts.Running)
	assert.Equal(t, 0, ts.Failed)
	assert.Equal(t, 0, ts.Skipped)
}

func TestRunRetrySuccess(t *testing.T) {
	count := 0
	handler := func(t *Task) error {
		count++
		if count == 1 {
			return ErrTaskFailed
		}
		return nil
	}

	e := NewTaskEngine(3, handler)
	task := &Task{
		ID:         1,
		Name:       "Task-1",
		MaxRetries: 2,
	}
	err := e.AddTask(task)
	assert.NoError(t, err)
	err = e.Run()
	assert.NoError(t, err)
	assert.Equal(t, 2, count)
	assert.Equal(t, Success, task.Status)
	assert.Equal(t, 1, task.Retry)
	assert.Equal(t, 1, e.failedCount)
	ts := e.GetSummary()
	assert.Equal(t, 1, ts.Success)
	assert.Equal(t, 0, ts.Failed)
}

func TestRunRetryExhausted(t *testing.T) {
	count := 0
	count2 := 0
	handler := func(t *Task) error {
		if t.ID == 1 {
			count++
			return ErrTaskFailed
		}
		count2++
		return nil
	}

	e := NewTaskEngine(10, handler)
	task := &Task{
		ID:         1,
		Name:       "Task-1",
		MaxRetries: 2,
	}
	task2 := &Task{
		ID:         2,
		Name:       "Task-2",
		MaxRetries: 2,
	}
	err := e.AddTask(task)
	assert.NoError(t, err)
	err = e.AddTask(task2)
	assert.NoError(t, err)

	err = e.Run()
	ts := e.GetSummary()
	assert.NoError(t, err)
	assert.Equal(t, 3, count)
	assert.Equal(t, 1, count2)
	assert.Equal(t, 3, e.failedCount)
	assert.Equal(t, Failed, task.Status)
	assert.Equal(t, 2, task.Retry)
	assert.Equal(t, Success, task2.Status)
	assert.Equal(t, 0, task2.Retry)
	assert.Equal(t, 1, ts.Failed)
	assert.Equal(t, 1, ts.Success)
}

func TestRunMaxFailures(t *testing.T) {
	count := 0
	handler := func(t *Task) error {
		count++
		return ErrTaskFailed
	}

	e := NewTaskEngine(2, handler)
	task1 := &Task{
		ID:         1,
		Name:       "Task-1",
		MaxRetries: 3,
	}
	task2 := &Task{
		ID:         2,
		Name:       "Task-2",
		MaxRetries: 3,
	}
	task3 := &Task{
		ID:         3,
		Name:       "Task-3",
		MaxRetries: 3,
	}
	err := e.AddTask(task1)
	assert.NoError(t, err)
	err = e.AddTask(task2)
	assert.NoError(t, err)
	err = e.AddTask(task3)
	assert.NoError(t, err)

	err = e.Run()
	assert.ErrorIs(t, err, ErrMaxFailuresExceed)
	assert.Equal(t, 2, count)
	assert.Equal(t, 2, e.failedCount)
	assert.Equal(t, Failed, task1.Status)
	assert.Equal(t, 1, task1.Retry)
	assert.Equal(t, Pending, task2.Status)
	assert.Equal(t, Pending, task3.Status)
	ts := e.GetSummary()
	assert.Equal(t, 1, ts.Failed)
	assert.Equal(t, 2, ts.Skipped)
}

func TestRunMaxFailuresZero(t *testing.T) {
	count := 0
	handler := func(t *Task) error {
		count++
		return ErrTaskFailed
	}
	e := NewTaskEngine(0, handler)
	task1 := &Task{
		ID:         1,
		Name:       "Task-1",
		MaxRetries: 3,
	}
	task2 := &Task{
		ID:         2,
		Name:       "Task-2",
		MaxRetries: 3,
	}
	task3 := &Task{
		ID:         3,
		Name:       "Task-3",
		MaxRetries: 3,
	}

	err := e.AddTask(task1)
	assert.NoError(t, err)
	err = e.AddTask(task2)
	assert.NoError(t, err)
	err = e.AddTask(task3)
	assert.NoError(t, err)

	err = e.Run()
	ts := e.GetSummary()
	assert.ErrorIs(t, err, ErrMaxFailuresExceed)
	assert.Equal(t, 0, count)
	assert.Equal(t, 0, e.failedCount)
	assert.Equal(t, Pending, task1.Status)
	assert.Equal(t, Pending, task2.Status)
	assert.Equal(t, 3, ts.Skipped)
	assert.Equal(t, 0, ts.Failed)
}

// TestAddTaskNilEngine：nil Engine 调用 AddTask，返回 ErrNilReceiver。
func TestAddTaskNilEngine(t *testing.T) {
	var e *TaskEngine
	err := e.AddTask(nil)
	assert.ErrorIs(t, err, ErrNilReceiver)
}

// TestRunNilEngine：nil Engine 调用 Run，返回 ErrNilReceiver。
func TestRunNilEngine(t *testing.T) {
	var e *TaskEngine
	err := e.Run()
	assert.ErrorIs(t, err, ErrNilReceiver)
}

// TestNewTaskEngineNegativeMaxFailures：传入负数，验证 MaxFailures 被归零。
func TestNewTaskEngineNegativeMaxFailures(t *testing.T) {
	e := NewTaskEngine(-1, DefaultHandler)
	assert.Equal(t, 0, e.MaxFailures)
}

// TestNewTaskEngineNilHandler：传入 nil Handler，验证默认使用 DefaultHandler。
func TestNewTaskEngineNilHandler(t *testing.T) {
	e := NewTaskEngine(2, nil)
	task := &Task{
		ID:   2,
		Name: "task2",
	}
	err := e.Handler(task)
	assert.NoError(t, err)
}
