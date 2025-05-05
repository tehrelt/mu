package grpc

import (
	"context"

	"github.com/google/uuid"
	"github.com/tehrelt/mu/notification-service/internal/dto"
	"github.com/tehrelt/mu/notification-service/internal/models/otp"
	"github.com/tehrelt/mu/notification-service/pkg/pb/notificationpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// LinkTelegram implements notificationlpb.NotificationServiceServer.
func (s *Server) LinkTelegram(ctx context.Context, req *notificationpb.LinkTelegramRequest) (*notificationpb.LinkTelegramResponse, error) {

	userid, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user_id=%s", req.UserId)
	}

	in := &dto.LinkTelegram{
		UserId:         userid,
		Code:           otp.OTP(req.Otp),
		TelegramChatId: req.ChatId,
	}

	if err := s.uc.LinkTelegram(ctx, in); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to link telegram: %v", err)
	}

	return &notificationpb.LinkTelegramResponse{}, nil
}

// TelegramOtp implements notificationlpb.NotificationServiceServer.
func (s *Server) TelegramOtp(ctx context.Context, req *notificationpb.TelegramOtpRequest) (*notificationpb.TelegramOtpResponse, error) {
	userid, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user_id=%s", req.UserId)
	}

	otp, err := s.uc.CreateTelegramOtp(ctx, userid)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create telegram otp: %v", err)
	}

	return &notificationpb.TelegramOtpResponse{Otp: string(otp)}, nil
}

// UnlinkTelegram implements notificationlpb.NotificationServiceServer.
func (s *Server) UnlinkTelegram(ctx context.Context, req *notificationpb.UnlinkTelegramRequest) (*notificationpb.UnlinkTelegramResponse, error) {
	panic("unimplemented")
}
