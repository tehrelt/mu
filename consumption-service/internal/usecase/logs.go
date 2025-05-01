package usecase

import (
	"context"

	"github.com/tehrelt/mu/consumption-service/internal/dto"
	"github.com/tehrelt/mu/consumption-service/internal/models"
)

func (u *UseCase) Logs(ctx context.Context, filters *dto.LogsFilters) ([]*models.ConsumptionLog, uint64, error) {

	logs, err := u.storage.Logs(ctx, filters)
	if err != nil {
		return nil, 0, err
	}

	total, err := u.storage.CountLogs(ctx, filters)
	if err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}
