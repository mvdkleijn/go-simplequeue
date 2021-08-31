# go-simplequeue

Simple locking queue system with workers.

## Usage

```golang
// Define a job that conforms to the simplequeue.Job interface
type MyJob struct {
    id int
}

func (mj *MyJob) ID() int64 {
    return int64(mj.id)
}

func (mj *MyJob) Do() {
    // Lets just pause the job for a little time
    ms := time.Duration(rand.Intn(1000)+1) * time.Millisecond
    time.Sleep(ms)
    fmt.Printf("Job %d executing\n", mj.ID())
}

// Create some jobs for our test
func createJobs(number int) []*MyJob {
    jobs := make([]*MyJob, 0)

    for i := 1; i <= number; i++ {
        jobs = append(jobs, &MyJob{id: i})
    }

    return jobs
}

// Run our program
func main() {
    ctx := context.Background()

    // How much we want of each
    numWorkers := 15
    numJobs := 200

    // Create some jobs with a helper function
    jobs := createJobs(numJobs)

    // Create a queue
    q := sq.CreateQueue(ctx)

    // Initialize the workers
    workers := sq.InitializeWorkers(ctx, numWorkers)

    fmt.Printf("Number of workers in pool: %d\n", len(workers))
    fmt.Printf("Number of jobs for queue: %d\n", len(jobs))

    // Push the jobs onto the Queue
    for _, job := range jobs {
        q.Push(job)
    }

    // Process the queue with some workers
    q.Process(ctx, workers)

    // Show some stats afterwards
    var totalJobsHandled int64 = 0
    for _, w := range workers {
        totalJobsHandled += w.Handled()

        fmt.Printf("Worker %d processed a total of %d jobs\n", w.ID(), w.Handled())
    }

    fmt.Printf("Total jobs handled: %d\n", totalJobsHandled)
    fmt.Printf("Total workers: %d\n", len(workers))
}
```
