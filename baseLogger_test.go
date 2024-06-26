package logging

import (
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestBaseLogger_AddAndDeleteContext(t *testing.T) {
	logger := NewBaseLogger()
	logger.AddContext("key1", "value1")
	logger.AddContext("key2", 2)

	// Test AddContexts
	logger.AddContexts(map[string]interface{}{"key3": 3.0, "key4": true})

	// Verify all contexts are added
	expectedContexts := map[string]string{"key1": "value1", "key2": "2", "key3": "3", "key4": "true"}
	assert.Equal(t, expectedContexts, logger.GetAllContexts(), "Contexts should match expected")

	// Test DeleteContext
	logger.DeleteContext("key1")
	delete(expectedContexts, "key1")
	assert.Equal(t, expectedContexts, logger.GetAllContexts(), "Context after deletion should match expected")
}

func TestBaseLogger_SetLevel(t *testing.T) {
	logger := NewBaseLogger()
	logger.SetLevel(DebugLevel)
	assert.Equal(t, logrus.DebugLevel, logger.logrus.GetLevel(), "Log level should be set to Debug")
}

func TestBaseLogger_TimerFunctionality(t *testing.T) {
	logger := NewBaseLogger()
	label := "testTimer"

	logger.TimerStart(label)
	time.Sleep(10 * time.Millisecond) // Sleep to simulate elapsed time

	duration, ok := logger.TimerDuration(label)
	require.True(t, ok, "Timer should exist")
	assert.GreaterOrEqual(t, duration, 10*time.Millisecond, "Duration should be greater or equal than 10ms")
	assert.LessOrEqual(t, duration, 12*time.Millisecond, "Duration should be less or equal than 11ms")
}

func TestBaseLogger_Clone(t *testing.T) {
	originalLogger := NewBaseLogger()
	originalLogger.AddContext("key", "value")

	clonedLogger := originalLogger.Clone().(*BaseLogger)

	// Test if contexts are cloned
	require.Equal(t, originalLogger.context, clonedLogger.context, "Contexts should be equal")

	// Modify clone and ensure original is unaffected
	clonedLogger.AddContext("newKey", "newValue")
	assert.NotEqual(t, originalLogger.context, clonedLogger.context, "Original logger's context should remain unchanged")
}

func TestBaseLogger_Output(t *testing.T) {
	logger := NewBaseLogger()
	logger.SetLevel(TraceLevel)
	logger.Trace("This is trace log")
	logger.Debug("This is debug log")
	logger.Info("This is info log")
	logger.Warn("This is warn log")
	logger.Error("This is error log")
}
