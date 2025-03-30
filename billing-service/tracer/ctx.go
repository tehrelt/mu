package tracer

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

type ctxKey string

const (
	key ctxKey = "tracer"
)

func withTracer(ctx context.Context, t trace.Tracer) context.Context {
	return context.WithValue(ctx, key, t)
}

func FromContext(ctx context.Context) trace.Tracer {
	return ctx.Value(key).(trace.Tracer)
}
