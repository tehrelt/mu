package usecase

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/tehrelt/mu-lib/sl"
	"github.com/tehrelt/mu/notification-service/internal/models"
	"github.com/tehrelt/mu/notification-service/internal/storage/pg/integrationstorage"
	"github.com/tehrelt/mu/notification-service/internal/storage/redis/otpstorage"
	"github.com/tehrelt/mu/notification-service/internal/storage/rmq"
	"github.com/tehrelt/mu/notification-service/pkg/pb/accountpb"
	"github.com/tehrelt/mu/notification-service/pkg/pb/ticketpb"
	"github.com/tehrelt/mu/notification-service/pkg/pb/userpb"
)

type UseCase struct {
	otpstorage         *otpstorage.Storage
	integrationstorage *integrationstorage.Storage
	logger             *slog.Logger

	broker     *rmq.Broker
	ticketapi  ticketpb.TicketServiceClient
	userapi    userpb.UserServiceClient
	accountapi accountpb.AccountServiceClient
}

func New(
	otpstorage *otpstorage.Storage,
	integrationstorage *integrationstorage.Storage,
	ticketapi ticketpb.TicketServiceClient,
	broker *rmq.Broker,
	userapi userpb.UserServiceClient,
	accountapi accountpb.AccountServiceClient,
) *UseCase {
	return &UseCase{
		otpstorage:         otpstorage,
		integrationstorage: integrationstorage,
		logger:             slog.With(sl.Module("use_case")),
		broker:             broker,
		ticketapi:          ticketapi,
		userapi:            userapi,
		accountapi:         accountapi,
	}
}

func (uc *UseCase) Find(ctx context.Context, userId uuid.UUID) (*models.Integration, error) {
	return uc.integrationstorage.Find(ctx, userId)
}
