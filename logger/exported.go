package logger

import "strings"

func Fatal(args ...interface{}) {
	logger.Fatal(args)
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

func SetLogger(l Logger) {
	logger = l
}
