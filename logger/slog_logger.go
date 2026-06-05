package logger

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const defaultLogLevel string = "debug"
const minBodyFieldWidth = 5
const levelWidth = 5
const tidWidth = 5
const jobNameWidth = 18

type slogLogger struct {
	writer      io.Writer
	level       Level
	fields      []Field
	mutex       *sync.Mutex
	fieldWidths *sync.Map
}

func newDefaultLogger() Logger {
	l := NewSlogLogger(os.Stdout)
	if envLevel := os.Getenv(logger_level_env_name); envLevel != "" {
		l.SetLogLevel(convertLogLevel(envLevel))
	}
	return l
}

func NewSlogLogger(writer io.Writer) Logger {
	return &slogLogger{
		writer:      writer,
		level:       convertLogLevel(defaultLogLevel),
		fields:      nil,
		mutex:       &sync.Mutex{},
		fieldWidths: &sync.Map{},
	}
}

func (l *slogLogger) Trace(args ...interface{}) {
	l.log(TraceLevel, fmt.Sprint(args...))
}

func (l *slogLogger) Tracef(format string, args ...interface{}) {
	l.log(TraceLevel, fmt.Sprintf(format, args...))
}

func (l *slogLogger) Debug(args ...interface{}) {
	l.log(DebugLevel, fmt.Sprint(args...))
}

func (l *slogLogger) Debugf(format string, args ...interface{}) {
	l.log(DebugLevel, fmt.Sprintf(format, args...))
}

func (l *slogLogger) Info(args ...interface{}) {
	l.log(InfoLevel, fmt.Sprint(args...))
}

func (l *slogLogger) Infof(format string, args ...interface{}) {
	l.log(InfoLevel, fmt.Sprintf(format, args...))
}

func (l *slogLogger) Warn(args ...interface{}) {
	l.log(WarnLevel, fmt.Sprint(args...))
}

func (l *slogLogger) Warnf(format string, args ...interface{}) {
	l.log(WarnLevel, fmt.Sprintf(format, args...))
}

func (l *slogLogger) Error(args ...interface{}) {
	l.log(ErrorLevel, fmt.Sprint(args...))
}

func (l *slogLogger) Errorf(format string, args ...interface{}) {
	l.log(ErrorLevel, fmt.Sprintf(format, args...))
}

func (l *slogLogger) Panic(args ...interface{}) {
	msg := fmt.Sprint(args...)
	l.log(PanicLevel, msg)
	panic(msg)
}

func (l *slogLogger) Panicf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	l.log(PanicLevel, msg)
	panic(msg)
}

func (l *slogLogger) Fatal(args ...interface{}) {
	msg := fmt.Sprint(args...)
	l.log(FatalLevel, msg)
	os.Exit(1)
}

func (l *slogLogger) Fatalf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	l.log(FatalLevel, msg)
	os.Exit(1)
}

func (l *slogLogger) SetLogLevel(level Level) {
	l.level = level
}

func (l *slogLogger) WithFields(fields ...Field) Logger {
	combinedFields := make([]Field, 0, len(l.fields)+len(fields))
	combinedFields = append(combinedFields, l.fields...)
	combinedFields = append(combinedFields, fields...)

	return &slogLogger{
		writer:      l.writer,
		level:       l.level,
		fields:      combinedFields,
		mutex:       l.mutex,
		fieldWidths: l.fieldWidths,
	}
}

func (l *slogLogger) log(level Level, message string) {
	if level < l.level {
		return
	}

	l.mutex.Lock()
	defer l.mutex.Unlock()

	line := buildLogLine(level, message, l.fields, l.fieldWidths)
	_, _ = io.WriteString(l.writer, line)
}

type writerLogger struct {
	level Level
}

func (w *writerLogger) Write(p []byte) (int, error) {
	message := string(p)
	switch w.level {
	case TraceLevel:
		Trace(message)
	case DebugLevel:
		Debug(message)
	case InfoLevel:
		Info(message)
	case WarnLevel:
		Warn(message)
	default:
		Error(message)
	}
	return len(p), nil
}

func buildLogLine(level Level, message string, fields []Field, fieldWidths *sync.Map) string {
	var builder strings.Builder
	headerFields := make([]Field, 0, len(fields))
	bodyFields := make([]Field, 0, len(fields))

	for _, field := range fields {
		if field.Placement == FieldPlacementHeader {
			headerFields = append(headerFields, field)
			continue
		}
		bodyFields = append(bodyFields, field)
	}

	builder.WriteString(fmt.Sprintf("%-*s", levelWidth, level.String()))
	builder.WriteByte(' ')
	builder.WriteString(Now().Format(loggerTimeLayout))
	for index, field := range headerFields {
		builder.WriteByte(' ')
		switch index {
		case 0:
			builder.WriteString(fmt.Sprintf("%-*s", tidWidth, fmt.Sprint(field.Value)))
		case 1:
			builder.WriteString(fmt.Sprintf("%-*s", jobNameWidth, fmt.Sprint(field.Value)))
		default:
			builder.WriteString(fmt.Sprint(field.Value))
		}
	}
	builder.WriteByte('\t')
	builder.WriteString(message)

	for _, field := range bodyFields {
		builder.WriteByte(' ')
		builder.WriteString(formatBodyField(field, fieldWidths))
	}

	builder.WriteByte('\n')
	return builder.String()
}

func (l Level) String() string {
	switch l {
	case TraceLevel:
		return "TRACE"
	case DebugLevel:
		return "DEBUG"
	case InfoLevel:
		return "INFO"
	case WarnLevel:
		return "WARN"
	case ErrorLevel:
		return "ERROR"
	case PanicLevel:
		return "PANIC"
	case FatalLevel:
		return "FATAL"
	default:
		return "INFO"
	}
}

var Now = time.Now

func formatBodyField(field Field, fieldWidths *sync.Map) string {
	value := fmt.Sprint(field.Value)
	if !field.EnableMinWidth {
		return fmt.Sprintf("%s:%s", field.Key, value)
	}

	width := loadFieldWidth(fieldWidths, field.Key)
	if len(value) > width {
		updateFieldWidth(fieldWidths, field.Key, len(value))
		width = loadFieldWidth(fieldWidths, field.Key)
	}

	return fmt.Sprintf("%s:%-*s", field.Key, width, value)
}

func loadFieldWidth(fieldWidths *sync.Map, key string) int {
	if fieldWidths == nil {
		return minBodyFieldWidth
	}

	value, ok := fieldWidths.Load(key)
	if !ok {
		return minBodyFieldWidth
	}

	width, ok := value.(*atomic.Int64)
	if !ok {
		return minBodyFieldWidth
	}

	current := int(width.Load())
	if current < minBodyFieldWidth {
		return minBodyFieldWidth
	}

	return current
}

func updateFieldWidth(fieldWidths *sync.Map, key string, candidate int) {
	if fieldWidths == nil {
		return
	}

	value, _ := fieldWidths.LoadOrStore(key, &atomic.Int64{})
	width, ok := value.(*atomic.Int64)
	if !ok {
		return
	}

	if width.Load() < minBodyFieldWidth {
		width.CompareAndSwap(width.Load(), minBodyFieldWidth)
	}

	for {
		current := width.Load()
		if current >= int64(candidate) {
			return
		}
		if width.CompareAndSwap(current, int64(candidate)) {
			return
		}
	}
}
