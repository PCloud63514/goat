package schedule

import (
	"context"
	"github.com/sirupsen/logrus"
	"polygon-collector/internal/application"
	"sync"
)

var (
	cronInstance *cron.Cron
)

func init() {
	cronInstance = cron.New(cron.WithSeconds())
	application.AddStartEventCallback(start)
	application.AddDestroyEventCallback(stop)
}

func AddJob(spec string, job Job) {
	cronInstance.AddJob(spec, newMutexJobWrapper(job))
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

func start() {
	cronInstance.Start()
	logrus.Infof("Cron Shcedule Start")
}

func stop(ctx context.Context) {
	c := cronInstance.Stop()
	<-c.Done()
	logrus.Infof("Cron Shcedule Stop")
}
