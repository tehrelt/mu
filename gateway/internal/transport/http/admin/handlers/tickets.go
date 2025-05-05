package handlers

import (
	"io"
	"log/slog"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/tehrelt/mu/gateway/internal/dto"
	"github.com/tehrelt/mu/gateway/internal/transport/http"
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

		req := &ticketpb.ListRequest{}

		status := c.Query("status", "")
		if status != "" {
			for k, v := range ticketpb.TicketStatus_value {
				if strings.Compare(k, status) == 0 {
					req.Status = ticketpb.TicketStatus(v)
				}
			}
		}

		ttype := c.Query("type", "")
		if ttype != "" {
			for k, v := range ticketpb.TicketType_value {
				if strings.Compare(k, ttype) == 0 {
					req.Type = ticketpb.TicketType(v)
				}
			}
		}

		stream, err := ticketer.List(ctx, req)
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

func TicketDetailsHandler(ticketer ticketpb.TicketServiceClient) fiber.Handler {
	return func(c *fiber.Ctx) error {

		ctx := c.UserContext()

		id := c.Params("id")
		if id == "" {
			return fiber.ErrBadRequest
		}

		ticket, err := ticketer.Find(ctx, &ticketpb.FindRequest{Id: id})
		if err != nil {
			return err
		}

		return c.JSON(dto.MarshalTicket(ticket.Ticket.Header, ticket.Ticket.Payload))
	}
}

type TicketStatusPatchRequest struct {
	Status string `json:"status"`
}

func TicketStatusPatchHandler(ticketer ticketpb.TicketServiceClient) fiber.Handler {
	return func(c *fiber.Ctx) error {

		ctx := c.UserContext()

		id := c.Params("id")
		if id == "" {
			return fiber.ErrBadRequest
		}

		var req TicketStatusPatchRequest
		if err := c.BodyParser(&req); err != nil {
			return err
		}

		status := dto.ParseTicketStatus(req.Status)

		if _, err := ticketer.UpdateTicketStatus(ctx, &ticketpb.UpdateTicketStatusRequest{
			Id:     id,
			Status: status,
		}); err != nil {
			return err
		}

		return c.SendStatus(fiber.StatusNoContent)
	}
}
