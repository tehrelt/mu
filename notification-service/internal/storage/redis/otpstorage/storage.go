package otpstorage

import (
	"context"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type Storage struct {
	client *redis.Client
}

func NewStorage(client *redis.Client) *Storage {
	return &Storage{
		client: client,
	}
}

func (s *Storage) Set(ctx context.Context, userId uuid.UUID, hashedOtp string) error {
	return s.client.Set(ctx, userId.String(), hashedOtp, 0).Err()
}

func (s *Storage) Get(ctx context.Context, userId uuid.UUID) (string, error) {
	res, err := s.client.Get(ctx, userId.String()).Result()
	if err != nil {
		return "", err
	}

	return res, nil
}

func (s *Storage) Delete(ctx context.Context, userId string) error {
	return s.client.Del(ctx, userId).Err()
}
