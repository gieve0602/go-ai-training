package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEngineSummary(t *testing.T) {
	e := NewTaskEngine(2)
	for id := range 10 {
		id++
		task := &Task{
			ID:   id,
			Name: fmt.Sprintf("Task-%d", id),
		}
		err := e.AddTask(task)
		assert.NoError(t, err)
	}

	e.Run()
	ts := e.GetSummary()
	assert.Equal(t, 2, ts.Success)
	assert.Equal(t, 5, ts.Running)
	assert.Equal(t, 2, ts.Failed)
	assert.Equal(t, 1, ts.Skipped)
}

func TestEngineAddNilTask(t *testing.T) {
	var task *Task
	e := NewTaskEngine(2)
	err := e.AddTask(task)
	assert.ErrorIs(t, err, ErrNilTask)
}

func TestEngineMaxFailed(t *testing.T) {
	e := NewTaskEngine(1)
	for id := range 10 {
		id++
		task := &Task{
			ID:   id,
			Name: fmt.Sprintf("Task-%d", id),
		}
		err := e.AddTask(task)
		assert.NoError(t, err)
	}
	e.Run()
	ts := e.GetSummary()
	assert.Equal(t, 1, ts.Failed)
	assert.Equal(t, 3, ts.Running)
	assert.Equal(t, 1, ts.Success)
	assert.Equal(t, 5, ts.Skipped)
}
