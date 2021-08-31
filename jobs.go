package simplequeue

type Job interface {
	ID() int64 // Returns the job's ID
	Do()       // Actually executes the job
}
