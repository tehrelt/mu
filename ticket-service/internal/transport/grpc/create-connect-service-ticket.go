package grpc

import (
	"context"

	"github.com/tehrelt/mu/ticket-service/internal/models"
	"github.com/tehrelt/mu/ticket-service/pkg/pb/ticketpb"
)

// CreateConnectServiceTicket implements ticketpb.TicketServiceServer.
func (s *Server) CreateConnectServiceTicket(ctx context.Context, req *ticketpb.TicketConnectService) (*ticketpb.CreateResponse, error) {
	ticket := models.ConnectServiceTicket{
		TicketHeader: models.TicketHeader{
			TicketType: models.TicketTypeConnectService,
			Status:     models.TicketStatusPending,
		},
		UserId:    req.UserId,
		AccountId: req.AccountId,
		ServiceId: req.ServiceId,
	}

	if err := s.storage.Create(ctx, &ticket); err != nil {
		return nil, err
	}

	return &ticketpb.CreateResponse{
		Id: ticket.Id,
	}, nil
}
