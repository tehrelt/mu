package models

import (
	"time"

	"github.com/tehrelt/mu/account-service/pkg/pb/accountpb"
)

type Account struct {
	Id        string     `db:"id"`
	UserId    string     `db:"user_id"`
	HouseId   string     `db:"house_id"`
	Balance   int64      `db:"balance"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
}

func (p *Account) ToProto() *accountpb.Account {
	pb := &accountpb.Account{
		Id:        p.Id,
		UserId:    p.UserId,
		HouseId:   p.HouseId,
		Balance:   p.Balance,
		CreatedAt: p.CreatedAt.Unix(),
	}

	if p.UpdatedAt != nil {
		pb.UpdatedAt = p.UpdatedAt.Unix()
	}

	return pb
}
