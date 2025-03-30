package interceptors

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func UnaryServerInterceptor(t trace.Tracer) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		startTime := time.Now()

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		}

		propagator := otel.GetTextMapPropagator()
		ctx = propagator.Extract(ctx, &metadataSupplier{metadata: &md})

		ctx, span := t.Start(
			ctx,
			info.FullMethod,
			trace.WithSpanKind(trace.SpanKindServer),
			trace.WithAttributes(
				attribute.String("rpc.service", info.FullMethod),
			),
		)
		defer span.End()

		resp, err := handler(ctx, req)
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

func StreamServerInterceptor(tracer trace.Tracer) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := ss.Context()
		startTime := time.Now()

		// Extract trace context from incoming metadata
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		}

		// Extract the span context from the metadata
		propagator := otel.GetTextMapPropagator()
		ctx = propagator.Extract(ctx, &metadataSupplier{metadata: &md})

		// Start a new span
		ctx, span := tracer.Start(
			ctx,
			info.FullMethod,
			trace.WithSpanKind(trace.SpanKindServer),
			trace.WithAttributes(
				attribute.String("rpc.system", "grpc"),
				attribute.String("rpc.service", info.FullMethod),
			),
		)
		defer span.End()

		// Wrap the server stream with our context
		wrappedStream := &tracedServerStream{
			ServerStream: ss,
			ctx:          ctx,
			span:         span,
		}

		err := handler(srv, wrappedStream)

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

		return err
	}
}
