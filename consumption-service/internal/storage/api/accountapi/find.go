package accountapi

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/tehrelt/mu-lib/sl"
	"github.com/tehrelt/mu/consumption-service/internal/models"
	"github.com/tehrelt/mu/consumption-service/pkg/pb/accountpb"
)

func (a *Api) Find(ctx context.Context, id uuid.UUID) (*models.Account, error) {
	fn := "Find"
	logger := a.logger.With(sl.Method(fn))

	req := &accountpb.FindRequest{
		Id: id.String(),
	}

	res, err := a.client.Find(ctx, req)
	if err != nil {
		logger.Error("failed find account", sl.Err(err))
		return nil, err
	}
	logger.Debug("found account", slog.Any("account", res))

	return &models.Account{
		Id: uuid.MustParse(res.Account.Id),
	}, nil
}
