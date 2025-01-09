package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

const (
	currentLogsPath = "_logs"
	serviceName     = "notification-service"
)

func New(logLevel string) *zap.Logger {
	config := zap.NewProductionEncoderConfig()

	config.MessageKey = "message"
	config.TimeKey = "timestamp"
	config.EncodeTime = zapcore.ISO8601TimeEncoder

	fileEncoder := zapcore.NewJSONEncoder(config)
	consoleEncoder := zapcore.NewConsoleEncoder(config)

	level, err := zapcore.ParseLevel(logLevel)
	if err != nil {
		return zap.New(zapcore.NewCore(consoleEncoder, os.Stdout, zap.DebugLevel))
	}

	if _, err = os.Stat(currentLogsPath); os.IsNotExist(err) {
		if err = os.Mkdir(currentLogsPath, 0777); err != nil {
			return zap.New(zapcore.NewCore(consoleEncoder, os.Stdout, zap.DebugLevel))
		}
	}

	filename := currentLogsPath + "/" + serviceName + ".log"
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return zap.New(zapcore.NewCore(consoleEncoder, os.Stdout, zap.DebugLevel))
	}

	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, zapcore.AddSync(file), level),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), level),
	)

	return zap.New(core, zap.AddCaller())
}
