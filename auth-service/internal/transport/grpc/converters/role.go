package converters

import (
	"github.com/tehrelt/mu/auth-service/internal/models"
	"github.com/tehrelt/mu/auth-service/pkg/pb/authpb"
)

func RoleFromPb(r authpb.Role) models.Role {
	switch r {
	case authpb.Role_ROLE_ADMIN:
		return models.Role_Admin
	default:
		return models.Role_Regular
	}
}

func RoleToPb(r models.Role) authpb.Role {
	switch r {
	case models.Role_Admin:
		return authpb.Role_ROLE_ADMIN
	default:
		return authpb.Role_ROLE_REGULAR
	}
}
