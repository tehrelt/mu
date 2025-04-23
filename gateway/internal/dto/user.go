package dto

import (
	"github.com/google/uuid"
	"github.com/tehrelt/mu/gateway/pkg/pb/authpb"
)

type UserProfile struct {
	Id         uuid.UUID
	LastName   string
	FirstName  string
	MiddleName string
	Email      string
	Roles      []Role
}

func (u *UserProfile) FromProto(p *authpb.ProfileResponse) (err error) {
	u.Id, err = uuid.Parse(p.Id)
	if err != nil {
		return
	}

	u.LastName = p.LastName
	u.FirstName = p.FirstName
	u.MiddleName = p.MiddleName
	u.Email = p.Email

	u.Roles = make([]Role, len(p.Roles))
	for i, role := range p.Roles {
		u.Roles[i] = u.Roles[i].FromProto(role)
	}

	return
}
