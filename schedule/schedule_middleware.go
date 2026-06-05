package schedule

import (
	"context"
	"github.com/PCloud63514/goat/logger"
	"github.com/PCloud63514/goat/tid"
)

type ContextMiddleware struct{}

func NewContextMiddleware() *ContextMiddleware {
	return &ContextMiddleware{}
}

func (m *ContextMiddleware) Fields(ctx context.Context) []logger.Field {
	tidValue := tid.TIDFromContext(ctx)
	if tidValue == "" {
		tidValue = "-"
	}

	jobName := JobNameFromContext(ctx)
	if jobName == "" {
		jobName = "-"
	}

	return []logger.Field{
		logger.HeaderString("tid", tidValue),
		logger.HeaderString("jobName", jobName),
	}
}
