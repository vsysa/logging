package logging

import (
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTestLogger_AddAndDeleteContext(t *testing.T) {
	logger := NewTestLogger()
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

func TestTestLogger_SetLevel(t *testing.T) {
	logger := NewTestLogger()
	logger.SetLevel(DebugLevel)
	assert.Equal(t, logrus.DebugLevel, logger.logrus.GetLevel(), "Log level should be set to Debug")
}

func TestTestLogger_TimerFunctionality(t *testing.T) {
	logger := NewTestLogger()
	label := "testTimer"

	logger.TimerStart(label)
	time.Sleep(10 * time.Millisecond) // Sleep to simulate elapsed time

	duration, ok := logger.TimerDuration(label)
	require.True(t, ok, "Timer should exist")
	assert.GreaterOrEqual(t, duration, 10*time.Millisecond, "Duration should be greater or equal than 10ms")
	assert.LessOrEqual(t, duration, 12*time.Millisecond, "Duration should be less or equal than 11ms")
}

func TestTestLogger_Clone(t *testing.T) {
	originalLogger := NewTestLogger()
	originalLogger.AddContext("key", "value")

	clonedLogger := originalLogger.Clone().(*TestLogger)

	// Test if contexts are cloned
	require.Equal(t, originalLogger.context, clonedLogger.context, "Contexts should be equal")

	// Modify clone and ensure original is unaffected
	clonedLogger.AddContext("newKey", "newValue")
	assert.NotEqual(t, originalLogger.context, clonedLogger.context, "Original logger's context should remain unchanged")
}
