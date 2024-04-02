package job

import (
	"context"
	"fmt"
	"sync"
)

type job interface {
	Work(ctx context.Context)
}
type WorkerPool struct {
	numWorkers int
	JobQueue   chan job
	results    chan int
	wg         sync.WaitGroup
}

func NewWorkerPool(numWorkers, jobQueueSize int) *WorkerPool {
	return &WorkerPool{
		numWorkers: numWorkers,
		JobQueue:   make(chan job, jobQueueSize),
		results:    make(chan int, jobQueueSize),
	}
}

func (wp *WorkerPool) Worker(id int) {
	defer wp.wg.Done()
	for job := range wp.JobQueue {
		job.Work(context.Background())
		wp.results <- id
	}
}

func (wp *WorkerPool) Start() {
	for i := 1; 1 <= wp.numWorkers; i++ {
		wp.wg.Add(1)
		go wp.Worker(i)
	}
}

func (wp *WorkerPool) CollectResults() {
	for result := range wp.results {
		fmt.Printf("Result received for job %d\n", result)
	}
}

func (wp *WorkerPool) Add(j job) {
	wp.JobQueue <- j
}

func (wp *WorkerPool) Wait() {
	wp.wg.Wait()
	close(wp.results)
}
