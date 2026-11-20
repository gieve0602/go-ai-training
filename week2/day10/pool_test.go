package main

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// 请新增一个测试：TestConcurrentShutdown。
// 要求：
// - 创建一个包含多个任务的 Pool。
// - 启动 Pool。
// - 使用多个 Goroutine（例如 10 个）同时调用 Shutdown(context.Background())。
// - 等待所有调用者返回。
// - 验证所有调用均返回 nil。
// - 验证 Pool 最终状态为 Terminated。
// - 使用 go test -race 验证。
func TestConcurrentShutdown(t *testing.T) {
	tasks := make([]Task, 0)
	for i := 1; i <= 3; i++ {
		tasks = append(tasks, Task{
			ID:   i,
			Name: fmt.Sprintf("Task-%d", i),
		})
	}

	p, err := NewPool(3, tasks)
	assert.NoError(t, err)
	err = p.Run()
	assert.NoError(t, err)

	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	var wg sync.WaitGroup
	wg.Add(10)
	for range 10 {
		go func() {
			defer wg.Done()
			err := p.Shutdown(ctx)
			assert.NoError(t, err)
		}()
	}

	wg.Wait()
	p.mux.Lock()
	assert.Equal(t, Terminated, p.status)
	p.mux.Unlock()
}

// TestShutdownCallerTimeout
// 要求：
// 1. 创建 Pool，至少 3 个任务。
// 2. 启动 Pool。
// 3. A 使用 context.WithTimeout，设置一个较短的超时时间（例如 100ms）。
// 4. B 使用 context.Background() 调用 Shutdown。
// 5. 验证 A 返回 context.DeadlineExceeded。
// 6. 验证 B 最终返回 nil。
// 7. 验证 Pool 最终状态是 Terminated。
// 8. 使用 go test -race 验证。
func TestShutdownCallerTimeout(t *testing.T) {
	tasks := make([]Task, 0)
	for i := 1; i <= 3; i++ {
		tasks = append(tasks, Task{
			ID:   i,
			Name: fmt.Sprintf("Task-%d", i),
		})
	}

	p, err := NewPool(3, tasks)
	assert.NoError(t, err)
	err = p.Run()
	assert.NoError(t, err)

	ctxA, cancelA := context.WithTimeout(t.Context(), 100*time.Millisecond)
	ctxB, cancelB := context.WithCancel(t.Context())
	defer cancelA()
	defer cancelB()
	errAChan := make(chan error)
	go func() {
		errAChan <- p.Shutdown(ctxA)
	}()

	errBChan := make(chan error)
	go func() {
		errBChan <- p.Shutdown(ctxB)
	}()
	err = <-errAChan
	assert.Equal(t, context.DeadlineExceeded, err)
	err = <-errBChan
	assert.NoError(t, err)
	p.mux.Lock()
	assert.Equal(t, Terminated, p.status)
	p.mux.Unlock()
}
