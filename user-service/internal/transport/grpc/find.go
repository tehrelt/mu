package grpc

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/tehrelt/moi-uslugi/user-service/internal/models"
	"github.com/tehrelt/moi-uslugi/user-service/pkg/pb/userspb"
	"github.com/tehrelt/moi-uslugi/user-service/pkg/sl"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Find implements userspb.UserServiceServer.
func (s *Server) Find(ctx context.Context, req *userspb.FindRequest) (*userspb.FindResponse, error) {

	log := slog.With(sl.Method("Find"))

	var user *models.User
	var err error

	// Handle the oneof field
	switch req.SearchBy.(type) {
	case *userspb.FindRequest_Id:
		rawId := req.GetId()

		id, err := uuid.Parse(rawId)
		if err != nil {
			log.Info("failed to parse uuid", slog.String("rawId", rawId), sl.Err(err))
			return nil, status.Error(codes.InvalidArgument, "invalid id")
		}

		user, err = s.users.provider.UserById(ctx, id)
	case *userspb.FindRequest_Email:
		email := req.GetEmail()
		// TODO ADD EMAIL VALIDATION
		user, err = s.users.provider.UserByEmail(ctx, email)
	default:
		return nil, status.Error(codes.InvalidArgument, "either id or email must be provided")
	}

	if err != nil {
		log.Error("failed to find user", sl.Err(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	response := &userspb.FindResponse{
		User: &userspb.User{
			Id:    user.Id.String(),
			Email: user.Email,
			Fio: &userspb.FIO{
				Lastname:   user.LastName,
				Firstname:  user.FirstName,
				Middlename: user.MiddleName,
			},
			PersonalData: &userspb.PersonalData{
				Snils: user.PersonalData.Snils,
				Phone: user.PersonalData.Phone,
				Passport: &userspb.Passport{
					Series: int32(user.PersonalData.Passport.Series),
					Number: int32(user.PersonalData.Passport.Number),
				},
			},
			CreatedAt: user.CreatedAt.Unix(),
		},
	}

	if user.UpdatedAt != nil {
		response.User.UpdatedAt = user.UpdatedAt.Unix()
	}

	return response, nil
}
