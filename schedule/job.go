package schedule

import (
	"context"
	"github.com/PCloud63514/goat/logger"
	"github.com/PCloud63514/goat/tid"
	"sync"
)

type contextKey string

const JobNameContextKey contextKey = "X-JOBNAME"

type Job interface {
	Process(ctx context.Context)
	Name() string
	Spec() string
}

type MutexJobWrapper struct {
	job              Job
	mutex            sync.Mutex
	schedulerCtxFunc func() context.Context
}

func (w *MutexJobWrapper) Run() {
	if w.mutex.TryLock() {
		defer w.mutex.Unlock()
		defer func() {
			if r := recover(); r != nil {
				logger.Error(r)
			}
		}()
		ctx, cancel := context.WithCancel(w.schedulerCtxFunc())
		defer cancel()

		ctx, _ = tid.WithContextTID(ctx)
		ctx = context.WithValue(ctx, JobNameContextKey, normalizeJobName(w.job.Name()))
		w.job.Process(ctx)
	}
}

func JobNameFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	if value, ok := ctx.Value(JobNameContextKey).(string); ok {
		return value
	}

	return ""
}

func normalizeJobName(value string) string {
	if value == "" {
		return "-"
	}

	return value
}
