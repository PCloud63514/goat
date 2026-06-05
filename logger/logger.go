package logger

import (
	"encoding/json"
	"fmt"
	"time"
)

type Level int8

const (
	TraceLevel Level = iota - 2
	DebugLevel
	InfoLevel
	WarnLevel
	ErrorLevel
	PanicLevel
	FatalLevel
	_minLevel = TraceLevel
	_maxLevel = FatalLevel
)

const loggerTimeLayout string = "15:04:05.000"

const (
	logger_level_env_name          = "LOG_LEVEL"
	logger_level_env_default_value = "debug"
)

var (
	loggerLevelMap = map[string]Level{
		"trace": TraceLevel,
		"debug": DebugLevel,
		"info":  InfoLevel,
		"warn":  WarnLevel,
		"error": ErrorLevel,
		"panic": PanicLevel,
		"fatal": FatalLevel,
		"TRACE": TraceLevel,
		"DEBUG": DebugLevel,
		"INFO":  InfoLevel,
		"WARN":  WarnLevel,
		"ERROR": ErrorLevel,
		"PANIC": PanicLevel,
		"FATAL": FatalLevel,
	}
)

type Logger interface {
	Trace(args ...interface{})
	Tracef(format string, args ...interface{})
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
	WithFields(fields ...Field) Logger
}

type FieldPlacement string

const (
	FieldPlacementBody   FieldPlacement = "body"
	FieldPlacementHeader FieldPlacement = "header"
)

type Field struct {
	Key            string
	Value          interface{}
	Placement      FieldPlacement
	EnableMinWidth bool
}

const (
	nullFieldValue   = "null"
	emptyHeaderValue = "-"
)

func String(key string, value string) Field {
	return Field{Key: key, Value: value, Placement: FieldPlacementBody, EnableMinWidth: true}
}

func PString(key string, value *string) Field {
	if value == nil {
		return String(key, nullFieldValue)
	}

	return String(key, *value)
}

func HeaderString(key string, value string) Field {
	return Field{Key: key, Value: value, Placement: FieldPlacementHeader}
}

func PHeaderString(key string, value *string) Field {
	if value == nil {
		return HeaderString(key, emptyHeaderValue)
	}

	return HeaderString(key, *value)
}

func Int(key string, value int) Field {
	return Field{Key: key, Value: value, Placement: FieldPlacementBody, EnableMinWidth: true}
}

func PInt(key string, value *int) Field {
	if value == nil {
		return String(key, nullFieldValue)
	}

	return Int(key, *value)
}

func Uint64(key string, value uint64) Field {
	return Field{Key: key, Value: value, Placement: FieldPlacementBody, EnableMinWidth: true}
}

func PUint64(key string, value *uint64) Field {
	if value == nil {
		return String(key, nullFieldValue)
	}

	return Uint64(key, *value)
}

func Duration(key string, value time.Duration) Field {
	return Field{Key: key, Value: value, Placement: FieldPlacementBody, EnableMinWidth: true}
}

func PDuration(key string, value *time.Duration) Field {
	if value == nil {
		return String(key, nullFieldValue)
	}

	return Duration(key, *value)
}

func Any(key string, value interface{}) Field {
	return Field{Key: key, Value: value, Placement: FieldPlacementBody}
}

func PAny(key string, value interface{}) Field {
	if value == nil {
		return String(key, nullFieldValue)
	}

	return Any(key, value)
}

func JSON(key string, value interface{}) Field {
	return Field{Key: key, Value: ToJSON(value), Placement: FieldPlacementBody}
}

func PJSON(key string, value interface{}) Field {
	if value == nil {
		return String(key, nullFieldValue)
	}

	return JSON(key, value)
}

func ToJSON(value interface{}) string {
	encoded, err := json.Marshal(value)
	if err != nil {
		return fmt.Sprintf("{\"marshalError\":%q}", err.Error())
	}
	return string(encoded)
}

func Bool(key string, value bool) Field {
	return Field{Key: key, Value: value, Placement: FieldPlacementBody, EnableMinWidth: true}
}

func PBool(key string, value *bool) Field {
	if value == nil {
		return String(key, nullFieldValue)
	}

	return Bool(key, *value)
}

func convertLogLevel(level string) Level {
	v, ok := loggerLevelMap[level]
	if !ok {
		return InfoLevel
	}
	return v
}
