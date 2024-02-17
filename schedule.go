package goat

import (
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"sync"
)

var (
	scheduler = cron.New(cron.WithSeconds())
)

func AddSchedule(spec string, job Job) {
	scheduler.AddJob(spec, newMutexJobWrapper(job))
}

type Job interface {
	Run()
}

type MutexJobWrapper struct {
	job   Job
	mutex sync.Mutex
}

func (wrapper *MutexJobWrapper) Run() {
	if wrapper.mutex.TryLock() {
		defer wrapper.mutex.Unlock()
		defer func() {
			if r := recover(); r != nil {
				logrus.Error(r)
			}
		}()
		wrapper.job.Run()
	}
}

func newMutexJobWrapper(job Job) *MutexJobWrapper {
	mutex := sync.Mutex{}
	return &MutexJobWrapper{
		job:   job,
		mutex: mutex,
	}
}

func StartScheduler() {
	scheduler.Start()
}

func StopScheduler() {
	c := scheduler.Stop()
	<-c.Done()
}
