package main

import (
	"go.uber.org/zap"
)

type AppLogger struct {
	Level LoggerLevel
	log   *zap.Logger
}

func NewDefaultAppLogger(level LoggerLevel) (*AppLogger, error) {
	zcore, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	if level == LoggerLevelDebug {
		zcore, err = zap.NewDevelopment()
		if err != nil {
			return nil, err
		}
	}

	logger := zap.New(zcore.Core())

	return &AppLogger{
		Level: level,
		log:   logger,
	}, nil
}

func (l *AppLogger) Info(msg string, fields ...zap.Field) {
	if l.Level != LoggerLevelSilent {
		l.log.Info(msg, fields...)
	}
}

func (l *AppLogger) Error(msg string, fields ...zap.Field) {
	l.log.Error(msg, fields...)
}

func (l *AppLogger) Debug(msg string, fields ...zap.Field) {
	l.log.Debug(msg, fields...)
}

func (l *AppLogger) Panic(msg string, fields ...zap.Field) {
	l.log.Panic(msg, fields...)
}
