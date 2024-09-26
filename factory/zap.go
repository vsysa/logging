package factory

import (
	"github.com/vsysa/logging"
	"github.com/vsysa/logging/logger/zaplog"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

type ZapLoggerFactory struct {
	zapLogger   *zap.Logger
	atomicLevel zap.AtomicLevel
}

func NewZapLoggerFactory(zapLogger *zap.Logger, atomicLevel zap.AtomicLevel) *ZapLoggerFactory {
	return &ZapLoggerFactory{zapLogger: zapLogger, atomicLevel: atomicLevel}
}

func NewZapLoggerDefault() (*zap.Logger, zap.AtomicLevel) {
	// Настройка цветного вывода для консоли
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder // CapitalColorLevelEncoder добавляет цвет в зависимости от уровня
	atomicLevel := zap.NewAtomicLevelAt(zap.DebugLevel)

	consoleCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.AddSync(zapcore.Lock(os.Stdout)),
		atomicLevel,
	)

	return zap.New(consoleCore, zap.AddCaller(), zap.AddCallerSkip(2)), atomicLevel
}

func (r *ZapLoggerFactory) CreateLogger() logging.Logger {
	zl := r.zapLogger
	zal := r.atomicLevel
	if zl == nil {
		zl, zal = NewZapLoggerDefault()
	}
	return zaplog.NewZapLogger(zl, zal)
}

var _ LoggerFactory = &ZapLoggerFactory{}
