package grpc

import (
	"log/slog"

	"github.com/tehrelt/mu-lib/sl"
	"github.com/tehrelt/mu/user-service/internal/models"
	"github.com/tehrelt/mu/user-service/pkg/pb/userpb"
	"google.golang.org/grpc"
)

func (s *Server) List(req *userpb.ListRequest, srv grpc.ServerStreamingServer[userpb.ListResponse]) error {

	filters := &models.UserFilters{
		Limit:  req.Limit,
		Offset: req.Offset,
	}

	usersCh, err := s.users.provider.List(srv.Context(), filters)
	if err != nil {
		return err
	}

	chunkSize := 10

	chunk := make([]*userpb.User, 0, chunkSize)

	for {
		u, ok := <-usersCh
		if !ok {
			slog.Info("users channel closed")
			if len(chunk) > 0 {
				slog.Info("sending last chunk of users", slog.Int("len", len(chunk)))
				if err := srv.Send(&userpb.ListResponse{
					UsersChunk: chunk,
				}); err != nil {
					return err
				}
			}
			break
		}

		slog.Info("appending user to chunk", sl.UUID("userId", u.Id))
		user := &userpb.User{
			Id: u.Id.String(),
			Fio: &userpb.FIO{
				Lastname:   u.LastName,
				Firstname:  u.FirstName,
				Middlename: u.MiddleName,
			},
			Email: u.Email,
			PersonalData: &userpb.PersonalData{
				Passport: &userpb.Passport{
					Series: int32(u.PersonalData.Passport.Series),
					Number: int32(u.PersonalData.Passport.Number),
				},
				Snils: u.PersonalData.Snils,
				Phone: u.PersonalData.Phone,
			},
			CreatedAt: u.CreatedAt.Unix(),
		}

		if u.UpdatedAt != nil {
			user.UpdatedAt = u.UpdatedAt.Unix()
		}

		chunk = append(chunk, user)

		if len(chunk) == chunkSize {
			slog.Info("chunk is full, sending chunk of users", slog.Int("len", len(chunk)))
			if err := srv.Send(&userpb.ListResponse{
				UsersChunk: chunk,
			}); err != nil {
				return err
			}

			chunk = make([]*userpb.User, 0, chunkSize)
		}
	}

	return nil
}
