package test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vsysa/logging/logger/logruslog"
	"testing"
)

func TestTestLogger_AddAndDeleteContext(t *testing.T) {
	logger := NewTestLogger(logruslog.NewLogrusLogger())
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
}

func TestTestLogger_SetLevel(t *testing.T) {
	//logger := NewTestLogger()
	//logger.SetLevel(DebugLevel)
	//assert.Equal(t, DebugLevel, logger.logrus.GetLevel(), "Log level should be set to Debug")
}

func TestTestLogger_Clone(t *testing.T) {
	originalLogger := NewTestLogger(logruslog.NewLogrusLogger())
	originalLogger.AddContext("key", "value")

	clonedLogger := originalLogger.Clone().(*TestLogger)

	// Test if contexts are cloned
	require.Equal(t, originalLogger.context, clonedLogger.context, "Contexts should be equal")

	// Modify clone and ensure original is unaffected
	clonedLogger.AddContext("newKey", "newValue")
	assert.NotEqual(t, originalLogger.context, clonedLogger.context, "Original logger's context should remain unchanged")
}
