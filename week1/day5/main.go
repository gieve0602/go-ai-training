package main

import "fmt"

func main() {
	e := NewTaskEngine(2)
	for id := range 10 {
		id++
		task := &Task{
			ID:   id,
			Name: fmt.Sprintf("Task-%d", id),
		}
		err := e.AddTask(task)
		if err != nil {
			fmt.Printf("AddTask error: %v\n", err)
			return
		}
	}

	e.Run()
	e.PrintSummary()
}
