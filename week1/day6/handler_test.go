package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultHandler(t *testing.T) {
	task1 := &Task{
		ID: 1,
	}
	task2 := &Task{
		ID: 2,
	}
	err1 := DefaultHandler(task1)
	err2 := DefaultHandler(task2)
	assert.ErrorIs(t, err1, ErrTaskFailed)
	assert.NoError(t, err2)
}

func TestDefaultHandlerNilTask(t *testing.T) {
	err := DefaultHandler(nil)
	assert.ErrorIs(t, err, ErrTaskNil)
}

func TestAlternateHandler(t *testing.T) {
	task1 := &Task{
		ID: 2,
	}
	task2 := &Task{
		ID: 3,
	}
	err1 := AlternateHandler(task1)
	err2 := AlternateHandler(task2)
	assert.ErrorIs(t, err1, ErrTaskFailed)
	assert.NoError(t, err2)
}

func TestAlternateHandlerNilTask(t *testing.T) {
	err := AlternateHandler(nil)
	assert.ErrorIs(t, err, ErrTaskNil)
}
