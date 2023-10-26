package logger

import (
	"go.uber.org/zap"
)

var globalLogger *zap.SugaredLogger

func GetLogger() *zap.SugaredLogger {
	if globalLogger == nil {
		logger, _ := zap.NewDevelopment()

		defer logger.Sync() // flushes buffer, if any
		globalLogger = logger.Sugar()

	}

	return globalLogger
}
