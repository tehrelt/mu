package dto

import (
	"github.com/google/uuid"
	"github.com/tehrelt/mu/auth-service/internal/models"
)

type UserRoles struct {
	UserId uuid.UUID
	Roles  []models.Role
}
