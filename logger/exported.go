package logger

import (
	"context"
	"log"
	"strings"
)

var logger Logger = newDefaultLogger()

func init() {
	if logger == nil {
		log.Fatal("logger must not be nil")
	}
}

func Trace(args ...interface{}) {
	logger.Trace(args...)
}

func Tracef(format string, args ...interface{}) {
	logger.Tracef(format, args...)
}

func Fatal(args ...interface{}) {
	logger.Fatal(args...)
}

func Fatalf(format string, args ...interface{}) {
	logger.Fatalf(format, args...)
}

func Error(args ...interface{}) {
	logger.Error(args...)
}

func Errorf(format string, args ...interface{}) {
	logger.Errorf(format, args...)
}

func Debug(args ...interface{}) {
	logger.Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	logger.Debugf(format, args...)
}

func Warn(args ...interface{}) {
	logger.Warn(args...)
}

func Warnf(format string, args ...interface{}) {
	logger.Warnf(format, args...)
}

func Info(args ...interface{}) {
	logger.Info(args...)
}

func Infof(format string, args ...interface{}) {
	logger.Infof(format, args...)
}

func Panic(args ...interface{}) {
	logger.Panic(args...)
}

func Panicf(format string, args ...interface{}) {
	logger.Panicf(format, args...)
}

func ToLevel(s string) Level {
	_s := strings.ToLower(s)
	if level, ok := loggerLevelMap[_s]; ok {
		return level
	}
	return InfoLevel
}

func Initialize(cfg *Config) {
	logger.SetLogLevel(ToLevel(cfg.Level))
}

func SetLogger(l Logger) {
	logger = l
}

func WithFields(fields ...Field) Logger {
	return logger.WithFields(fields...)
}

func WithContext(ctx context.Context) Logger {
	return logger.WithFields(contextFields(ctx)...)
}

func StdLogger(level Level) *log.Logger {
	return log.New(&writerLogger{level: level}, "", 0)
}
