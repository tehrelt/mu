package middlewares

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/tehrelt/mu-lib/tracer"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
)

func Trace(c *fiber.Ctx) error {

	ctx := c.UserContext()
	t := otel.Tracer(tracer.TracerKey)
	ctx, span := t.Start(ctx, fmt.Sprintf("%s %s", c.Method(), c.Path()))
	defer span.End()

	span.SetAttributes(attribute.String("http.method", c.Method()))

	propagator := propagation.TraceContext{}
	carrier := propagation.MapCarrier{}
	propagator.Inject(ctx, carrier)

	c.SetUserContext(ctx)

	return c.Next()
}
