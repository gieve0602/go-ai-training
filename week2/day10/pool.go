package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

const DefaultWorkerSize = 3

var (
	ErrReceiverNil        = errors.New("method receiver cannot be nil")
	ErrTaskListEmpty      = errors.New("task list cannot be empty")
	ErrPoolAlreadyRunning = errors.New("pool already running")
	ErrPoolTerminated     = errors.New("pool already terminated")
	ErrPoolShutdown       = errors.New("pool already shutdown")
	ErrPollUnknownStatus  = errors.New("pool status unknown")
)

type PoolStatus int

const (
	Init PoolStatus = iota // 标明pool的初始状态
	Running
	Draining
	Terminated
)

type Pool struct {
	workerSize   int
	tasks        []Task
	jobs         chan Task
	results      chan Result
	status       PoolStatus
	mux          sync.Mutex
	ctx          context.Context
	cancel       context.CancelFunc
	shutdownDone chan struct{}
	shutdownOnce sync.Once
	pwg          sync.WaitGroup
	wwg          sync.WaitGroup
	cwg          sync.WaitGroup
}

func NewPool(workerSize int, tasks []Task) (*Pool, error) {
	if workerSize <= 0 {
		workerSize = DefaultWorkerSize
	}
	if len(tasks) == 0 {
		return nil, ErrTaskListEmpty
	}
	copyTasks := make([]Task, len(tasks))
	copy(copyTasks, tasks)
	ctx, cancel := context.WithCancel(context.Background())
	return &Pool{
		workerSize:   workerSize,
		tasks:        copyTasks,
		jobs:         make(chan Task, workerSize),
		results:      make(chan Result, workerSize),
		status:       Init,
		ctx:          ctx,
		cancel:       cancel,
		shutdownDone: make(chan struct{}),
	}, nil
}

func (p *Pool) shutdown() {
	p.pwg.Wait()
	close(p.jobs)
	p.wwg.Wait()
	close(p.results)
	p.cwg.Wait()
	p.cancel()
	p.mux.Lock()
	p.status = Terminated
	p.mux.Unlock()
	close(p.shutdownDone)
}

func (p *Pool) Shutdown(ctx context.Context) error {
	if p == nil {
		return ErrReceiverNil
	}
	if ctx == nil {
		return ErrContextNil
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	p.mux.Lock()
	switch p.status {
	case Init:
		p.status = Terminated
		p.mux.Unlock()
		p.shutdownOnce.Do(func() {
			close(p.shutdownDone)
			p.cancel()
		})
		return nil
	case Running:
		p.status = Draining
		p.mux.Unlock()
		p.shutdownOnce.Do(func() {
			go p.shutdown()
		})
	case Draining:
		p.mux.Unlock()
	case Terminated:
		p.mux.Unlock()
	default:
		p.mux.Unlock()
		return ErrPollUnknownStatus
	}

	select {
	case <-p.shutdownDone:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (p *Pool) producer() {
	defer p.pwg.Done()
	for _, task := range p.tasks {
		select {
		case <-p.ctx.Done():
			return
		case p.jobs <- task:
		}
	}
}

func (p *Pool) worker(id int) {
	defer p.wwg.Done()
	for {
		select {
		case <-p.ctx.Done():
			fmt.Printf("Worker-%d: exit with cancel\n", id)
			return
		case job, ok := <-p.jobs:
			if !ok {
				fmt.Printf("Worker-%d: exit with closed jobs queue\n", id)
				return
			}
			result := ExecuteTask(p.ctx, &job)
			p.results <- result
		}
	}
}

func (p *Pool) collector() {
	defer p.cwg.Done()
	for result := range p.results {
		if result.Err != nil {
			fmt.Printf("Task-%d: run error %v", result.TaskID, result.Err)
		} else {
			fmt.Printf("%v", result.Output)
		}
	}
}

func (p *Pool) Run() error {
	if p == nil {
		return ErrReceiverNil
	}
	p.mux.Lock()
	defer p.mux.Unlock()
	switch p.status {
	case Init:
		p.status = Running
		p.cwg.Add(1)
		go p.collector()

		p.wwg.Add(p.workerSize)
		for i := 1; i <= p.workerSize; i++ {
			go p.worker(i)
		}

		p.pwg.Add(1)
		go p.producer()
		return nil
	case Running:
		return ErrPoolAlreadyRunning
	case Draining:
		return ErrPoolShutdown
	case Terminated:
		return ErrPoolTerminated
	default:
		return ErrPollUnknownStatus
	}
}
