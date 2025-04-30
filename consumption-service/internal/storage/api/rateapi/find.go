package rateapi

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/tehrelt/mu-lib/sl"
	"github.com/tehrelt/mu/consumption-service/internal/models"
	"github.com/tehrelt/mu/consumption-service/pkg/pb/ratepb"
)

// Find implements usecase.ServiceProvider.
func (a *Api) Find(ctx context.Context, id uuid.UUID) (*models.Service, error) {
	fn := "Find"
	logger := a.logger.With(sl.Method(fn))

	req := &ratepb.FindRequest{
		Id: id.String(),
	}

	res, err := a.client.Find(ctx, req)
	if err != nil {
		logger.Error("failed find rate", sl.Err(err))
		return nil, err
	}
	logger.Debug("found rate", slog.Any("rate", res))

	return &models.Service{
		Id:   uuid.MustParse(res.Id),
		Rate: uint64(res.Rate),
	}, nil
}
