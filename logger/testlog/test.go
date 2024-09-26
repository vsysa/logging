package test

import (
	"context"
	"fmt"
	"github.com/vsysa/logging"
	"github.com/vsysa/logging/internal/helper"
	"os"
	"sync"
)

type logStoreStruct struct {
	level   logging.Level
	message string
	context map[string]interface{}
}

type TestLogger struct {
	outLogger logging.Logger

	contextMu sync.RWMutex
	context   map[string]interface{}

	logStoreMu   sync.RWMutex
	logStore     []logStoreStruct
	parentLogger *TestLogger
	rootLogger   *TestLogger
}

func NewTestLogger(outLogger logging.Logger) *TestLogger {
	return &TestLogger{
		outLogger: outLogger,
		context:   make(map[string]interface{}),
	}
}

func NewTestLoggerWithParams(outLogger logging.Logger, params map[string]interface{}) *TestLogger {
	newLogger := NewTestLogger(outLogger)
	newLogger.AddContexts(params)
	return newLogger
}

//	CONTEXT

// AddContext добавляет ключ и значение в LoggerContext.
func (r *TestLogger) AddContext(key string, value interface{}) logging.Logger {
	r.contextMu.Lock()
	defer r.contextMu.Unlock()
	r.context[key] = fmt.Sprintf("%v", value)
	return r
}

// AddContexts добавляет несколько ключей и значений в LoggerContext.
func (r *TestLogger) AddContexts(contexts map[string]interface{}) logging.Logger {
	r.contextMu.Lock()
	defer r.contextMu.Unlock()
	for key, value := range contexts {
		r.context[key] = fmt.Sprintf("%v", value)
	}
	return r
}

// DeleteContext удаляет значение по ключу из LoggerContext.
func (r *TestLogger) DeleteContext(key string) logging.Logger {
	r.contextMu.Lock()
	defer r.contextMu.Unlock()
	delete(r.context, key)
	return r
}

func (r *TestLogger) GetAllContexts() map[string]interface{} {
	r.contextMu.RLock()
	defer r.contextMu.RUnlock()
	clone := make(map[string]interface{})
	for key, value := range r.context {
		clone[key] = value
	}
	return clone
}

//	LOGGING

func (r *TestLogger) SetLevel(level logging.Level) {
	//r.outLogger.SetLevel(level)
}

func (r *TestLogger) Trace(message string, a ...any) {
	r.log(logging.TraceLevel, fmt.Sprintf(message, a...))
}

func (r *TestLogger) Debug(message string, a ...any) {
	r.log(logging.DebugLevel, fmt.Sprintf(message, a...))
}

func (r *TestLogger) Info(message string, a ...any) {
	r.log(logging.InfoLevel, fmt.Sprintf(message, a...))
}

func (r *TestLogger) Warn(message string, a ...any) {
	r.log(logging.WarnLevel, fmt.Sprintf(message, a...))
}

func (r *TestLogger) Error(message string, a ...any) {
	r.log(logging.ErrorLevel, fmt.Sprintf(message, a...))
}

func (r *TestLogger) Fatal(message string, a ...any) {
	r.log(logging.FatalLevel, fmt.Sprintf(message, a...))
}

func (r *TestLogger) ErrorCatch(err error, message string, a ...any) {
	msg := fmt.Sprintf(message, a...)
	if err != nil {
		msg += ": " + err.Error()
	}
	r.log(logging.ErrorLevel, msg)
}

func (r *TestLogger) FatalCatch(err error, message string, a ...any) {
	msg := fmt.Sprintf(message, a...)
	if err != nil {
		msg += ": " + err.Error()
	}
	r.log(logging.FatalLevel, msg)
	os.Exit(1)
}

// BASE

func (r *TestLogger) Clone() logging.Logger {
	r.contextMu.RLock()
	defer r.contextMu.RUnlock()

	rootLogger := r
	if r.rootLogger != nil {
		rootLogger = r.rootLogger
	}
	return &TestLogger{
		outLogger:    r.outLogger,
		context:      helper.CopyMapContext(r.context),
		parentLogger: r,
		rootLogger:   rootLogger,
	}
}

func (r *TestLogger) ShowStoredLogs() {
	r.logStoreMu.RLock()
	defer r.logStoreMu.RUnlock()
	fmt.Printf("\n**** LOGS BEFORE ERROR\n")
	for _, storedLog := range r.logStore {
		// TODO что-то придумать бы поэлегантней
		localLogger := r.outLogger.Clone()
		localLogger.AddContexts(storedLog.context)
		switch storedLog.level {
		case logging.DebugLevel:
			localLogger.Debug(storedLog.message)
		case logging.InfoLevel:
			localLogger.Info(storedLog.message)
		case logging.WarnLevel:
			localLogger.Warn(storedLog.message)
		case logging.ErrorLevel:
			localLogger.Error(storedLog.message)
		case logging.FatalLevel:
			localLogger.Fatal(storedLog.message)
		}
	}
}

func (r *TestLogger) storeLog(level logging.Level, message string, context map[string]interface{}) {
	if r.parentLogger != nil {
		r.parentLogger.storeLog(level, message, context)
		return
	}
	interfaceContext := make(map[string]interface{})
	for k, v := range context {
		interfaceContext[k] = v
	}
	// можно поменять местами, если нужно чтоб каждый логер хранил в себе информацию о его логах и логах его дочерних логеров
	r.logStoreMu.Lock()
	r.logStore = append(r.logStore, logStoreStruct{
		level:   level,
		message: message,
		context: interfaceContext,
	})
	r.logStoreMu.Unlock()
}

func (r *TestLogger) log(level logging.Level, message string) {
	r.storeLog(level, message, r.GetAllContexts())
}

func (r *TestLogger) SetCtx(ctx context.Context) logging.Logger {
	return r
}

var _ logging.Logger = &TestLogger{}
