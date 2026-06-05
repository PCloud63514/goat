package restful

import (
	"context"
	"fmt"
	"github.com/PCloud63514/goat/logger"
	tid "github.com/PCloud63514/goat/tid"
	"net/http"
	"runtime/debug"
	"time"
)

const (
	tidHeaderName       = "X-REQSEQ"
	legacyTidHeaderName = "X-TID"
)

type HTTPContextMiddleware struct{}

func NewHTTPContextMiddleware() *HTTPContextMiddleware {
	return &HTTPContextMiddleware{}
}

func (m *HTTPContextMiddleware) Fields(ctx context.Context) []logger.Field {
	value := tid.TIDFromContext(ctx)
	if value == "" {
		value = "-"
	}

	return []logger.Field{
		logger.HeaderString("tid", value),
	}
}

type responseRecorder struct {
	http.ResponseWriter
	statusCode   int
	bytesWritten int
	wroteHeader  bool
}

func newResponseRecorder(writer http.ResponseWriter) *responseRecorder {
	return &responseRecorder{
		ResponseWriter: writer,
		statusCode:     http.StatusOK,
	}
}

func (r *responseRecorder) WriteHeader(statusCode int) {
	r.statusCode = statusCode
	r.wroteHeader = true
	r.ResponseWriter.WriteHeader(statusCode)
}

func (r *responseRecorder) Write(body []byte) (int, error) {
	if !r.wroteHeader {
		r.WriteHeader(http.StatusOK)
	}
	size, err := r.ResponseWriter.Write(body)
	r.bytesWritten += size
	return size, err
}

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		startedAt := time.Now()
		recorder := newResponseRecorder(writer)
		request, tid := WithTID(request)
		recorder.Header().Set(tidHeaderName, tid)
		requestLogger := logger.WithFields(
			logger.HeaderString("tid", tid),
			logger.String("method", request.Method),
			logger.String("path", request.URL.Path),
			logger.String("host", request.Host),
			logger.String("remoteAddr", request.RemoteAddr),
			logger.String("userAgent", request.UserAgent()),
		)

		defer func() {
			if recovered := recover(); recovered != nil {
				recorder.statusCode = http.StatusInternalServerError
				requestLogger.WithFields(
					logger.Int("status", recorder.statusCode),
					logger.Any("panic", recovered),
					logger.String("stack", string(debug.Stack())),
				).Error("http request panic")

				if !recorder.wroteHeader {
					http.Error(recorder, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				}
			}

			requestLogger.WithFields(
				logger.Int("status", recorder.statusCode),
				logger.Duration("duration", time.Since(startedAt)),
				logger.Int("bytes", recorder.bytesWritten),
			).Info("http request completed")
		}()

		next.ServeHTTP(recorder, request)
	})
}

func ErrorHandler(writer http.ResponseWriter, request *http.Request, statusCode int, message string) {
	logger.WithFields(
		logger.HeaderString("tid", tid.TIDFromContext(request.Context())),
		logger.String("method", request.Method),
		logger.String("path", request.URL.Path),
		logger.Int("status", statusCode),
	).Warn(message)
	http.Error(writer, message, statusCode)
}

func RecoverAsError(err interface{}) error {
	return fmt.Errorf("%v", err)
}

func WithTID(request *http.Request) (*http.Request, string) {
	_tid := request.Header.Get(tidHeaderName)
	if _tid == "" {
		_tid = request.Header.Get(legacyTidHeaderName)
	}
	if _tid == "" {
		_tid = tid.NewTID()
	}
	return request.WithContext(context.WithValue(request.Context(), tid.TidContextKey, _tid)), _tid
}
