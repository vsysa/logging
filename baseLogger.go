package logging

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"sync"
	"time"
)

func NewBaseLoggerWithParams(params ContextParams) *BaseLogger {
	newLogger := NewBaseLogger()
	newLogger.AddContexts(params)
	return newLogger
}
func NewBaseLogger() *BaseLogger {
	l := logrus.New()
	l.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02T15:04:05.999",
		ForceColors:     true,
		FullTimestamp:   true,
	})

	newLogger := &BaseLogger{
		logrus:  l,
		context: make(map[string]string),
		timers:  make(map[string]time.Time),
	}

	return newLogger
}

type BaseLogger struct {
	logrus    *logrus.Logger
	contextMu sync.RWMutex
	context   map[string]string
	timerMu   sync.RWMutex
	timers    map[string]time.Time
}

//	CONTEXT

// AddContext добавляет ключ и значение в LoggerContext.
func (r *BaseLogger) AddContext(key string, value interface{}) Logger {
	r.contextMu.Lock()
	defer r.contextMu.Unlock()
	r.context[key] = fmt.Sprintf("%v", value)
	return r
}

// AddContexts добавляет несколько ключей и значений в LoggerContext.
func (r *BaseLogger) AddContexts(contexts map[string]interface{}) Logger {
	r.contextMu.Lock()
	defer r.contextMu.Unlock()
	for key, value := range contexts {
		r.context[key] = fmt.Sprintf("%v", value)
	}
	return r
}

// DeleteContext удаляет значение по ключу из LoggerContext.
func (r *BaseLogger) DeleteContext(key string) Logger {
	r.contextMu.Lock()
	defer r.contextMu.Unlock()
	delete(r.context, key)
	return r
}

// GetAllContext возвращает все ключи и значения из LoggerContext.
func (r *BaseLogger) GetAllContexts() map[string]string {
	r.contextMu.RLock()
	defer r.contextMu.RUnlock()
	copy := make(map[string]string)
	for key, value := range r.context {
		copy[key] = value
	}
	return copy
}

//	LOGGING

func (r *BaseLogger) SetLevel(level Level) {
	r.logrus.SetLevel(level)
}

func (r *BaseLogger) getLogrusFields(params map[string]interface{}) logrus.Fields {
	fields := logrus.Fields{}
	for key, value := range r.GetAllContexts() { // Преобразование params в logrus.Fields
		fields[key] = value
	}
	if params != nil {
		for key, value := range params {
			fields[key] = value
		}
	}
	return fields
}

func (r *BaseLogger) log(level logrus.Level, message string, params map[string]interface{}) {
	entry := r.logrus.WithFields(r.getLogrusFields(params)) // Использование преобразованных fields
	entry.Log(level, message)
}

func (r *BaseLogger) Debug(message string, a ...any) {
	r.log(logrus.DebugLevel, fmt.Sprintf(message, a...), nil)
}

func (r *BaseLogger) Info(message string, a ...any) {
	r.log(logrus.InfoLevel, fmt.Sprintf(message, a...), nil)
}

func (r *BaseLogger) Warn(message string, a ...any) {
	r.log(logrus.WarnLevel, fmt.Sprintf(message, a...), nil)
}

func (r *BaseLogger) Error(message string, a ...any) {
	r.log(logrus.ErrorLevel, fmt.Sprintf(message, a...), nil)
}

func (r *BaseLogger) Fatal(message string, a ...any) {
	r.log(logrus.FatalLevel, fmt.Sprintf(message, a...), nil)
}

func (r *BaseLogger) ErrorCatch(message string, err error) {
	if err != nil {
		message += ": " + err.Error()
	}
	if err != nil {
		r.log(logrus.ErrorLevel, message, nil)
	}
}

func (r *BaseLogger) FatalCatch(message string, err error) {
	if err != nil {
		message += ": " + err.Error()
	}
	if err != nil {
		r.log(logrus.FatalLevel, message, nil)
	}
	os.Exit(1)
}

//	TIMER

func (r *BaseLogger) TimerStart(label string) {
	r.timerMu.Lock()
	defer r.timerMu.Unlock()
	r.timers[label] = time.Now()
}

func (r *BaseLogger) TimerPrint(label string) {
	if duration, ok := r.TimerDuration(label); ok {
		r.Warn(fmt.Sprintf("Timer %s: %s", label, duration))
	} else {
		r.Warn(fmt.Sprintf("Unknown timer with lable \"%s\"", label))
	}
}

func (r *BaseLogger) TimerDuration(label string) (duration time.Duration, ok bool) {
	r.timerMu.RLock()
	defer r.timerMu.RUnlock()
	if startTime, ok := r.timers[label]; ok {
		return time.Since(startTime), true // Возвращаем время в миллисекундах
	}
	return 0, false
}

// BASE

func (r *BaseLogger) Clone() Logger {
	newContext := make(map[string]string)
	for key, value := range r.context {
		newContext[key] = value
	}
	newTimers := make(map[string]time.Time)
	for key, value := range r.timers {
		newTimers[key] = value
	}
	newLogger := &BaseLogger{
		logrus:  r.logrus,
		context: newContext,
		timers:  newTimers,
	}

	return newLogger
}

var _ Logger = &BaseLogger{}
