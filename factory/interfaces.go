package factory

import (
	"github.com/vsysa/logging"
)

type LoggerFactory interface {
	CreateLogger() logging.Logger
}
