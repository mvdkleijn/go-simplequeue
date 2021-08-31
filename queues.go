package simplequeue

import (
	"context"
	"fmt"
	"runtime/trace"
	"sync"
)

type QueueI interface {
	Pop() *Job                                   // Remove a job from the queue and return it
	Push(job *Job)                               // Add a job onto the queue
	Jobs() int                                   // Returns the number of jobs on the queue
	Process(ctx context.Context, pool []*Worker) // Process the queue with a pool of workers
}

type Queue struct {
	jobs []*Job
	Lock *sync.Mutex
}

// Remove a job from the queue and return it
func (q *Queue) Pop() *Job {
	if len(q.jobs) > 0 {
		job := q.jobs[0]
		q.jobs = q.jobs[1:]

		return job
	}

	return nil
}

// Add a job onto the queue
func (q *Queue) Push(job Job) {
	q.jobs = append(q.jobs, &job)
}

// Returns the number of jobs on the queue
func (q *Queue) Jobs() int {
	return len(q.jobs)
}

// Process the queue with a pool of workers
func (q *Queue) Process(ctx context.Context, pool []*Worker) {
	wg := &sync.WaitGroup{}

	trace.Log(ctx, "process", fmt.Sprintf("detected %d workers in pool", len(pool)))

	for _, worker := range pool {
		wg.Add(1)
		go process(ctx, q, worker, wg)
	}
	wg.Wait()
}

func process(ctx context.Context, q *Queue, w *Worker, wg *sync.WaitGroup) {
	defer wg.Done()
	var job *Job

	q.Lock.Lock()
	queuedJobs := q.Jobs()
	if queuedJobs > 0 {
		job = q.Pop()
	}
	q.Lock.Unlock()

	for queuedJobs > 0 {
		ctx, task := trace.NewTask(ctx, "process")
		trace.Log(ctx, "process", fmt.Sprintf("processing jobs using worker %d", w.ID()))
		w.process(ctx, *job)
		task.End()

		q.Lock.Lock()
		queuedJobs = q.Jobs()
		job = q.Pop()
		q.Lock.Unlock()
	}
}

// Create a queue
func CreateQueue(ctx context.Context) *Queue {
	_, task := trace.NewTask(ctx, "initialize queue")
	q := &Queue{}
	q.jobs = []*Job{}
	q.Lock = &sync.Mutex{}
	task.End()
	return q
}
