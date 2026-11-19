package main

import "fmt"

const (
	Pending = iota
	Running
	Success
	Failed
)

const (
	MaxFailed = 2
	MaxTasks  = 15
)

func printSummary(skipped, running, success, failed int, successTaskIDs []int) {
	fmt.Println("=========================")
	fmt.Printf("Task skipped: %d\n", skipped)
	fmt.Printf("Task running: %d\n", running)
	fmt.Printf("Task success: %d\n", success)
	fmt.Printf("Task failed: %d\n", failed)
	fmt.Printf("Task success sum: %d\n", sum(successTaskIDs...))
	fmt.Println("=========================")
}

func getTaskStatus(id int) (int, error) {
	if id < 0 {
		return 0, fmt.Errorf("invalid task id: %d", id)
	}
	switch {
	case id%9 == 0:
		return Pending, nil
	case id%5 == 0:
		return Failed, nil
	case id%3 == 0:
		return Success, nil
	default:
		return Running, nil
	}
}

func getStatusName(status int) (string, error) {
	switch status {
	case Pending:
		return "Pending", nil
	case Running:
		return "Running", nil
	case Success:
		return "Success", nil
	case Failed:
		return "Failed", nil
	default:
		return "", fmt.Errorf("status is invalid: %v", status)
	}
}

func processTask(id int) (int, error) {
	status, err := getTaskStatus(id)
	if err != nil {
		return status, err
	}
	if status == Pending {
		fmt.Printf("Task-%d: Skipped (Pending)\n", id)
		return status, nil
	}
	statusName, err := getStatusName(status)
	if err != nil {
		return status, err
	}
	fmt.Printf("Task-%d: statusCode=%d, statusName=%v\n", id, status, statusName)
	return status, nil
}

func runScheduler(maxTasks int) error {
	if maxTasks < 1 || maxTasks > MaxTasks {
		return fmt.Errorf("maxTasks is invalud: %v", maxTasks)
	}
	skippedCount := 0
	runningCount := 0
	successCount := 0
	failedCount := 0
	var successTaskIDs []int

	for i := 1; i <= maxTasks; i++ {
		status, err := processTask(i)
		if err != nil {
			return err
		}
		switch status {
		case Pending:
			skippedCount++
		case Running:
			runningCount++
		case Success:
			successCount++
			successTaskIDs = append(successTaskIDs, i)
		case Failed:
			failedCount++
		}
		if failedCount >= 2 {
			fmt.Printf("failure threshold reached, stopping scheduler\n")
			break
		}
	}
	printSummary(skippedCount, runningCount, successCount, failedCount, successTaskIDs)
	return nil
}
