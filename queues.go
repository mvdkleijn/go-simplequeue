/*
	Copyright (C) 2023  Martijn van der Kleijn

	This file is part of the go-simplequeue library.

    This Source Code Form is subject to the terms of the Mozilla Public
  	License, v. 2.0. If a copy of the MPL was not distributed with this
  	file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package simplequeue

import (
	"context"
	"fmt"
	"runtime/trace"
	"sync"
)

// QueueI is intended for possible future expansion
type QueueI interface {
	Pop() *Job                                   // Remove a job from the queue and return it
	Push(job *Job)                               // Add a job onto the queue
	Jobs() int                                   // Returns the number of jobs on the queue
	Process(ctx context.Context, pool []*Worker) // Process the queue with a pool of workers
}

// Queue is set of Jobs that can be worked on by one or more Workers
type Queue struct {
	jobs []*Job
	Lock *sync.Mutex
}

// Pop a job off of the queue and return it
func (q *Queue) Pop() *Job {
	if len(q.jobs) > 0 {
		job := q.jobs[0]
		q.jobs = q.jobs[1:]

		return job
	}

	return nil
}

// Push a job onto the queue
func (q *Queue) Push(job Job) {
	q.jobs = append(q.jobs, &job)
}

// Jobs returns the number of jobs on the queue
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

// CreateQueue initializes and returns an empty queue
func CreateQueue(ctx context.Context) *Queue {
	_, task := trace.NewTask(ctx, "initialize queue")
	q := &Queue{}
	q.jobs = []*Job{}
	q.Lock = &sync.Mutex{}
	task.End()
	return q
}
