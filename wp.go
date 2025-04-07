package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Job struct {
	ID      int
	Payload string
}

type WorkerPool struct {
	Count    int
	JobQueue chan Job
	wg       *sync.WaitGroup
	ctx      context.Context
}

func NewWorkerPool(ctx context.Context, n int) *WorkerPool {
	return &WorkerPool{
		Count:    n,
		JobQueue: make(chan Job, n),
		wg:       &sync.WaitGroup{},
		ctx:      ctx,
	}
}

func (wp *WorkerPool) Worker(n int) {
	defer wp.wg.Done()

	for {
		select {
		case <-wp.ctx.Done():
			return
		case job := <-wp.JobQueue:
			fmt.Println("worker", n, "doing job", job.ID, "with payload", job.Payload)
			time.Sleep(1 * time.Second)
		}
	}
}

func (wp *WorkerPool) Start() {
	for i := 0; i < wp.Count; i++ {
		wp.wg.Add(1)
		go wp.Worker(i)
	}
}

func (wp *WorkerPool) AddJob(job Job) {
	select {
	case wp.JobQueue <- job:
	default:
		fmt.Println("job queue is full, dropping job:", job.Payload)
	}
}

func (wp *WorkerPool) Wait() {
	wp.wg.Wait()
}
