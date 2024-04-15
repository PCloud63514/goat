package schedule

import (
	"context"
	"github.com/PCloud63514/goat"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"sync"
)

var (
	scheduler *cron.Cron
)

func init() {
	scheduler = cron.New(cron.WithSeconds())
	goat.AddHandlerFunc(start, goat.HandlerType_Starting)
	goat.AddHandlerFunc(stop, goat.HandlerType_Stop)
}

func start(ctx context.Context, env *goat.Environment) {
	scheduler.Start()
	logrus.Infof("Shcedule start")
}

func stop(ctx context.Context, env *goat.Environment) {
	c := scheduler.Stop()
	<-c.Done()
	logrus.Infof("Shcedule Stop Complate")
}

func AddJob(spec string, job Job) {
	scheduler.AddJob(spec, &MutexJobWrapper{
		job:   job,
		mutex: sync.Mutex{},
	})
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
