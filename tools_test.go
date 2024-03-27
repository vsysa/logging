package logging

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestTools_CtxWithAndFromLogger(t *testing.T) {
	logger := NewBaseLogger()
	ctx := context.Background()

	ctxWithLogger := CtxWithLogger(ctx, logger)
	retrievedLogger := LoggerFromCtx(ctxWithLogger)

	require.NotNil(t, retrievedLogger, "Logger should be retrieved successfully from context")

	// Ensure that retrieving a logger from a context without a logger returns a default logger
	ctxWithoutLogger := context.Background()
	defaultLogger := LoggerFromCtx(ctxWithoutLogger)
	assert.NotNil(t, defaultLogger, "Default logger should be returned")
	assert.NotEqual(t, defaultLogger, logger, "Default logger should not equal the original logger")
}

func TestTools_BackgroundCtxLogger(t *testing.T) {
	logger := NewBaseLogger()
	ctx := BackgroundCtxWithLogger(logger)

	retrievedLogger := LoggerFromCtx(ctx)
	require.NotNil(t, retrievedLogger, "Logger should be retrieved successfully from background context")

	retrievedLogger.AddContext("testKey", "testValue")
	originalLoggerContext := logger.GetAllContexts()
	_, exists := originalLoggerContext["testKey"]
	assert.True(t, exists, "Modification in retrievedLogger should reflect in original logger, indicating same instance")
}

func TestTools_LoggerCloneFromCtx(t *testing.T) {
	originalLogger := NewBaseLogger()
	originalLogger.AddContext("key", "originalValue")
	ctx := BackgroundCtxWithLogger(originalLogger)

	clonedLogger := LoggerCloneFromCtx(ctx)
	clonedLogger.AddContext("key", "clonedValue")

	originalContext := originalLogger.GetAllContexts()
	clonedContext := clonedLogger.GetAllContexts()

	assert.NotEqual(t, originalContext["key"], clonedContext["key"], "Cloned logger should have independent state")
}

func TestTools_CtxWithCloneLogger(t *testing.T) {
	originalLogger := NewBaseLogger()
	originalLogger.AddContext("key", "originalValue")
	ctx := CtxWithLogger(context.Background(), originalLogger)

	clonedCtx := CtxWithCloneLogger(ctx)
	clonedLogger := LoggerFromCtx(clonedCtx)

	clonedLogger.AddContext("key", "clonedValue")

	originalContext := originalLogger.GetAllContexts()
	clonedContext := clonedLogger.GetAllContexts()

	assert.Equal(t, "originalValue", originalContext["key"], "Original logger's context should remain unchanged")
	assert.Equal(t, "clonedValue", clonedContext["key"], "Cloned logger's context should reflect the new value")
}

func TestTools_ConvertContextToLogrusFields(t *testing.T) {
	contextMap := map[string]string{"key1": "value1", "key2": "value2"}
	logrusFields := convertContextToLogrusFields(contextMap)

	expectedFields := logrus.Fields{"key1": "value1", "key2": "value2"}
	assert.Equal(t, expectedFields, logrusFields, "Logrus fields should match context map")
}

func TestTools_CopyMapContext(t *testing.T) {
	original := map[string]string{"key1": "value1", "key2": "value2"}
	copied := copyMapContext(original)

	assert.Equal(t, original, copied, "Copied map should equal the original")

	// Ensure modifications to the copy do not affect the original
	copied["key3"] = "value3"
	assert.NotEqual(t, original, copied, "Original should not be modified when the copy is changed")
}

func TestTools_CopyMapTime(t *testing.T) {
	now := time.Now()
	original := map[string]time.Time{"timer1": now}
	copied := copyMapTime(original)

	assert.Equal(t, original, copied, "Copied map should equal the original")

	// Ensure modifications to the copy do not affect the original
	copied["timer2"] = now.Add(1 * time.Hour)
	assert.NotEqual(t, original, copied, "Original should not be modified when the copy is changed")
}
