package usecase

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/tehrelt/mu-lib/sl"
	"github.com/tehrelt/mu/consumption-service/internal/dto"
	"github.com/tehrelt/mu/consumption-service/internal/models"
)

//go:generate mockgen -destination=mocks/mock_service_provider.go -package=mocks github.com/tehrelt/mu/consumption-service/internal/usecase ServiceProvider
type ServiceProvider interface {
	Find(ctx context.Context, id uuid.UUID) (*models.Service, error)
}

type AccountProvider interface {
	Find(ctx context.Context, id uuid.UUID) (*models.Account, error)
}

//go:generate mockgen -destination=mocks/mock_consumption_storage.go -package=mocks github.com/tehrelt/mu/consumption-service/internal/usecase ConsumptionStorage
type ConsumptionStorage interface {
	Create(ctx context.Context, in *dto.NewCabinet) (*models.Cabinet, error)
	Find(ctx context.Context, criteria *dto.FindCabinet) (*models.Cabinet, error)
	Update(ctx context.Context, in *dto.UpdateCabinet) (*models.Cabinet, error)
	Log(ctx context.Context, log *dto.NewConsumeLog) (*models.ConsumptionLog, error)
}

//go:generate mockgen -destination=mocks/mock_billing_provider.go -package=mocks github.com/tehrelt/mu/consumption-service/internal/usecase BillingProvider
type BillingProvider interface {
	Charge(ctx context.Context, bill *dto.Charge) (uuid.UUID, error)
}

type UseCase struct {
	storage         ConsumptionStorage
	serviceProvider ServiceProvider
	billingProvider BillingProvider
	accountProvider AccountProvider
	logger          *slog.Logger
}

func New(storage ConsumptionStorage, serviceProvider ServiceProvider, billingProvider BillingProvider, accountProvider AccountProvider) *UseCase {
	return &UseCase{
		storage:         storage,
		serviceProvider: serviceProvider,
		billingProvider: billingProvider,
		accountProvider: accountProvider,
		logger:          slog.With(sl.Method("usecase")),
	}
}

func (uc *UseCase) find(ctx context.Context, in *dto.FindCabinet) (*models.Cabinet, error) {
	fn := "usecase.find"
	logger := uc.logger.With(sl.Method(fn))

	logger.Debug("try to found consumption", slog.Any("in", in))
	consumption, err := uc.storage.Find(ctx, in)
	if err != nil {
		logger.Error("unexpected error", sl.Err(err))
		return nil, err
	}

	if consumption == nil {
		logger.Error("consumption not found")
		return nil, ErrConsumptionNotFound
	}

	logger.Debug("consumption found", slog.Any("consumption", consumption))
	return consumption, nil
}

func (uc *UseCase) findService(ctx context.Context, id uuid.UUID) (*models.Service, error) {
	service, err := uc.serviceProvider.Find(ctx, id)
	if err != nil {
		return nil, err
	}
	if service == nil {
		return nil, ErrServiceNotFound
	}

	return service, nil
}

func (uc *UseCase) findAccount(ctx context.Context, id uuid.UUID) (*models.Account, error) {
	fn := "findAccount"
	logger := uc.logger.With(sl.Method(fn))

	logger.Debug("try to found account", slog.String("id", id.String()))
	account, err := uc.accountProvider.Find(ctx, id)
	if err != nil {
		return nil, err
	}
	if account == nil {
		return nil, ErrAccountNotFound
	}

	return account, nil
}
