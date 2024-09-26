package factory

import (
	"github.com/vsysa/logging"
	"github.com/vsysa/logging/logger/logruslog"
)

type LogrusLoggerFactory struct {
}

func (r *LogrusLoggerFactory) CreateLogger() logging.Logger {
	return logruslog.NewLogrusLogger()
}

var _ LoggerFactory = &LogrusLoggerFactory{}
