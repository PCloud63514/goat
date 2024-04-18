package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var (
	encoderConfig zapcore.EncoderConfig = makeDefaultEncoderConfig()
	encoder       zapcore.Encoder       = zapcore.NewConsoleEncoder(encoderConfig)
)

func init() {
	logger = newZapLogger()
}

type zapLogger struct {
	sugaredLogger *zap.SugaredLogger
}

func newZapLogger() Logger {
	return &zapLogger{
		sugaredLogger: makeSugaredLogger(convertZapLogLevel(DebugLevel)),
	}
}

func (l *zapLogger) Fatal(args ...interface{}) {
	l.sugaredLogger.Fatal(args...)
}

func (l *zapLogger) Fatalf(format string, args ...interface{}) {
	l.sugaredLogger.Fatalf(format, args...)
}

func (l *zapLogger) Error(args ...interface{}) {
	l.sugaredLogger.Error(args...)
}

func (l *zapLogger) Errorf(format string, args ...interface{}) {
	l.sugaredLogger.Errorf(format, args...)
}

func (l *zapLogger) Debug(args ...interface{}) {
	l.sugaredLogger.Debug(args...)
}

func (l *zapLogger) Debugf(format string, args ...interface{}) {
	l.sugaredLogger.Debugf(format, args...)
}

func (l *zapLogger) Warn(args ...interface{}) {
	l.sugaredLogger.Warn(args...)
}

func (l *zapLogger) Warnf(format string, args ...interface{}) {
	l.sugaredLogger.Warnf(format, args...)
}

func (l *zapLogger) Info(args ...interface{}) {
	l.sugaredLogger.Info(args...)
}

func (l *zapLogger) Infof(format string, args ...interface{}) {
	l.sugaredLogger.Infof(format, args...)
}

func (l *zapLogger) Panic(args ...interface{}) {
	l.sugaredLogger.Panic(args...)
}

func (l *zapLogger) Panicf(format string, args ...interface{}) {
	l.sugaredLogger.Panicf(format, args...)
}

func (l *zapLogger) SetLogLevel(level Level) {
	l.sugaredLogger = makeSugaredLogger(convertZapLogLevel(level))
}

func convertZapLogLevel(level Level) zapcore.Level {
	switch level {
	case DebugLevel:
		return zapcore.DebugLevel
	case InfoLevel:
		return zapcore.InfoLevel
	case WarnLevel:
		return zapcore.WarnLevel
	case ErrorLevel:
		return zapcore.ErrorLevel
	case PanicLevel:
		return zapcore.PanicLevel
	case FatalLevel:
		return zapcore.FatalLevel
	}
	return zapcore.DebugLevel
}

func makeDefaultEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "M",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout(LOGGER_TIME_LAYOUT),
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

func makeSugaredLogger(l zapcore.Level) *zap.SugaredLogger {
	logWriter := zapcore.AddSync(os.Stderr)
	core := zapcore.NewCore(encoder, logWriter, zap.NewAtomicLevelAt(l))
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	defer logger.Sync()
	sugaredLogger := logger.Sugar()
	return sugaredLogger
}
