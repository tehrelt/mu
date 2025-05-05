package usecase

import (
	"context"

	"github.com/tehrelt/mu/consumption-service/internal/dto"
	"github.com/tehrelt/mu/consumption-service/internal/models"
)

func (uc *UseCase) CreateCabinet(ctx context.Context, in *dto.NewCabinet) (*models.Cabinet, error) {

	if _, err := uc.findAccount(ctx, in.AccountId); err != nil {
		return nil, err
	}

	if _, err := uc.findService(ctx, in.ServiceId); err != nil {
		return nil, err
	}

	cabinet, err := uc.storage.Create(ctx, in)
	if err != nil {
		return nil, err
	}

	return cabinet, nil
}
