package logruslog

import (
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/vsysa/logging"
	"testing"
)

func TestBaseLogger_SetLevel(t *testing.T) {
	logger := NewLogrusLogger()
	logger.SetLevel(logging.DebugLevel)
	assert.Equal(t, logrus.DebugLevel, logger.logrus.GetLevel(), "Log level should be set to Debug")
}

//func TestBaseLogger_Output(t *testing.T) {
//	logger := NewLogrusLogger()
//	logger.SetLevel(logging.TraceLevel)
//	logger.Trace("This is trace log")
//	logger.Debug("This is debug log")
//	logger.Info("This is info log")
//	logger.Warn("This is warn log")
//	logger.Error("This is error log")
//}
