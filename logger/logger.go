package logger

type Level int8

const (
	DebugLevel Level = iota - 1
	InfoLevel
	WarnLevel
	ErrorLevel
	PanicLevel
	FatalLevel
	_minLevel = DebugLevel
	_maxLevel = FatalLevel
)

const (
	LOGGER_TIME_LAYOUT string = "15:04:05"
)

var (
	loggerLevelMap = map[string]Level{
		"debug": DebugLevel,
		"info":  InfoLevel,
		"warn":  WarnLevel,
		"error": ErrorLevel,
		"panic": PanicLevel,
		"fatal": FatalLevel,
	}
)

type Logger interface {
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Panic(args ...interface{})
	Panicf(format string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	SetLogLevel(level Level)
}

var logger Logger
