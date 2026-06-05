package logger

import (
	"context"
	"sync"
)

type ContextMiddleware interface {
	Fields(ctx context.Context) []Field
}

var (
	contextMiddlewares   []ContextMiddleware
	contextMiddlewaresMu sync.RWMutex
)

func RegisterMiddleware(middlewares ...ContextMiddleware) {
	contextMiddlewaresMu.Lock()
	defer contextMiddlewaresMu.Unlock()

	contextMiddlewares = append(contextMiddlewares, middlewares...)
}

func contextFields(ctx context.Context) []Field {
	contextMiddlewaresMu.RLock()
	defer contextMiddlewaresMu.RUnlock()

	fields := make([]Field, 0)
	for _, middleware := range contextMiddlewares {
		fields = append(fields, middleware.Fields(ctx)...)
	}

	return fields
}
