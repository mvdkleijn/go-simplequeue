package simplequeue

import (
	"context"
	"fmt"
	"runtime/trace"
)

type WorkerI interface {
	ID() int64
	Handled() int64
	Process(job Job)
}

type Worker struct {
	id          int64
	jobsHandled int64
}

// ID returns the worker's ID
func (w *Worker) ID() int64 {
	return w.id
}

// Handled returns the number of jobs the worker handled
func (w *Worker) Handled() int64 {
	return w.jobsHandled
}

func (w *Worker) process(ctx context.Context, job Job) {
	job.Do()
	w.jobsHandled++
	trace.Log(ctx, "process", fmt.Sprintf("worker %v finished Job %v.", w.ID(), job.ID()))
}

// Initialize and return a pool of workers
func InitializeWorkers(ctx context.Context, num int) []*Worker {
	_, task := trace.NewTask(ctx, "initialize workers")
	workers := make([]*Worker, 0)

	for i := 1; i <= num; i++ {
		workers = append(workers, &Worker{id: int64(i)})
	}

	task.End()
	return workers
}
