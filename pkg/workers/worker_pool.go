package workers

import (
	"log"
	"sync"
)

type Task struct {
	Execute func()
}

type WorkerPool struct {
	workerCount int
	jobQueue    chan Task
	wg          sync.WaitGroup
}

func NewWorkerPool(workerCount int, jobQueueSize int) *WorkerPool {
	return &WorkerPool{
		workerCount: workerCount,
		jobQueue:    make(chan Task, jobQueueSize),
	}
}

func (wp *WorkerPool) Start() {
	for i := range wp.workerCount {
		go wp.worker(i)
	}
}

func (wp *WorkerPool) worker(id int) {
	for task := range wp.jobQueue {
		wp.wg.Add(1)
		func() {
			defer wp.wg.Done()
			defer func() {
				if r := recover(); r != nil {
					log.Printf("Worker %d recovered from panic: %v", id, r)
				}
			}()
			task.Execute()
		}()
	}
}

func (wp *WorkerPool) AddTask(task Task) {
	wp.jobQueue <- task
}

func (wp *WorkerPool) Stop() {
	close(wp.jobQueue)
	wp.wg.Wait()
}
