package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueueConcurrent(t *testing.T) {
	var wg sync.WaitGroup
	var consumed [100]atomic.Int64
	stopSignal := -1
	q := NewQueue()

	wg.Add(3)
	for i := 1; i <= 3; i++ {
		go func(id int) {
			defer wg.Done()
			for {
				task := q.Get()
				if task == stopSignal {
					fmt.Printf("consumer-%d: exit with receive cancel signal\n", id)
					break
				}
				consumed[task].Add(1)
				fmt.Printf("consumer-%d: task=%d\n", id, task)
			}
		}(i)
	}
	for task := range 100 {
		q.Put(task)
	}
	q.Put(stopSignal)
	q.Put(stopSignal)
	q.Put(stopSignal)
	wg.Wait()

	for i := 0; i < len(consumed); i++ {
		assert.Equal(t, int64(1), consumed[i].Load())
	}
}
