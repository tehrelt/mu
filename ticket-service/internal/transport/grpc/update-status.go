package grpc

import (
	"context"
	"log/slog"
	"time"

	"github.com/tehrelt/mu-lib/sl"
	"github.com/tehrelt/mu/ticket-service/internal/dto"
	"github.com/tehrelt/mu/ticket-service/internal/models"
	"github.com/tehrelt/mu/ticket-service/internal/transport/grpc/marshaler"
	"github.com/tehrelt/mu/ticket-service/pkg/pb/ticketpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) UpdateTicketStatus(ctx context.Context, req *ticketpb.UpdateTicketStatusRequest) (*ticketpb.UpdateTicketStatusResponse, error) {
	st := marshaler.UnmarshalTicketStatus(req.Status)
	if !st.IsValid() {
		return nil, status.Errorf(codes.InvalidArgument, "invalid status")
	}

	if err := s.storage.Update(ctx, req.Id, st); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update ticket status")
	}

	t, err := s.storage.Find(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to find ticket")
	}

	event := &dto.EventTicketStatusChanged{
		TicketId:  req.Id,
		Ticket:    t,
		Status:    st,
		Timestamp: time.Now(),
	}

	if t.Header().TicketType == models.TicketTypeAccount {
		if err := s.broker.PublishStatusNewAccount(ctx, event); err != nil {
			slog.Error("failed to publish status changed event", sl.Err(err))
			return nil, status.Errorf(codes.Internal, "failed to publish status changed event")
		}
	} else if t.Header().TicketType == models.TicketTypeConnectService {
		if err := s.broker.PublishStatusConnectService(ctx, event); err != nil {
			slog.Error("failed to publish status changed event", sl.Err(err))
			return nil, status.Errorf(codes.Internal, "failed to publish status changed event")
		}
	}

	return &ticketpb.UpdateTicketStatusResponse{
		Id: req.Id,
	}, nil
}
