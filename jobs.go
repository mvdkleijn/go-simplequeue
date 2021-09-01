package simplequeue

// Job should be implemented by any struct intended for queueing
type Job interface {
	ID() int64 // Returns the job's ID
	Do()       // Actually executes the job
}
