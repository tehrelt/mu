package usecase

import (
	"context"

	"github.com/tehrelt/mu/consumption-service/internal/dto"
	"github.com/tehrelt/mu/consumption-service/internal/models"
)

func (uc *UseCase) FindCabinet(ctx context.Context, criteria *dto.FindCabinet) (*models.Cabinet, error) {
	return uc.storage.Find(ctx, criteria)
}
