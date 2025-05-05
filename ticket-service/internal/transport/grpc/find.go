package grpc

import (
	"context"

	"github.com/tehrelt/mu/ticket-service/internal/transport/grpc/marshaler"
	"github.com/tehrelt/mu/ticket-service/pkg/pb/ticketpb"
)

// Find implements ticketpb.TicketServiceServer.
func (s *Server) Find(ctx context.Context, req *ticketpb.FindRequest) (*ticketpb.FindResponse, error) {
	ticket, err := s.storage.Find(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	marshaled := marshaler.MarshalTicket(ticket)

	return &ticketpb.FindResponse{Ticket: marshaled}, nil
}
