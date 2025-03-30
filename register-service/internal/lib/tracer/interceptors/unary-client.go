package interceptors

import (
	"context"
	"time"

	"github.com/tehrelt/mu/register-service/internal/lib/tracer"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func UnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		start := time.Now()

		t := otel.Tracer(tracer.TracerKey)
		ctx, span := t.Start(ctx, method)
		defer span.End()

		md, ok := metadata.FromOutgoingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		}

		propagator := otel.GetTextMapPropagator()
		propagator.Inject(ctx, &metadataSupplier{metadata: &md})

		ctx = metadata.NewOutgoingContext(ctx, md)

		err := invoker(ctx, method, req, reply, cc, opts...)
		if err != nil {
			span.RecordError(err)
		}

		span.SetAttributes(
			attribute.String("rpc.method", method),
			attribute.String("rpc.system", "grpc"),
			attribute.String("rpc.peer_address", cc.Target()),
			attribute.Int64("rpc.duration_ms", time.Since(start).Milliseconds()),
		)

		return err
	}
}

// StreamClientInterceptor returns a stream client interceptor for OpenTelemetry tracing
func StreamClientInterceptor() grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		start := time.Now()

		ctx, span := otel.Tracer(tracer.TracerKey).Start(
			ctx,
			method,
			trace.WithSpanKind(trace.SpanKindClient),
		)

		md, ok := metadata.FromOutgoingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		}

		propagator := otel.GetTextMapPropagator()
		propagator.Inject(ctx, &metadataSupplier{metadata: &md})

		ctx = metadata.NewOutgoingContext(ctx, md)

		stream, err := streamer(ctx, desc, cc, method, opts...)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			span.End()
			return nil, err
		}

		span.SetAttributes(
			attribute.String("rpc.method", method),
			attribute.String("rpc.system", "grpc"),
			attribute.String("rpc.peer_address", cc.Target()),
			attribute.Int64("rpc.duration_ms", time.Since(start).Milliseconds()),
		)

		return &tracingClientStream{
			ClientStream: stream,
			span:         span,
		}, nil
	}
}

// tracingClientStream wraps grpc.ClientStream to properly close the span
type tracingClientStream struct {
	grpc.ClientStream
	span trace.Span
}

func (s *tracingClientStream) RecvMsg(m interface{}) error {
	err := s.ClientStream.RecvMsg(m)
	if err != nil {
		s.span.RecordError(err)
		s.span.SetStatus(codes.Error, err.Error())
	}
	return err
}

func (s *tracingClientStream) SendMsg(m interface{}) error {
	err := s.ClientStream.SendMsg(m)
	if err != nil {
		s.span.RecordError(err)
		s.span.SetStatus(codes.Error, err.Error())
	}
	return err
}

func (s *tracingClientStream) CloseSend() error {
	err := s.ClientStream.CloseSend()
	if err != nil {
		s.span.RecordError(err)
		s.span.SetStatus(codes.Error, err.Error())
	}
	return err
}

func (s *tracingClientStream) Header() (metadata.MD, error) {
	md, err := s.ClientStream.Header()
	if err != nil {
		s.span.RecordError(err)
		s.span.SetStatus(codes.Error, err.Error())
	}
	return md, err
}

func (s *tracingClientStream) Trailer() metadata.MD {
	return s.ClientStream.Trailer()
}

func (s *tracingClientStream) Context() context.Context {
	return s.ClientStream.Context()
}

// metadataSupplier implements the TextMapCarrier interface for metadata propagation
type metadataSupplier struct {
	metadata *metadata.MD
}

func (s *metadataSupplier) Get(key string) string {
	values := s.metadata.Get(key)
	if len(values) == 0 {
		return ""
	}
	return values[0]
}

func (s *metadataSupplier) Set(key, value string) {
	s.metadata.Set(key, value)
}

func (s *metadataSupplier) Keys() []string {
	out := make([]string, 0, len(*s.metadata))
	for key := range *s.metadata {
		out = append(out, key)
	}
	return out
}
