package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTaskPendingToRunning(t *testing.T) {
	task := &Task{
		ID:         1,
		Name:       "Task-1",
		MaxRetries: 2,
	}
	err := task.Start()
	assert.NoError(t, err)
	assert.Equal(t, Running, task.Status)
}

func TestTaskRunningToSuccess(t *testing.T) {
	task := &Task{
		ID:         1,
		Name:       "Task-1",
		MaxRetries: 2,
		Status:     Running,
	}
	err := task.Complete()
	assert.NoError(t, err)
	assert.Equal(t, Success, task.Status)
}

func TestTaskRunningToFailed(t *testing.T) {
	task := &Task{
		ID:         1,
		Name:       "Task-1",
		MaxRetries: 2,
		Status:     Running,
	}
	err := task.Fail()
	assert.NoError(t, err)
	assert.Equal(t, Failed, task.Status)
}

func TestTaskSuccessToRunning(t *testing.T) {
	task := &Task{
		ID:         1,
		Name:       "Task-1",
		MaxRetries: 2,
		Status:     Success,
	}
	err := task.Start()
	assert.NotNil(t, err)
}

func TestTaskPendingToSuccess(t *testing.T) {
	task := &Task{
		ID:         1,
		Name:       "Task-1",
		MaxRetries: 2,
	}
	err := task.Complete()
	assert.NotNil(t, err)
}

func TestTaskRetrySuccess(t *testing.T) {
	task := &Task{
		ID:         1,
		Name:       "Task-1",
		MaxRetries: 2,
	}
	err := task.Start()
	assert.NoError(t, err)
	err = task.Fail()
	assert.NoError(t, err)
	err = task.RetryTask()
	assert.NoError(t, err)
	assert.Equal(t, Running, task.Status)
	assert.Equal(t, 1, task.Retry)
}

func TestMaxRetries(t *testing.T) {
	task := &Task{
		ID:         1,
		Name:       "Task-1",
		MaxRetries: 2,
	}

	err := task.Start()
	assert.NoError(t, err)
	for idx := range task.MaxRetries {
		err = task.Fail()
		assert.NoError(t, err)
		err = task.RetryTask()
		assert.NoError(t, err)
		assert.Equal(t, Running, task.Status)
		assert.Equal(t, idx+1, task.Retry)
	}

	err = task.Fail()
	assert.NoError(t, err)
	err = task.RetryTask()
	assert.ErrorIs(t, err, ErrMaxRetriesExceed)
	assert.Equal(t, Failed, task.Status)
	assert.Equal(t, 2, task.Retry)
}

func TestNilReceiver(t *testing.T) {
	var task *Task
	err := task.Start()
	assert.ErrorIs(t, err, ErrNilReceiver)
}
