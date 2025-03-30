package tracer

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func UnaryServerInterceptor(t trace.Tracer) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		startTime := time.Now()

		md, _ := metadata.FromIncomingContext(ctx)

		carrier := propagation.HeaderCarrier{}

		for k, v := range md {
			carrier[k] = v
		}

		pg := propagation.TraceContext{}
		ctx = pg.Extract(ctx, carrier)

		ctx, span := t.Start(
			ctx,
			info.FullMethod,
			trace.WithSpanKind(trace.SpanKindServer),
			trace.WithAttributes(
				attribute.String("rpc.service", info.FullMethod),
			),
		)
		defer span.End()

		resp, err := handler(withTracer(ctx, t), req)
		if err != nil {
			if e, ok := status.FromError(err); ok {
				span.SetStatus(codes.Error, e.Message())
			} else {
				span.SetStatus(codes.Error, err.Error())
			}

			span.RecordError(err)
		}

		span.SetAttributes(
			attribute.Int64("rpc.duration_ms", time.Since(startTime).Milliseconds()),
		)

		return resp, err
	}
}
