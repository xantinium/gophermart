package logger

import (
	"fmt"

	"go.uber.org/zap"
)

var logger *zap.SugaredLogger

// Init инициализирует логгер.
func Init(isDev bool) {
	var (
		lg  *zap.Logger
		err error
	)

	if isDev {
		lg, err = zap.NewDevelopment()
	} else {
		lg, err = zap.NewProduction()
	}
	if err != nil {
		panic(fmt.Errorf("failed to init logger: %w", err))
	}

	logger = lg.Sugar()
}

// Destroy уничтожает логгер, записывая оставшиеся данные.
func Destroy() {
	logger.Sync()
}

// Field поле, которое можно использовать для стилизации логов.
type Field struct {
	Name  string
	Value any
}

// Info пишет структурированный лог уровня INFO.
func Info(msg string, fields ...Field) {
	log(InfoLevel, msg, fields...)
}

// Infof пишет форматированный лог уровня INFO.
func Infof(format string, args ...any) {
	logger.Infof(format, args)
}

// Warn пишет структурированный лог уровня WARN.
func Warn(msg string, fields ...Field) {
	log(WarnLevel, msg, fields...)
}

// Warnf пишет форматированный лог уровня WARN.
func Warnf(format string, args ...any) {
	logger.Warnf(format, args)
}

// Error пишет структурированный лог уровня ERROR.
func Error(msg string, fields ...Field) {
	log(ErrorLevel, msg, fields...)
}

// Errorf пишет форматированный лог уровня ERROR.
func Errorf(format string, args ...any) {
	logger.Errorf(format, args)
}

// LogLevel уровень логирования.
type LogLevel = uint8

const (
	// InfoLevel уровень нейтральных уведомлений.
	InfoLevel LogLevel = iota
	// WarnLevel уровень подозрительных событий.
	WarnLevel
	// ErrorLevel уровень ошибок.
	ErrorLevel
)

func log(level LogLevel, msg string, fields ...Field) {
	args := make([]any, len(fields)*2)

	for i, field := range fields {
		args[2*i] = field.Name
		args[2*i+1] = field.Value
	}

	switch level {
	case InfoLevel:
		logger.Infow(msg, args...)
	case WarnLevel:
		logger.Warnw(msg, args...)
	case ErrorLevel:
		logger.Errorw(msg, args...)
	}
}
