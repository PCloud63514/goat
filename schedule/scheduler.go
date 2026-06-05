package schedule

import (
	"context"
	"github.com/PCloud63514/goat/logger"
	"sync"

	"github.com/robfig/cron/v3"
)

type Scheduler interface {
	RegisterJob(job Job) error
	Context() context.Context
	Start()
	Stop()
}

type CronScheduler struct {
	cronInstance *cron.Cron
	ctx          context.Context
	cancel       context.CancelFunc
	running      bool
	runningMu    sync.Mutex
}

func NewCronScheduler() *CronScheduler {
	return &CronScheduler{
		cronInstance: cron.New(cron.WithSeconds()),
	}
}

func (s *CronScheduler) RegisterJob(job Job) error {
	logger.Infof("CronSceduler.RegisterJob, jobName: %s, jobSpec: %s", job.Name(), job.Spec())
	_, err := s.cronInstance.AddJob(job.Spec(), &MutexJobWrapper{
		schedulerCtxFunc: s.Context,
		job:              job,
	})
	return err
}

func (s *CronScheduler) Context() context.Context {
	return s.ctx
}

func (s *CronScheduler) Start() {
	s.runningMu.Lock()
	defer s.runningMu.Unlock()
	if s.running {
		return
	}
	s.running = true
	ctx, cancel := context.WithCancel(context.Background())
	s.ctx = ctx
	s.cancel = cancel
	s.cronInstance.Start()
}

func (s *CronScheduler) Stop() {
	s.runningMu.Lock()
	defer s.runningMu.Unlock()
	if !s.running {
		return
	}
	s.running = false
	if s.cancel != nil {
		s.cancel()
	}
	stopCtx := s.cronInstance.Stop()
	<-stopCtx.Done()
}
