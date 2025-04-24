package dto

import (
	"github.com/google/uuid"
	"github.com/tehrelt/mu/auth-service/internal/models"
)

type SaveCredentials struct {
	UserId         uuid.UUID
	HashedPassword string
	Roles          []models.Role
}
