package dto

import (
	"github.com/tehrelt/mu/gateway/pkg/pb/authpb"
)

type Role string

func (r Role) Validate() bool {
	switch r {
	case RoleAdmin, RoleRegular:
		return true
	default:
		return false
	}
}

func (r Role) ToProto() authpb.Role {
	switch r {
	case RoleAdmin:
		return authpb.Role_ROLE_ADMIN
	case RoleRegular:
		return authpb.Role_ROLE_REGULAR
	default:
		return authpb.Role_ROLE_UNKNOWN
	}
}

func (r Role) FromProto(role authpb.Role) Role {
	switch role {
	case authpb.Role_ROLE_ADMIN:
		return RoleAdmin
	case authpb.Role_ROLE_REGULAR:
		return RoleRegular
	default:
		return RoleUnknown
	}
}

const (
	RoleAdmin   Role = "admin"
	RoleRegular Role = "regular"
	RoleUnknown Role = "unknown"
)
