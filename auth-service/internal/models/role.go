package models

import "github.com/tehrelt/mu/auth-service/pkg/pb/authpb"

type Role string

const (
	Role_Regular Role = "regular"
	Role_Admin   Role = "admin"
	Role_Unknown Role = "unknown"
)

func (r Role) Valid() bool {
	switch r {
	case Role_Regular:
		return true
	case Role_Admin:
		return true
	}

	return false
}

func (role Role) FromProto(r authpb.Role) Role {
	switch r {
	case authpb.Role_ROLE_ADMIN:
		role = Role_Admin
	case authpb.Role_ROLE_REGULAR:
		role = Role_Regular
	default:
		role = Role_Unknown
	}
	return role
}
