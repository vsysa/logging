package logruslog

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/vsysa/logging"
	"github.com/vsysa/logging/internal/helper"
	"go.opentelemetry.io/otel/trace"
	"os"
	"sync"
)

var levelMap = map[logging.Level]logrus.Level{
	logging.TraceLevel: logrus.TraceLevel,
	logging.DebugLevel: logrus.DebugLevel,
	logging.InfoLevel:  logrus.InfoLevel,
	logging.WarnLevel:  logrus.WarnLevel,
	logging.ErrorLevel: logrus.ErrorLevel,
	logging.FatalLevel: logrus.FatalLevel,
}

type LogrusLogger struct {
	logrus    *logrus.Logger
	contextMu sync.RWMutex
	context   map[string]interface{}
	timerMu   sync.RWMutex
}

func NewLogrusLogger() *LogrusLogger {
	l := logrus.New()
	l.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02T15:04:05.999",
		ForceColors:     true,
		FullTimestamp:   true,
	})

	l.SetOutput(&levelAwareWriter{
		InfoWriter:  os.Stdout,
		ErrorWriter: os.Stderr,
	})

	newLogger := &LogrusLogger{
		logrus:  l,
		context: make(map[string]interface{}),
	}

	return newLogger
}

func NewLogrusLoggerWithParams(params map[string]interface{}) *LogrusLogger {
	newLogger := NewLogrusLogger()
	newLogger.AddContexts(params)
	return newLogger
}

//	CONTEXT

// AddContext добавляет ключ и значение в LoggerContext.
func (r *LogrusLogger) AddContext(key string, value interface{}) logging.Logger {
	r.contextMu.Lock()
	defer r.contextMu.Unlock()
	r.context[key] = fmt.Sprintf("%v", value)
	return r
}

// AddContexts добавляет несколько ключей и значений в LoggerContext.
func (r *LogrusLogger) AddContexts(contexts map[string]interface{}) logging.Logger {
	r.contextMu.Lock()
	defer r.contextMu.Unlock()
	for key, value := range contexts {
		r.context[key] = fmt.Sprintf("%v", value)
	}
	return r
}

// DeleteContext удаляет значение по ключу из LoggerContext.
func (r *LogrusLogger) DeleteContext(key string) logging.Logger {
	r.contextMu.Lock()
	defer r.contextMu.Unlock()
	delete(r.context, key)
	return r
}

// GetAllContext возвращает все ключи и значения из LoggerContext.
func (r *LogrusLogger) GetAllContexts() map[string]interface{} {
	r.contextMu.RLock()
	defer r.contextMu.RUnlock()
	newContext := make(map[string]interface{}, len(r.context))
	for key, value := range r.context {
		newContext[key] = value
	}
	return newContext
}

//	LOGGING

func (r *LogrusLogger) SetLevel(level logging.Level) {
	if logrusLevel, ok := levelMap[level]; ok {
		r.logrus.SetLevel(logrusLevel)
	} else {
		r.Warn("Invalid logging level: %v", level)
		return
	}
}

func (r *LogrusLogger) Trace(message string, a ...any) {
	r.log(logrus.TraceLevel, fmt.Sprintf(message, a...))
}

func (r *LogrusLogger) Debug(message string, a ...any) {
	r.log(logrus.DebugLevel, fmt.Sprintf(message, a...))
}

func (r *LogrusLogger) Info(message string, a ...any) {
	r.log(logrus.InfoLevel, fmt.Sprintf(message, a...))
}

func (r *LogrusLogger) Warn(message string, a ...any) {
	r.log(logrus.WarnLevel, fmt.Sprintf(message, a...))
}

func (r *LogrusLogger) Error(message string, a ...any) {
	r.log(logrus.ErrorLevel, fmt.Sprintf(message, a...))
}

func (r *LogrusLogger) Fatal(message string, a ...any) {
	r.log(logrus.FatalLevel, fmt.Sprintf(message, a...))
}

func (r *LogrusLogger) ErrorCatch(err error, message string, a ...any) {
	msg := fmt.Sprintf(message, a...)
	if err != nil {
		msg += ": " + err.Error()
	}
	r.log(logrus.ErrorLevel, msg)
}

func (r *LogrusLogger) FatalCatch(err error, message string, a ...any) {
	msg := fmt.Sprintf(message, a...)
	if err != nil {
		msg += ": " + err.Error()
	}
	r.log(logrus.FatalLevel, msg)
	os.Exit(1)
}

// BASE

func (r *LogrusLogger) Clone() logging.Logger {
	return &LogrusLogger{
		logrus:  r.logrus,
		context: helper.CopyMapContext(r.context),
	}
}

func (r *LogrusLogger) SetCtx(ctx context.Context) logging.Logger {
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		traceID := span.SpanContext().TraceID().String()
		spanID := span.SpanContext().SpanID().String()

		r.AddContexts(map[string]interface{}{
			"traceID": traceID,
			"spanID":  spanID,
		})
	}
	return r
}

func (r *LogrusLogger) log(level logrus.Level, message string) {
	entry := r.logrus.WithFields(r.getLogrusFields()) // Использование преобразованных fields
	entry.Log(level, message)
}

func (r *LogrusLogger) getLogrusFields() logrus.Fields {
	fields := logrus.Fields{}
	for key, value := range r.GetAllContexts() { // Преобразование params в logrus.Fields
		fields[key] = value
	}
	return fields
}

var _ logging.Logger = &LogrusLogger{}
