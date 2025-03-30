package grpc

import (
	"context"
	"errors"
	"log/slog"

	"github.com/tehrelt/mu/user-service/internal/models"
	"github.com/tehrelt/mu/user-service/internal/storage"
	"github.com/tehrelt/mu/user-service/pkg/pb/userpb"
	"github.com/tehrelt/mu/user-service/pkg/sl"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Create implements userspb.UserServiceServer.
func (s *Server) Create(ctx context.Context, req *userpb.CreateRequest) (*userpb.CreateResponse, error) {

	log := slog.With(sl.Method("Create"))

	candidate := &models.CreateUser{
		FirstName:  req.Fio.Firstname,
		LastName:   req.Fio.Lastname,
		MiddleName: req.Fio.Middlename,
		Email:      req.Email,
		PersonalData: models.PersonalData{
			Phone: req.PersonalData.Phone,
			Snils: req.PersonalData.Snils,
			Passport: models.Passport{
				Series: int(req.PersonalData.Passport.Series),
				Number: int(req.PersonalData.Passport.Number),
			},
		},
	}

	id, err := s.users.creator.Create(ctx, candidate)
	if err != nil {
		log.Error("failed to create user", sl.Err(err))

		if errors.Is(err, storage.ErrUserAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &userpb.CreateResponse{
		Id: id.String(),
	}, nil
}
