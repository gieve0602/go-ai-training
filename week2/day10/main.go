package main

import (
	"context"
	"fmt"
)

func main() {
	tasks := make([]Task, 0)
	for i := 1; i <= 4; i++ {
		tasks = append(tasks, Task{
			ID:   i,
			Name: fmt.Sprintf("Task-%d", i),
		})
	}

	p, err := NewPool(3, tasks)
	if err != nil {
		fmt.Printf("Pool: create failed %v\n", err)
		return
	}
	err = p.Run()
	if err != nil {
		fmt.Printf("Pool: run failed %v\n", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err = p.Shutdown(ctx)
	if err != nil {
		fmt.Printf("Pool: shutdown failed %v\n", err)
	}
}
