package utils

import (
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

// JobTask job object to be executed by pool
type JobTask interface {
	Execute() error
}

type JobReporter interface {
	Done() error
	Close() error
}

// Pool worker pool of specified size to process JobTask
type Pool struct {
	size      int
	in        chan JobTask
	wg        sync.WaitGroup
	reporter  JobReporter
	mx        sync.RWMutex
	queued    int
	processed int
	running   bool
}

// NewPool create a worker pool of given size
func NewPool(size int) *Pool {
	return &Pool{
		size: size,
		in:   make(chan JobTask),
		wg:   sync.WaitGroup{},
	}
}

func (p *Pool) SetJobReporter(reporter JobReporter) {
	p.reporter = reporter
}

// Start the worker pool
func (p *Pool) Start() {
	p.running = true
	go p.report()

	for w := 1; w <= p.size; w++ {
		p.wg.Add(1)
		go p.process()
	}

}

func (p *Pool) report() {
	for p.running {
		p.mx.RLock()
		log.Infof("processed %d of %d queued tasks (%0.4f%%)", p.processed, p.queued, float64(p.processed)/float64(p.queued)*100)
		p.mx.RUnlock()

		time.Sleep(10 * time.Second)
	}

	log.Info("ended reporting")
}

// Add another task to the pool
func (p *Pool) Add(t JobTask) {
	p.mx.Lock()
	p.queued += 1
	p.mx.Unlock()

	p.in <- t
}

func (p *Pool) AddBulk(tasks []JobTask) {
	p.mx.Lock()
	p.queued += len(tasks)
	p.mx.Unlock()

	for _, t := range tasks {
		p.in <- t
	}
}

func (p *Pool) process() {
	defer p.wg.Done()

	for task := range p.in {
		err := task.Execute()

		if err != nil {
			log.Error(err)
		}

		p.done()
	}
}

// Close the pool, wait for WaitGroup to complete
func (p *Pool) Close() {
	close(p.in)
	p.wg.Wait()
	if p.reporter != nil {
		if err := p.reporter.Close(); err != nil {
			log.Error(err)
		}
	}
	p.mx.Lock()
	defer p.mx.Unlock()
	p.running = false
}

func (p *Pool) done() {
	p.mx.Lock()
	defer p.mx.Unlock()

	p.processed += 1

	if p.reporter != nil {
		if err := p.reporter.Done(); err != nil {
			log.Error(err)
		}
	}
}
