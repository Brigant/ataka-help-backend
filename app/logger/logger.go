package logger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger represents logger.
type Logger struct{ log *zap.SugaredLogger }

// New initialize logger.
func New(logLevel string) (*Logger, error) {
	level, err := zapcore.ParseLevel(logLevel)
	if err != nil {
		return nil, fmt.Errorf("error with logger level parsing: %w", err)
	}

	cfg := zap.Config{
		Encoding:          "console",
		Level:             zap.NewAtomicLevelAt(level),
		Development:       zap.NewDevelopmentConfig().Development,
		DisableCaller:     zap.NewDevelopmentConfig().DisableCaller,
		DisableStacktrace: zap.NewDevelopmentConfig().DisableStacktrace,
		Sampling:          nil,
		InitialFields:     zap.NewDevelopmentConfig().InitialFields,
		OutputPaths:       []string{"stderr"},
		ErrorOutputPaths:  []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			NameKey: "logger",
			// StacktraceKey:  "stacktrace",
			LevelKey:       "level",
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			MessageKey:     "message",
			TimeKey:        "time",
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeDuration: zapcore.SecondsDurationEncoder,
		},
	}

	logger, err := cfg.Build()
	if err != nil {
		return nil, fmt.Errorf("can't build: %w", err)
	}

	return &Logger{logger.Sugar()}, nil
}

// Flush will flush any buffered log entries.
func (l *Logger) Flush() {
	if err := l.log.Sync(); err != nil {
		l.log.Error(err)
	}
}

// Methods above will implement all needful logging behavior.
func (l *Logger) Errorf(msg string, val ...any) {
	l.log.Errorf(msg, val...)
}

func (l *Logger) Errorw(msg string, val ...any) {
	l.log.Errorw(msg, val...)
}

func (l *Logger) Debugf(msg string, val ...any) {
	l.log.Debugf(msg, val...)
}

func (l *Logger) Debugw(msg string, val ...any) {
	l.log.Debugw(msg, val...)
}

func (l *Logger) Infof(msg string, val ...any) {
	l.log.Infof(msg, val...)
}

func (l *Logger) Warnf(msg string, val ...any) {
	l.log.Warnf(msg, val...)
}

func (l *Logger) Warnw(msg string, val ...any) {
	l.log.Warnw(msg, val...)
}

func (l *Logger) Infow(msg string, val ...any) {
	l.log.Infow(msg, val...)
}
