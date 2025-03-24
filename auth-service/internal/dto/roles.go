package dto

import (
	"github.com/google/uuid"
	"github.com/tehrelt/moi-uslugi/auth-service/internal/models"
)

type UserRoles struct {
	UserId uuid.UUID
	Roles  []models.Role
}
