package logger

import (
	"context"

	"go.uber.org/zap"
)

type key string

const (
	loggerKey key = "logger"
	RequestID key = "request_id"
)

type Logger interface {
	Info(ctx context.Context, msg string, fields ...zap.Field)
	Error(ctx context.Context, msg string, fields ...zap.Field)
	Debug(ctx context.Context, msg string, fields ...zap.Field)
	Fatal(ctx context.Context, msg string, fields ...zap.Field)
	Stop() error
}

type logger struct {
	l *zap.Logger
}

func New() *logger {
	log, err := zap.NewProduction()
	if err != nil {
		return &logger{
			l: zap.L(),
		}
	}

	return &logger{
		l: log,
	}
}

func (l *logger) Stop() error {
	return l.l.Sync()
}

func (l *logger) Info(ctx context.Context, msg string, fields ...zap.Field) {
	fields = addID(ctx, fields)
	l.l.Info(msg, fields...)
}

func (l *logger) Error(ctx context.Context, msg string, fields ...zap.Field) {
	fields = addID(ctx, fields)
	l.l.Error(msg, fields...)
}

func (l *logger) Debug(ctx context.Context, msg string, fields ...zap.Field) {
	fields = addID(ctx, fields)
	l.l.Debug(msg, fields...)
}

func (l *logger) Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	fields = addID(ctx, fields)
	l.l.Fatal(msg, fields...)
}

func WithLogger(ctx context.Context, logger *logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

func FromContext(ctx context.Context) *logger {
	l, ok := ctx.Value(loggerKey).(logger)
	if !ok {
		return New()
	}

	return &l
}

func addID(ctx context.Context, fields []zap.Field) []zap.Field {
	id, ok := ctx.Value(RequestID).(string)
	if ok {
		fields = append(fields, zap.String("request_id", id))
	}

	return fields
}
