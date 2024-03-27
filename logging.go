package logging

import (
	"github.com/sirupsen/logrus"
	"time"
)

type Level = logrus.Level

const (
	DebugLevel Level = logrus.DebugLevel
	InfoLevel  Level = logrus.InfoLevel
	WarnLevel  Level = logrus.WarnLevel
	ErrorLevel Level = logrus.ErrorLevel
)

type ContextParams = map[string]interface{}

type Logger interface {
	AddContext(key string, value interface{}) Logger
	AddContexts(contexts ContextParams) Logger
	DeleteContext(key string) Logger
	GetAllContexts() map[string]string

	Debug(message string, a ...any)
	Info(message string, a ...any)
	Warn(message string, a ...any)
	Error(message string, a ...any)
	Fatal(message string, a ...any)
	ErrorCatch(message string, err error)
	FatalCatch(message string, err error)

	TimerStart(label string)
	TimerPrint(label string)
	TimerDuration(label string) (duration time.Duration, ok bool)

	SetLevel(level Level)
	Clone() Logger
}
