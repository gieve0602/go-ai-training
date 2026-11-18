package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalStatusTransformer(t *testing.T) {
	task := &Task{
		ID:   1,
		Name: "Task-1",
	}
	assert.NoError(t, task.Start())
	assert.NoError(t, task.Complete())
	assert.Equal(t, Success, task.Status)
	task2 := &Task{
		ID:   2,
		Name: "Task-2",
	}
	assert.NoError(t, task2.Start())
	assert.NoError(t, task2.Fail())
	assert.Equal(t, Failed, task2.Status)
}

func TestIllegalStatusTransformer(t *testing.T) {
	task := &Task{
		ID:   1,
		Name: "Task-1",
	}
	assert.NotNil(t, task.Fail())
	assert.Equal(t, Pending, task.Status)
	task2 := &Task{
		ID:   2,
		Name: "Task-2",
	}
	assert.NotNil(t, task2.Complete())
	assert.Equal(t, Pending, task2.Status)
	task3 := &Task{
		ID:     3,
		Name:   "Task-3",
		Status: Success,
	}
	assert.NotNil(t, task3.Start())
	task4 := &Task{
		ID:     4,
		Name:   "Task-4",
		Status: Failed,
		Retry:  MaxFailures,
	}
	assert.Error(t, task4.Start())
	assert.Equal(t, Failed, task4.Status)
	task5 := &Task{
		ID:     5,
		Name:   "Task-5",
		Status: Running,
	}
	assert.NotNil(t, task5.Start())
}

func TestTaskRetry(t *testing.T) {
	task := &Task{
		ID:   1,
		Name: "Task-1",
	}
	assert.NoError(t, task.Start())
	assert.NoError(t, task.Fail())
	assert.Equal(t, 1, task.Retry)
}

func TestGetStatusName(t *testing.T) {
	task := &Task{
		ID:   1,
		Name: "Task-1",
	}
	assert.Equal(t, "Pending", task.GetStatusName())
	task.Status = Running
	assert.Equal(t, "Running", task.GetStatusName())
	task.Status = Success
	assert.Equal(t, "Success", task.GetStatusName())
	task.Status = Failed
	assert.Equal(t, "Failed", task.GetStatusName())
	task.Status = -1
	assert.Equal(t, "Unknown", task.GetStatusName())
}

func TestPrintTo(t *testing.T) {
	task := Task{ID: 1, Status: Running}
	var buf bytes.Buffer
	err := task.PrintTo(&buf)
	assert.NoError(t, err)
	assert.Equal(t, "Task-1: Running\n", buf.String())

	buf.Reset()
	task.Status = Pending
	err = task.PrintTo(&buf)
	assert.NoError(t, err)
	assert.Equal(t, "Task-1: Pending (Skipped)\n", buf.String())

	buf.Reset()
	task.Status = Success
	err = task.PrintTo(&buf)
	assert.NoError(t, err)
	assert.Equal(t, "Task-1: Success\n", buf.String())

	buf.Reset()
	task.Status = Failed
	err = task.PrintTo(&buf)
	assert.NoError(t, err)
	assert.Equal(t, "Task-1: Failed\n", buf.String())

	buf.Reset()
	task.Status = -1
	err = task.PrintTo(&buf)
	assert.NoError(t, err)
	assert.Equal(t, "Task-1: Unknown\n", buf.String())
}

func TestNilTask(t *testing.T) {
	var task *Task
	assert.ErrorIs(t, task.Start(), ErrNilTask)
	assert.ErrorIs(t, task.Complete(), ErrNilTask)
	assert.ErrorIs(t, task.Fail(), ErrNilTask)
}
