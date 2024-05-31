package logging

import (
	"context"
	"github.com/sirupsen/logrus"
	"time"
)

// context tools

type ctxLogger struct{}

func CtxWithLogger(ctx context.Context, l Logger) context.Context {
	return context.WithValue(ctx, ctxLogger{}, l)
}
func BackgroundCtxWithLogger(l Logger) context.Context {
	return context.WithValue(context.Background(), ctxLogger{}, l)
}

func LoggerFromCtx(ctx context.Context) Logger {
	if l, ok := ctx.Value(ctxLogger{}).(Logger); ok {
		return l
	}

	return NewBaseLoggerWithParams(ContextParams{
		"type": "defaultLoggerFromCtx",
	})
}
func LoggerCloneFromCtx(ctx context.Context) Logger {
	return LoggerFromCtx(ctx).Clone()
}
func CtxWithCloneLogger(ctx context.Context) context.Context {
	return CtxWithLogger(ctx, LoggerCloneFromCtx(ctx))
}

func L(ctx context.Context) Logger {
	return LoggerFromCtx(ctx)
}

func TestCtxWithLogger() context.Context {
	return BackgroundCtxWithLogger(NewTestLogger())
}

// converters

func convertContextToLogrusFields(strContext map[string]string) logrus.Fields {
	fields := logrus.Fields{}
	for key, value := range strContext {
		fields[key] = value
	}
	return fields
}

// cloning

func copyMapContext(original map[string]string) map[string]string {
	copy := make(map[string]string, len(original))
	for key, value := range original {
		copy[key] = value
	}
	return copy
}

func copyMapTime(original map[string]time.Time) map[string]time.Time {
	copy := make(map[string]time.Time, len(original))
	for key, value := range original {
		copy[key] = value
	}
	return copy
}
