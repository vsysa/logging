package logging

import "context"

type Logger interface {
	AddContext(key string, value interface{}) Logger
	AddContexts(contexts map[string]interface{}) Logger
	DeleteContext(key string) Logger
	GetAllContexts() map[string]interface{}

	Trace(message string, a ...any)
	Debug(message string, a ...any)
	Info(message string, a ...any)
	Warn(message string, a ...any)
	Error(message string, a ...any)
	Fatal(message string, a ...any)
	ErrorCatch(err error, message string, a ...any)
	FatalCatch(err error, message string, a ...any)

	SetLevel(level Level)
	Clone() Logger
	SetCtx(ctx context.Context) Logger
}
