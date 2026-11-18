package main

import (
	"fmt"
)

const (
	Pending = iota
	Running
	Success
	Failed
)

func main() {
	maxConnections := 100
	currentConnections := 80
	maxTasks := 15
	schedulerEnabled := true
	debug := false

	successCount := 0
	failedCount := 0
	runningCount := 0
	skippedCount := 0

	if !schedulerEnabled {
		fmt.Println("Disabled")
		return
	}
	if currentConnections >= maxConnections {
		fmt.Println("Connection limit reached")
		return
	}

	if debug {
		fmt.Println("dubugging")
	}
	for i := 1; i <= maxTasks; i++ {
		taskID := i
		taskStatus := Running
		if taskID%9 == 0 {
			taskStatus = Pending
			skippedCount++
		} else if taskID%5 == 0 {
			taskStatus = Failed
			failedCount++
		} else if taskID%3 == 0 {
			taskStatus = Success
			successCount++
		} else {
			runningCount++
		}
		switch taskStatus {
		case Pending:
			fmt.Printf("Task-%d: Skipped (Pending)\n", taskID)
			continue
		case Running:
			fmt.Printf("Task-%d: Running\n", taskID)
		case Success:
			fmt.Printf("Task-%d: Success\n", taskID)
		case Failed:
			fmt.Printf("Task-%d: Failed\n", taskID)
		}
		if failedCount >= 2 {
			fmt.Println("Failure threshold reached, stopping scheduler")
			break
		}
	}
	fmt.Printf("Success task count: %v\n", successCount)
	fmt.Printf("Running task count: %v\n", runningCount)
	fmt.Printf("Failed task count: %v\n", failedCount)
	fmt.Printf("Skipped task count: %v\n", skippedCount)
}
