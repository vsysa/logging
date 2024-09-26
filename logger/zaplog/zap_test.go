package zaplog

import (
	"github.com/stretchr/testify/assert"
	"github.com/vsysa/logging"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"testing"
)

func getZapLogger() (*ZapLogger, zap.AtomicLevel) {
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder // CapitalColorLevelEncoder добавляет цвет в зависимости от уровня
	atomicLevel := zap.NewAtomicLevelAt(zap.DebugLevel)

	consoleCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.AddSync(zapcore.Lock(os.Stdout)),
		atomicLevel,
	)

	return NewZapLogger(zap.New(consoleCore, zap.AddCaller(), zap.AddCallerSkip(2)), atomicLevel), atomicLevel
}

func TestBaseLogger_SetLevel(t *testing.T) {
	logger, al := getZapLogger()
	logger.SetLevel(logging.InfoLevel)
	assert.Equal(t, levelMap[logging.InfoLevel], al.Level(), "Log level should be set to Debug")
}

//func TestBaseLogger_Output(t *testing.T) {
//	logger, _ := getZapLogger()
//	logger.SetLevel(logging.TraceLevel)
//	logger.Trace("This is trace log")
//	logger.Debug("This is debug log")
//	logger.Info("This is info log")
//	logger.Warn("This is warn log")
//	logger.Error("This is error log")
//}
