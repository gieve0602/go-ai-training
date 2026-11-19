package main

import (
	"testing"

	"github.com/go-playground/assert"
)

func TestCalculateFunc(t *testing.T) {
	x, y := 6, 8
	r1 := execute(x, y, add)
	r2 := execute(x, y, subtract)
	r3 := execute(x, y, multiply)
	assert.Equal(t, r1, 14)
	assert.Equal(t, r2, -2)
	assert.Equal(t, r3, 48)
}

func TestAnonymousFunc(t *testing.T) {
	calcuate := func(x, y int) int {
		return x + y
	}
	r := calcuate(10, 5)

	assert.Equal(t, r, 15)
}

func TestClosure(t *testing.T) {
	closure1 := createTaskCounter()
	closure2 := createTaskCounter()

	tests := []struct {
		name string
		got  int
		want int
	}{
		{"closure1", closure1(), 1},
		{"closure1", closure1(), 2},
		{"closure2", closure2(), 1},
	}

	for _, tt := range tests {
		if tt.got != tt.want {
			t.Errorf("%s: got=%d, want=%d\n", tt.name, tt.got, tt.want)
		}
	}
}
