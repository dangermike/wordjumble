package logging

import (
	"context"

	"go.uber.org/zap"
)

var loggerKey struct{}

func FromContext(ctx context.Context) *zap.Logger {
	if t, ok := ctx.Value(loggerKey).(*zap.Logger); ok {
		return t
	}
	return zap.L()
}

func NewContext(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}
