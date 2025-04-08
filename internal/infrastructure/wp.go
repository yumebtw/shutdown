package infrastructure

import (
	"context"
	"fmt"
	"sync"
)

type Task func()

type WorkerPool struct {
	wg  sync.WaitGroup
	sem chan struct{}
	ctx context.Context
}

func NewWorkerPool(limit int, ctx context.Context) *WorkerPool {
	return &WorkerPool{
		wg:  sync.WaitGroup{},
		sem: make(chan struct{}, limit),
		ctx: ctx,
	}
}

func (p *WorkerPool) AddJob(task Task) {
	select {
	case <-p.ctx.Done():
	case p.sem <- struct{}{}:
		p.wg.Add(1)
		go func() {
			defer func() {
				<-p.sem
				p.wg.Done()
			}()
			task()
		}()
	default:
		fmt.Println("Server is currently busy")
	}
}

func (p *WorkerPool) Wait() {
	p.wg.Wait()
}
