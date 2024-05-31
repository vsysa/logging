package logging

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"sync"
	"time"
)

type logStoreStruct struct {
	level   logrus.Level
	message string
	context map[string]string
}

func NewTestLoggerWithParams(params ContextParams) *TestLogger {
	newLogger := NewTestLogger()
	newLogger.AddContexts(params)
	return newLogger
}
func NewTestLogger() *TestLogger {
	l := logrus.New()
	l.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02T15:04:05.999",
		ForceColors:     true,
		FullTimestamp:   true,
	})

	newLogger := &TestLogger{
		logrus:  l,
		context: make(map[string]string),
		timers:  make(map[string]time.Time),
	}

	return newLogger
}

type TestLogger struct {
	logrus *logrus.Logger

	contextMu sync.RWMutex
	context   map[string]string

	timerMu sync.RWMutex
	timers  map[string]time.Time

	logStoreMu   sync.RWMutex
	logStore     []logStoreStruct
	parentLogger *TestLogger
	rootLogger   *TestLogger
}

//	CONTEXT

// AddContext добавляет ключ и значение в LoggerContext.
func (r *TestLogger) AddContext(key string, value interface{}) Logger {
	r.contextMu.Lock()
	defer r.contextMu.Unlock()
	r.context[key] = fmt.Sprintf("%v", value)
	return r
}

// AddContexts добавляет несколько ключей и значений в LoggerContext.
func (r *TestLogger) AddContexts(contexts map[string]interface{}) Logger {
	r.contextMu.Lock()
	defer r.contextMu.Unlock()
	for key, value := range contexts {
		r.context[key] = fmt.Sprintf("%v", value)
	}
	return r
}

// DeleteContext удаляет значение по ключу из LoggerContext.
func (r *TestLogger) DeleteContext(key string) Logger {
	r.contextMu.Lock()
	defer r.contextMu.Unlock()
	delete(r.context, key)
	return r
}

func (r *TestLogger) GetAllContexts() map[string]string {
	r.contextMu.RLock()
	defer r.contextMu.RUnlock()
	copy := make(map[string]string)
	for key, value := range r.context {
		copy[key] = value
	}
	return copy
}

//	LOGGING

func (r *TestLogger) SetLevel(level Level) {
	r.logrus.SetLevel(level)
}

func (r *TestLogger) Trace(message string, a ...any) {
	r.log(logrus.TraceLevel, fmt.Sprintf(message, a...))
}

func (r *TestLogger) Debug(message string, a ...any) {
	r.log(logrus.DebugLevel, fmt.Sprintf(message, a...))
}

func (r *TestLogger) Info(message string, a ...any) {
	r.log(logrus.InfoLevel, fmt.Sprintf(message, a...))
}

func (r *TestLogger) Warn(message string, a ...any) {
	r.log(logrus.WarnLevel, fmt.Sprintf(message, a...))
}

func (r *TestLogger) Error(message string, a ...any) {
	r.log(logrus.ErrorLevel, fmt.Sprintf(message, a...))
}

func (r *TestLogger) Fatal(message string, a ...any) {
	r.log(logrus.FatalLevel, fmt.Sprintf(message, a...))
}

func (r *TestLogger) ErrorCatch(message string, err error) {
	if err != nil {
		message += ": " + err.Error()
	}
	if err != nil {
		r.log(logrus.ErrorLevel, message)
	}
}

func (r *TestLogger) FatalCatch(message string, err error) {
	if err != nil {
		message += ": " + err.Error()
	}
	if err != nil {
		r.log(logrus.FatalLevel, message)
	}
	os.Exit(1)
}

//	TIMER

func (r *TestLogger) TimerStart(label string) {
	r.timerMu.Lock()
	defer r.timerMu.Unlock()
	r.timers[label] = time.Now()
}

func (r *TestLogger) TimerPrint(label string) {
	if duration, ok := r.TimerDuration(label); ok {
		r.Warn(fmt.Sprintf("Timer %s: %s", label, duration))
	} else {
		r.Warn(fmt.Sprintf("Unknown timer with lable \"%s\"", label))
	}
}

func (r *TestLogger) TimerDuration(label string) (duration time.Duration, ok bool) {
	r.timerMu.RLock()
	defer r.timerMu.RUnlock()
	if startTime, ok := r.timers[label]; ok {
		return time.Since(startTime), true
	}
	return 0, false
}

// BASE

func (r *TestLogger) Clone() Logger {
	r.timerMu.RLock()
	defer r.timerMu.RUnlock()
	r.contextMu.RLock()
	defer r.contextMu.RUnlock()

	rootLogger := r
	if r.rootLogger != nil {
		rootLogger = r.rootLogger
	}
	return &TestLogger{
		logrus:       r.logrus,
		context:      copyMapContext(r.context),
		timers:       copyMapTime(r.timers),
		parentLogger: r,
		rootLogger:   rootLogger,
	}
}

func (r *TestLogger) ShowStoredLogs() {
	r.logStoreMu.RLock()
	defer r.logStoreMu.RUnlock()
	fmt.Printf("\n**** LOGS BEFORE ERROR\n")
	for _, storedLog := range r.logStore {
		entry := r.logrus.WithFields(convertContextToLogrusFields(storedLog.context)) // Использование преобразованных fields
		entry.Log(storedLog.level, fmt.Sprintf("\t\t%s", storedLog.message))
	}
}

func (r *TestLogger) storeLog(level logrus.Level, message string, context map[string]string) {
	if r.parentLogger != nil {
		r.parentLogger.storeLog(level, message, context)
		return
	}
	// можно поменять местами, если нужно чтоб каждый логер хранил в себе информацию о его логах и логах его дочерних логеров
	r.logStoreMu.Lock()
	r.logStore = append(r.logStore, logStoreStruct{
		level:   level,
		message: message,
		context: context,
	})
	r.logStoreMu.Unlock()
}

func (r *TestLogger) log(level logrus.Level, message string) {
	r.storeLog(level, message, r.GetAllContexts())
}

var _ Logger = &TestLogger{}
