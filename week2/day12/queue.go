package main

import "sync"

type Queue struct {
	mu    sync.Mutex
	cond  *sync.Cond
	tasks []int
}

func NewQueue() *Queue {
	q := &Queue{
		tasks: make([]int, 0),
	}

	q.cond = sync.NewCond(&q.mu)
	return q
}

func (q *Queue) Put(task int) {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()
	q.tasks = append(q.tasks, task)
	q.cond.Signal()
}

func (q *Queue) Get() int {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()
	for len(q.tasks) == 0 {
		q.cond.Wait()
	}
	task := q.tasks[0]
	q.tasks = q.tasks[1:]
	return task
}
