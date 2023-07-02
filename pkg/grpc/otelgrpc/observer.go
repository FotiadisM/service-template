package otelgrpc

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
)

type ctxKey struct{}

type observer struct {
	skip bool

	msgSentCount    int
	msgReceiveCount int

	// isStreaming is used to avoid measuring duration in streaming RPCs
	isStreaming bool

	attrs []attribute.KeyValue
}

func observerFromCtx(ctx context.Context) *observer {
	o, ok := ctx.Value(ctxKey{}).(*observer)
	if !ok {
		return &observer{}
	}
	return o
}

func ctxWithObserver(ctx context.Context, o *observer) context.Context {
	return context.WithValue(ctx, ctxKey{}, o)
}
