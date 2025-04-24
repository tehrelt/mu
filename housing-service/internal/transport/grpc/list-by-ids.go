package grpc

import (
	"errors"
	"io"

	"github.com/google/uuid"
	"github.com/tehrelt/mu/housing-service/pkg/pb/housepb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ housepb.HouseServiceServer

func (s *Server) ListHousesByIds(srv grpc.BidiStreamingServer[housepb.ListHousesByIdsRequest, housepb.ListHousesResponse]) error {

	ids := make([]uuid.UUID, 0)

	for {
		req, err := srv.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return err
		}

		id, err := uuid.Parse(req.Id)
		if err != nil {
			return status.Errorf(codes.InvalidArgument, "invalid id = %s", req.Id)
		}

		ids = append(ids, id)
	}

	if len(ids) == 0 {
		return status.Errorf(codes.InvalidArgument, "no ids provided")
	}

	hh, err := s.storage.ListByIds(srv.Context(), ids)
	if err != nil {
		return err
	}

	for _, h := range hh {

		hresp := &housepb.House{
			Id:           h.Id.String(),
			Address:      h.Address,
			RoomsQty:     int64(h.RoomsQuantity),
			ResidentsQty: int64(h.ResidentsQuantity),
			CreatedAt:    h.CreatedAt.Unix(),
		}

		if h.UpdatedAt != nil {
			hresp.UpdatedAt = h.UpdatedAt.Unix()
		}

		if err := srv.Send(&housepb.ListHousesResponse{
			House: hresp,
		}); err != nil {
			return err
		}
	}

	return nil
}
