package models

import (
	"time"

	"github.com/tehrelt/mu/notification-service/pkg/pb/accountpb"
	"github.com/tehrelt/mu/notification-service/pkg/pb/housepb"
)

type Account struct {
	Id        string     `db:"id"`
	UserId    string     `db:"user_id"`
	HouseId   string     `db:"house_id"`
	Balance   int64      `db:"balance"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
}

func (a *Account) DeltaBalance(amount int64) {
	a.Balance += amount
}

func (p *Account) ToProto(house *housepb.House) *accountpb.Account {
	pb := &accountpb.Account{
		Id:        p.Id,
		UserId:    p.UserId,
		Balance:   p.Balance,
		CreatedAt: p.CreatedAt.Unix(),
		UpdatedAt: 0,
	}

	if house != nil {
		pb.House = &accountpb.House{
			Id:                  house.Id,
			Address:             house.Address,
			ConnectedServiceIds: house.ConnectedServices,
		}
	}

	if p.UpdatedAt != nil {
		pb.UpdatedAt = p.UpdatedAt.Unix()
	}

	return pb
}
