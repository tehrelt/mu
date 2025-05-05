package usecase

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/tehrelt/mu-lib/sl"
	"github.com/tehrelt/mu/notification-service/internal/dto"
	"github.com/tehrelt/mu/notification-service/internal/events"
	"github.com/tehrelt/mu/notification-service/internal/models"
	"github.com/tehrelt/mu/notification-service/internal/usecase/sender"
	"github.com/tehrelt/mu/notification-service/pkg/pb/ticketpb"
	"github.com/tehrelt/mu/notification-service/pkg/pb/userpb"
)

func (uc *UseCase) HandleTicketStatusChangedEvent(ctx context.Context, event *events.IncomingTicketStatusChanged) error {

	fn := "HandleTicketStatusChangedEvent"
	log := uc.logger.With(sl.Method(fn))

	log.Info("looking for ticket", slog.String("ticketId", event.TicketId))
	ticket, err := uc.ticketapi.Find(ctx, &ticketpb.FindRequest{
		Id: event.TicketId,
	})
	if err != nil {
		return err
	}

	userId := ticket.Ticket.Header.CreatedBy

	upcomingEvent := &events.TicketStatusChanged{
		EventHeader: events.NewEventHeader(events.EventTicketStatusChanged, userId),
		TicketId:    event.TicketId,
		NewStatus:   event.Status,
	}

	slog.Info("composing new event", slog.Any("event", upcomingEvent))
	if err := uc.send(ctx, upcomingEvent); err != nil {
		return err
	}

	return nil
}

func (uc *UseCase) send(ctx context.Context, event events.Event) error {
	fn := "send"
	log := uc.logger.With(sl.Method(fn))

	userId, err := uuid.Parse(event.Header().UserId)
	if err != nil {
		return fmt.Errorf("invalid user id: %s", event.Header().UserId)
	}

	log.Info("sending event", slog.Any("event", event))

	settings, err := uc.userSettings(ctx, userId)
	if err != nil {
		return err
	}

	event.SetSettings(settings)

	sender := sender.New(uc.broker, settings)

	if err := sender.Send(ctx, event); err != nil {
		return err
	}

	return nil
}

func (uc *UseCase) userSettings(ctx context.Context, userId uuid.UUID) (*dto.UserSettings, error) {

	user, err := uc.userapi.Find(ctx, &userpb.FindRequest{
		SearchBy: &userpb.FindRequest_Id{
			Id: userId.String(),
		},
	})
	if err != nil {
		return nil, err
	}

	settings, err := uc.integrationstorage.Find(ctx, userId)
	if err != nil {
		return nil, err
	}
	if settings == nil {
		settings = &models.Integration{
			UserId: userId,
		}

		if err := uc.integrationstorage.Create(ctx, settings); err != nil {
			slog.Error("failed to create integration", slog.String("user_id", userId.String()), sl.Err(err))
			return nil, err
		}
	}

	return &dto.UserSettings{
		Email:          user.User.Email,
		TelegramChatId: settings.TelegramChatId,
	}, nil
}
