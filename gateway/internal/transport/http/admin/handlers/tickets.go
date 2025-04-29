package handlers

import (
	"io"
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/tehrelt/mu/gateway/internal/transport/http"
	"github.com/tehrelt/mu/gateway/internal/transport/http/public/handlers/dto"
	"github.com/tehrelt/mu/gateway/pkg/pb/ticketpb"
	"go.opentelemetry.io/otel/trace"
)

type NewAccountRequest struct {
	Address string `json:"address"`
}

type NewTicketResponse struct {
	Id string `json:"id"`
}

type ListResponse struct {
	Tickets []dto.Ticket `json:"tickets"`
}

func TicketListHandler(ticketer ticketpb.TicketServiceClient) fiber.Handler {
	return func(c *fiber.Ctx) error {

		ctx := c.UserContext()

		span := trace.SpanFromContext(ctx)

		stream, err := ticketer.List(ctx, &ticketpb.ListRequest{})
		if err != nil {
			return err
		}

		res := &ListResponse{
			Tickets: make([]dto.Ticket, 0, 4),
		}
		for {
			chunk, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}

			t := dto.MarshalTicket(chunk.Header, chunk.Payload)
			if t == nil {
				slog.Error("failed to marshal ticket", slog.Any("chunk", chunk))
				span.AddEvent("failed marshal ticket", trace.WithAttributes(
					http.JsonAttribute("invalid_ticket", t),
				))
				continue
			}

			res.Tickets = append(res.Tickets, t)
		}

		return c.JSON(res)
	}
}
