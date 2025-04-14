package dto

import "github.com/tehrelt/mu/gateway/pkg/pb/authpb"

type Role string

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

const (
	RoleAdmin   Role = "admin"
	RoleRegular Role = "regular"
)
