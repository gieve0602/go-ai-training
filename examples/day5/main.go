package main

import "fmt"

type Task struct {
	ID     int
	Name   string
	Status int
}

func (t Task) Print() {
	fmt.Printf("Task=%#v\n", t)
	t.Status = 1
}

func (t *Task) Print2() {
	fmt.Printf("Task=%#v\n", t)
}

func main() {
	task := &Task{}
	task.Print()
	fmt.Printf("task status: %d\n", task.Status)
	task.Print2()
	task2 := Task{}
	task2.Print()
	task2.Print2()
}
