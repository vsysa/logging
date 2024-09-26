package logger

import (
	"github.com/stretchr/testify/assert"
	"github.com/vsysa/logging"
	"github.com/vsysa/logging/logger/logruslog"
	"github.com/vsysa/logging/logger/zaplog"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"testing"
)

func getLoggerForTest() []struct {
	name   string
	logger logging.Logger
} {
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder // CapitalColorLevelEncoder добавляет цвет в зависимости от уровня
	atomicLevel := zap.NewAtomicLevelAt(zap.DebugLevel)

	consoleCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.AddSync(zapcore.Lock(os.Stdout)),
		atomicLevel,
	)

	return []struct {
		name   string
		logger logging.Logger
	}{
		{"LogrusLogger", logruslog.NewLogrusLogger()},
		{"ZapLogger", zaplog.NewZapLogger(zap.New(consoleCore, zap.AddCaller(), zap.AddCallerSkip(2)), atomicLevel)},
	}
}

func TestBaseLogger_AddAndDeleteContext(t *testing.T) {
	tests := getLoggerForTest()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := tt.logger

			logger.AddContext("key1", "value1")
			logger.AddContext("key2", 2)

			// Test AddContexts
			logger.AddContexts(map[string]interface{}{"key3": 3.0, "key4": true})

			// Verify all contexts are added
			expectedContexts := map[string]interface{}{"key1": "value1", "key2": "2", "key3": "3", "key4": "true"}
			assert.Equal(t, expectedContexts, logger.GetAllContexts(), "Contexts should match expected")

			// Test DeleteContext
			logger.DeleteContext("key1")
			delete(expectedContexts, "key1")
			assert.Equal(t, expectedContexts, logger.GetAllContexts(), "Context after deletion should match expected")
		})
	}
}

func TestBaseLogger_Clone(t *testing.T) {
	tests := getLoggerForTest()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalLogger := tt.logger
			originalLogger.AddContext("key", "value")

			clonedLogger := originalLogger.Clone()

			// Test if contexts are cloned
			assert.Equal(t, originalLogger.GetAllContexts(), clonedLogger.GetAllContexts(), "Contexts should be equal")

			// Modify clone and ensure original is unaffected
			clonedLogger.AddContext("newKey", "newValue")
			assert.NotEqual(t, originalLogger.GetAllContexts(), clonedLogger.GetAllContexts(), "Original logger's context should remain unchanged")
		})
	}
}
