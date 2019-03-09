package components

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

func SetupLogger(isDebug bool) error {
	highPriority := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level >= zapcore.ErrorLevel
	})

	lowPriority := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level < zapcore.ErrorLevel
	})

	consoleDebugging := zapcore.Lock(os.Stdout)
	consoleErrors := zapcore.Lock(os.Stderr)

	var consoleEncoder zapcore.Encoder

	if isDebug {
		consoleEncoder = zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	} else {
		consoleEncoder = zapcore.NewConsoleEncoder(zap.NewProductionEncoderConfig())
	}

	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, consoleErrors, highPriority),
		zapcore.NewCore(consoleEncoder, consoleDebugging, lowPriority),
	)

	logger := zap.New(core)

	zap.ReplaceGlobals(logger)

	return nil
}
