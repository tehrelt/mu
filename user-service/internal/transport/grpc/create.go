package grpc

import (
	"context"
	"log/slog"

	"github.com/tehrelt/moi-uslugi/user-service/internal/models"
	"github.com/tehrelt/moi-uslugi/user-service/pkg/pb/userspb"
	"github.com/tehrelt/moi-uslugi/user-service/pkg/sl"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Create implements userspb.UserServiceServer.
func (s *Server) Create(ctx context.Context, req *userspb.CreateRequest) (*userspb.CreateResponse, error) {

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
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &userspb.CreateResponse{
		Id: id.String(),
	}, nil
}
