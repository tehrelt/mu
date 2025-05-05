package grpc

import (
	"github.com/tehrelt/mu/ticket-service/internal/models"
	"github.com/tehrelt/mu/ticket-service/internal/transport/grpc/marshaler"
	"github.com/tehrelt/mu/ticket-service/pkg/pb/ticketpb"
	"google.golang.org/grpc"
)

// List implements ticketpb.TicketServiceServer.
func (s *Server) List(req *ticketpb.ListRequest, stream grpc.ServerStreamingServer[ticketpb.Ticket]) error {

	filters := models.NewFilters()

	if req.UserId != "" {
		filters = filters.SetUserId(req.UserId)
	}

	if req.AccountId != "" {
		filters = filters.SetAccountId(req.AccountId)
	}

	if req.Type != ticketpb.TicketType_TicketTypeUnknown {
		filters = filters.SetType(models.TicketTypeFromProto(req.Type))
	}

	if req.Status != ticketpb.TicketStatus_TicketStatusUnknown {
		filters = filters.SetStatus(models.TicketStatusFromProto(req.Status))
	}

	ticketsCh, err := s.storage.List(stream.Context(), filters)
	if err != nil {
		return err
	}

	for ticket := range ticketsCh {
		t := marshaler.MarshalTicket(ticket)

		if err := stream.Send(t); err != nil {
			return err
		}
	}

	return nil

}
