package usecase

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/tehrelt/mu-lib/sl"
	"github.com/tehrelt/mu/notification-service/internal/dto"
	"github.com/tehrelt/mu/notification-service/internal/models"
	"github.com/tehrelt/mu/notification-service/internal/models/otp"
	"golang.org/x/crypto/bcrypt"
)

func (u *UseCase) CreateTelegramOtp(ctx context.Context, userId uuid.UUID) (otp.OTP, error) {
	fn := "CreateTelegramOtp"
	log := u.logger.With(sl.Method(fn))

	res, err := u.otpstorage.Get(ctx, userId)
	if err != nil {
		log.Error("failed to get otp", sl.Err(err))
	}
	log.Debug("try to get otp", slog.String("res", res))

	code := otp.New()
	log.Debug("created new otp", slog.String("code", string(code)))
	hash, err := code.Hash(10)
	if err != nil {
		log.Error("failed to hash otp", sl.Err(err))
		return otp.Nil, err
	}

	log.Debug("set hashed otp to user", slog.String("userId", userId.String()), slog.String("hash", string(hash)))
	if err = u.otpstorage.Set(ctx, userId, string(hash)); err != nil {
		return otp.Nil, err
	}

	return code, nil
}

func (u *UseCase) LinkTelegram(ctx context.Context, in *dto.LinkTelegram) error {
	fn := "LinkTelegram"
	log := u.logger.With(sl.Method(fn))

	user, err := u.integrationstorage.Find(ctx, in.UserId)
	if err != nil {
		log.Error("failed to find integration", sl.Err(err))
		return err
	}

	if user == nil {
		log.Info("user not found, creating new integration", slog.String("userId", in.UserId.String()))
		user = &models.Integration{
			UserId: in.UserId,
		}

		if err := u.integrationstorage.Create(ctx, user); err != nil {
			log.Error("failed to create user", sl.Err(err))
			return err
		}
	}

	hash, err := u.otpstorage.Get(ctx, in.UserId)
	if err != nil {
		log.Error("failed to get hashed otp", sl.Err(err))
		return err
	}
	log.Debug("get hash from otp store", slog.String("hash", hash))

	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(in.Code)); err != nil {
		log.Error("failed to compare hash", sl.Err(err))
		return err
	}
	log.Debug("otp code verified", slog.String("code", string(in.Code)))
	log.Debug("set telegram chat id for user", slog.String("userId", in.UserId.String()), slog.String("telegramChatId", in.TelegramChatId))
	user.SetTelegramChatId(in.TelegramChatId)
	if err := u.integrationstorage.Update(ctx, user); err != nil {
		log.Error("failed to update user", sl.Err(err))
		return err
	}

	log.Debug("delete otp from redis")
	if err := u.otpstorage.Delete(ctx, in.UserId.String()); err != nil {
		log.Error("failed to delete otp", sl.Err(err))
		return err
	}

	return nil
}
