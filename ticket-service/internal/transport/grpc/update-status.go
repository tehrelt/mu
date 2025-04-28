package grpc

import (
	"context"

	"github.com/tehrelt/mu/ticket-service/pkg/pb/ticketpb"
)

func (s *Server) UpdateTicketStatus(context.Context, *ticketpb.UpdateTicketStatusRequest) (*ticketpb.UpdateTicketStatusResponse, error) {
	panic("unimplemented")
}
