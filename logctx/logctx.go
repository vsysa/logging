package logctx

import (
	"context"
	"github.com/vsysa/logging"
	"github.com/vsysa/logging/factory"
)

type ctxLogger struct{}

var loggerFactory factory.LoggerFactory

func SetLoggerFactory(factory factory.LoggerFactory) {
	loggerFactory = factory
}

func CtxWithLogger(ctx context.Context, l logging.Logger) context.Context {
	return context.WithValue(ctx, ctxLogger{}, l)
}

func LoggerFromCtx(ctx context.Context) logging.Logger {
	if l, ok := ctx.Value(ctxLogger{}).(logging.Logger); ok {
		return l
	}

	if loggerFactory == nil {
		// Устанавливаем фабрику по умолчанию, если она не была установлена
		loggerFactory = factory.NewZapLoggerFactory(factory.NewZapLoggerDefault())
	}

	logger := loggerFactory.CreateLogger()
	logger.AddContext("type", "defaultLoggerFromCtx")
	return logger
}

func L(ctx context.Context) logging.Logger {
	return LoggerFromCtx(ctx)
}
