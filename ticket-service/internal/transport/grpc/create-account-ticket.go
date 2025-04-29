package grpc

import (
	"context"

	"github.com/tehrelt/mu/ticket-service/internal/models"
	"github.com/tehrelt/mu/ticket-service/pkg/pb/ticketpb"
)

// CreateAccountTicket implements ticketpb.TicketServiceServer.
func (s *Server) CreateAccountTicket(ctx context.Context, req *ticketpb.NewTicketAccount) (*ticketpb.CreateResponse, error) {
	ticket := models.NewAccountTicket{
		TicketHeader: models.TicketHeader{
			TicketType: models.TicketTypeAccount,
			Status:     models.TicketStatusPending,
			CreatedBy:  req.CreatedBy,
		},
		Address: req.HouseAdress,
	}

	if err := s.storage.Create(ctx, &ticket); err != nil {
		return nil, err
	}

	return &ticketpb.CreateResponse{
		Id: ticket.Id,
	}, nil
}
