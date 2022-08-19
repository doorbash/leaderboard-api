package util

import (
	"context"
	"time"
)

const (
	CONTEXT_TIMEOUT = 5 * time.Second
)

func GetContextWithTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, CONTEXT_TIMEOUT)
}

func GetContextWithThisTimeout(ctx context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, timeout)
}
