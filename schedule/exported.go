package schedule

func NewScheduler() Scheduler {
	return NewCronScheduler()
}

func RegisterJobs(scheduler Scheduler, jobs ...Job) {
	for _, job := range jobs {
		scheduler.RegisterJob(job)
	}
}
