package grpc

import (
	"context"

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

	return &ticketpb.UpdateTicketStatusResponse{
		Id: req.Id,
	}, nil
}
