package middlewares

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/tehrelt/mu-lib/tracer"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

func Trace(api string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()
		t := otel.Tracer(tracer.TracerKey)
		ctx, span := t.Start(ctx, fmt.Sprintf("%s %s", c.Method(), c.Path()), trace.WithAttributes())
		defer span.End()

		span.SetAttributes(
			attribute.String("http.api", api),
			attribute.String("http.path", c.Path()),
			attribute.String("http.method", c.Method()),
		)

		propagator := propagation.TraceContext{}
		carrier := propagation.MapCarrier{}
		propagator.Inject(ctx, carrier)

		c.SetUserContext(ctx)

		return c.Next()
	}
}
