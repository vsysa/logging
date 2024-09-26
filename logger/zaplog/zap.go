package zaplog

import (
	"context"
	"fmt"
	"github.com/vsysa/logging"
	"github.com/vsysa/logging/internal/helper"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"sync"
)

var levelMap = map[logging.Level]zapcore.Level{
	logging.TraceLevel: zapcore.DebugLevel, // Zap не имеет Trace уровня, используем Debug
	logging.DebugLevel: zapcore.DebugLevel,
	logging.InfoLevel:  zapcore.InfoLevel,
	logging.WarnLevel:  zapcore.WarnLevel,
	logging.ErrorLevel: zapcore.ErrorLevel,
	logging.FatalLevel: zapcore.FatalLevel,
}

type ZapLogger struct {
	zapLogger   *zap.Logger
	contextMu   sync.RWMutex
	context     map[string]interface{}
	atomicLevel zap.AtomicLevel
}

func NewZapLogger(zapLogger *zap.Logger, atomicLevel zap.AtomicLevel) *ZapLogger {
	if zapLogger == nil {
		var err error
		zapLogger, err = zap.NewProduction() // Можно использовать zap.NewProduction() или любую другую конфигурацию
		if err != nil {
			panic("Failed to create zap logger: " + err.Error())
		}
	}

	newLogger := &ZapLogger{
		zapLogger:   zapLogger,
		context:     make(map[string]interface{}),
		atomicLevel: atomicLevel,
	}

	return newLogger
}

// CONTEXT

func (r *ZapLogger) AddContext(key string, value interface{}) logging.Logger {
	r.contextMu.Lock()
	defer r.contextMu.Unlock()
	r.context[key] = fmt.Sprintf("%v", value)
	return r
}

func (r *ZapLogger) AddContexts(contexts map[string]interface{}) logging.Logger {
	r.contextMu.Lock()
	defer r.contextMu.Unlock()
	for key, value := range contexts {
		r.context[key] = fmt.Sprintf("%v", value)
	}
	return r
}

func (r *ZapLogger) DeleteContext(key string) logging.Logger {
	r.contextMu.Lock()
	defer r.contextMu.Unlock()
	delete(r.context, key)
	return r
}

func (r *ZapLogger) GetAllContexts() map[string]interface{} {
	r.contextMu.RLock()
	defer r.contextMu.RUnlock()
	newContext := make(map[string]interface{}, len(r.context))
	for key, value := range r.context {
		newContext[key] = value
	}
	return newContext
}

// LOGGING

func (r *ZapLogger) SetLevel(level logging.Level) {
	if zapLevel, ok := levelMap[level]; ok {
		r.atomicLevel.SetLevel(zapLevel)
	} else {
		r.Warn("Invalid logging level: %v", level)
		return
	}
}

func (r *ZapLogger) Trace(message string, a ...any) {
	r.log(zap.DebugLevel, fmt.Sprintf(message, a...))
}

func (r *ZapLogger) Debug(message string, a ...any) {
	r.log(zap.DebugLevel, fmt.Sprintf(message, a...))
}

func (r *ZapLogger) Info(message string, a ...any) {
	r.log(zap.InfoLevel, fmt.Sprintf(message, a...))
}

func (r *ZapLogger) Warn(message string, a ...any) {
	r.log(zap.WarnLevel, fmt.Sprintf(message, a...))
}

func (r *ZapLogger) Error(message string, a ...any) {
	r.log(zap.ErrorLevel, fmt.Sprintf(message, a...))
}

func (r *ZapLogger) Fatal(message string, a ...any) {
	r.log(zap.FatalLevel, fmt.Sprintf(message, a...))
}

func (r *ZapLogger) ErrorCatch(err error, message string, a ...any) {
	msg := fmt.Sprintf(message, a...)
	if err != nil {
		msg += ": " + err.Error()
	}
	r.log(zap.ErrorLevel, msg)
}

func (r *ZapLogger) FatalCatch(err error, message string, a ...any) {
	msg := fmt.Sprintf(message, a...)
	if err != nil {
		msg += ": " + err.Error()
	}
	r.log(zap.FatalLevel, msg)
}

// BASE

func (r *ZapLogger) Clone() logging.Logger {
	return &ZapLogger{
		zapLogger: r.zapLogger,
		context:   helper.CopyMapContext(r.context),
	}
}

func (r *ZapLogger) SetCtx(ctx context.Context) logging.Logger {
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

func (r *ZapLogger) log(level zapcore.Level, message string) {
	fields := r.getZapFields()
	r.zapLogger.With(fields...).Check(level, message).Write()
}

func (r *ZapLogger) getZapFields() []zap.Field {
	fields := make([]zap.Field, 0, len(r.context))
	for key, value := range r.GetAllContexts() {
		fields = append(fields, zap.Any(key, value))
	}
	return fields
}

var _ logging.Logger = &ZapLogger{}
