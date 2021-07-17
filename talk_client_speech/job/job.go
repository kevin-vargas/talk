package job

import (
	"fmt"
	"time"
)

const (
	MAX_JOB_QUANTITY = 100
	TIMEOUT          = 10
)

type WorkManager struct {
	jobs     chan []*Job
	jobBatch chan *Job
}

func New() *WorkManager {
	jobBatch := make(chan *Job, MAX_JOB_QUANTITY)
	jobs := make(chan []*Job, MAX_JOB_QUANTITY)
	go func() {
		for {
			sliceJob := append([]*Job{}, <-jobBatch)
		L:
			for {
				select {
				case job := <-jobBatch:
					sliceJob = append(sliceJob, job)
				case <-time.After(TIMEOUT * time.Second):
					fmt.Println("ADDING TIMEOUT")
					jobs <- sliceJob
					break L
				}

			}

		}
	}()
	return &WorkManager{
		jobs,
		jobBatch,
	}
}
func (wm *WorkManager) AddJob(fn func()) {
	job := Job{time.Now(), fn}
	wm.jobBatch <- &job
}

func (wm *WorkManager) ExecuteJob() {
	fmt.Println("EXECUTING JOBS")
	batchJobs := <-wm.jobs
	fmt.Println("QUANTIY JOBS ", len(batchJobs))
	syncChan := make(chan bool, 1)
	defer close(syncChan)
	for i := 0; i < len(batchJobs); i++ {
		syncChan <- true
		fmt.Println("EXECUTING JOB ", i+1)
		batchJobs[i].Execute()
		fmt.Println("FINISH JOB ", i+1)
		<-syncChan
	}
	fmt.Println("FINISH JOBS")
}

type Job struct {
	CreatedAt time.Time
	Execute   func()
}
