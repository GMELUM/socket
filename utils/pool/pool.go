package pool

import (
	"fmt"
	"sync"
)

type Pool struct {
	sem   chan struct{}
	work  chan func()
	size  int
	queue int
	mutex sync.Mutex
}

func New(size, queue int) *Pool {
	pool := &Pool{
		sem:   make(chan struct{}, size),
		work:  make(chan func(), queue+size),
		size:  size,
		queue: queue,
	}

	return pool
}

func (p *Pool) Schedule(task func()) error {

	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.queue != 0 && len(p.work) >= (p.queue+p.size) {
		return fmt.Errorf("queue is full")
	}

	select {
	case p.work <- task:
	case p.sem <- struct{}{}:
		go p.worker(task)
	}

	return nil
}

func (p *Pool) Ready() bool {
	return p.queue == 0 || len(p.work) < p.queue+p.size
}

func (p *Pool) NotReady() bool {
	return !p.Ready()
}

func (p *Pool) worker(task func()) {
	defer func() { <-p.sem }()
	task()
	for task := range p.work {
		task()
	}
}
