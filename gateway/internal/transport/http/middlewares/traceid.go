package middlewares

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

func Trace(c *fiber.Ctx) error {

	t := otel.Tracer("gateway")
	ctx, span := t.Start(c.Context(), "trace")
	defer span.End()

	propagator := propagation.TraceContext{}
	carrier := propagation.MapCarrier{}
	propagator.Inject(ctx, carrier)

	for k, v := range carrier {
		slog.Debug("setting fiber context key", slog.String("key", k), slog.String("value", v))
		c.Set(k, v)
	}

	return c.Next()
}
