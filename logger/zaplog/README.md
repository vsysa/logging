# ZapLogger

## Logger Initialization

### Example: Basic Zap Logger Initialization


```go
package main

import (
    "os"
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
    "github.com/yourusername/zaplog"
)

func main() {
    // Configuring colored output for the console
    encoderConfig := zap.NewDevelopmentEncoderConfig()
    encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder // Adds color based on log level

    atomicLevel := zap.NewAtomicLevelAt(zap.DebugLevel)

    consoleCore := zapcore.NewCore(
        zapcore.NewConsoleEncoder(encoderConfig),
        zapcore.AddSync(zapcore.Lock(os.Stdout)),
        atomicLevel,
    )

    logger := zaplog.NewZapLogger(consoleCore, atomicLevel)

    // Using the logger
    logger.Info("Application started")
}
```

### Example: Zap Logger with File Writing and Log Rotation


```go
package main

import (
	"os"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"github.com/yourusername/zaplog"
)

func main() {
	// Setting up log rotation with lumberjack
	rotator := &lumberjack.Logger{
		Filename:   "logs/myapp.log",
		MaxSize:    10,   // Max size in MB before the log file is rotated
		MaxBackups: 5,    // Max number of old log files to keep
		MaxAge:     28,   // Max number of days to retain old log files
		Compress:   true, // Compress old log files
	}

	atomicLevel := zap.NewAtomicLevelAt(zap.InfoLevel)

	// Core for file logging
	fileCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.AddSync(rotator),
		atomicLevel,
	)

	// Configuring colored output for the console
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder // Adds color based on log level

	// Core for console logging
	consoleCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.AddSync(zapcore.Lock(os.Stdout)),
		atomicLevel,
	)

	// Combining console and file logging
	combinedCore := zapcore.NewTee(fileCore, consoleCore)

	zapLogger := zap.New(combinedCore, zap.AddCaller(), zap.AddCallerSkip(2))

	// Creating an instance of ZapLogger with the configured zap.Logger
	logger := zaplog.NewZapLogger(zapLogger, atomicLevel)

	// Using the logger
	logger.Info("Application started")
}
```