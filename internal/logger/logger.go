// Package logger пакет для работы с логами
package logger

import "go.uber.org/zap"

// Log инициализация объекта для возможности логгировать действия
func Logger() *zap.SugaredLogger {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	return logger.Sugar()
}
