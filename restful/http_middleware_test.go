package restful

import (
	"bytes"
	"github.com/PCloud63514/goat/logger"
	"github.com/PCloud63514/goat/tid"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// 패닉이 발생한 요청을 500 응답으로 전환한다.
func TestMiddlewareRecoversPanic(t *testing.T) {
	logger.SetLogger(logger.NewSlogLogger(httptest.NewRecorder()))

	handler := Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("boom")
	}))

	request := httptest.NewRequest(http.MethodGet, "/panic", nil)
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, recorder.Code)
	}

	if recorder.Header().Get(tidHeaderName) == "" {
		t.Fatal("expected tid header to be set")
	}
}

// 상태 코드와 메시지를 함께 기록하고 에러 응답을 작성한다.
func TestErrorHandlerWritesExpectedResponse(t *testing.T) {
	logger.SetLogger(logger.NewSlogLogger(httptest.NewRecorder()))

	request := httptest.NewRequest(http.MethodGet, "/health", nil)
	recorder := httptest.NewRecorder()

	ErrorHandler(recorder, request, http.StatusMethodNotAllowed, "method not allowed")

	if recorder.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status %d, got %d", http.StatusMethodNotAllowed, recorder.Code)
	}
}

// 요청 로그를 한 줄 텍스트 포맷으로 출력하고 tid를 함께 기록한다.
func TestMiddlewareLogsTextWithTID(t *testing.T) {
	logBuffer := &bytes.Buffer{}
	logger.SetLogger(logger.NewSlogLogger(logBuffer))
	logger.Now = func() time.Time {
		return time.Date(2026, 4, 9, 14, 41, 41, 0, time.Local)
	}
	defer func() {
		logger.Now = time.Now
	}()

	handler := Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if tid.TIDFromContext(r.Context()) == "" {
			t.Fatal("expected tid in context")
		}
		w.WriteHeader(http.StatusAccepted)
	}))

	request := httptest.NewRequest(http.MethodGet, "/health", nil)
	request.Header.Set(tidHeaderName, "tid-from-client")
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if recorder.Header().Get(tidHeaderName) != "tid-from-client" {
		t.Fatalf("expected tid header %q, got %q", "tid-from-client", recorder.Header().Get(tidHeaderName))
	}

	logOutput := logBuffer.String()
	if !strings.Contains(logOutput, "INFO ") || !strings.Contains(logOutput, "14:41:41.000") {
		t.Fatalf("expected text log output with level and timestamp, got %s", logOutput)
	}
	if !strings.Contains(logOutput, "tid-from-client") || !strings.Contains(logOutput, "http request completed") {
		t.Fatalf("expected text log output, got %s", logOutput)
	}
	if !strings.Contains(logOutput, "status:202") {
		t.Fatalf("expected log output to contain status, got %s", logOutput)
	}
}
