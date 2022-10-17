package logger

import "go.uber.org/zap"

type Engine interface {
	LogError(message string, err error)
	LogInfo(message string, value ...interface{})
}

type Logger struct {
	zap *zap.Logger
}

func NewLogger(logger *zap.Logger) *Logger {
	return &Logger{
		zap: logger,
	}
}

func (l *Logger) LogError(message string, err error) {
	l.zap.Error(
		message,
		zap.Error(err),
	)
}

func (l *Logger) LogInfo(message string, value ...interface{}) {
	l.zap.Info(
		message,
		zap.Any(message, value),
	)
}
