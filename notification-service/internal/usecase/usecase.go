package usecase

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/tehrelt/mu-lib/sl"
	"github.com/tehrelt/mu/notification-service/internal/models"
	"github.com/tehrelt/mu/notification-service/internal/storage/pg/integrationstorage"
	"github.com/tehrelt/mu/notification-service/internal/storage/redis/otpstorage"
)

type UseCase struct {
	otpstorage         *otpstorage.Storage
	integrationstorage *integrationstorage.Storage
	logger             *slog.Logger
}

func New(otpstorage *otpstorage.Storage, integrationstorage *integrationstorage.Storage) *UseCase {
	return &UseCase{
		otpstorage:         otpstorage,
		integrationstorage: integrationstorage,
		logger:             slog.With(sl.Module("use_case")),
	}
}

func (uc *UseCase) Find(ctx context.Context, userId uuid.UUID) (*models.Integration, error) {
	return uc.integrationstorage.Find(ctx, userId)
}
