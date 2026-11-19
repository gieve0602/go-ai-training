package main

import "fmt"

func main() {
	e := NewTaskEngine(0, DefaultHandler)
	for id := range 10 {
		id++
		task := &Task{
			ID:         id,
			Name:       fmt.Sprintf("Task-%d", id),
			MaxRetries: 2,
		}
		err := e.AddTask(task)
		if err != nil {
			fmt.Printf("AddTask error: %v\n", err)
			return
		}
	}

	err := e.Run()
	if err != nil {
		fmt.Printf("engine run failed: %v\n", err)
	}
	e.PrintSummary()
}
