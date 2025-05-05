package dto

import (
	"github.com/google/uuid"
	"github.com/tehrelt/mu/account-service/pkg/pb/accountpb"
)

type PaymentStatus string

var (
	PaymentStatusPending  PaymentStatus = "pending"
	PaymentStatusCanceled PaymentStatus = "canceled"
	PaymentStatusPaid     PaymentStatus = "paid"
)

type CreateAccount struct {
	UserId  uuid.UUID
	HouseId uuid.UUID
}

func (cp *CreateAccount) FromProto(pb *accountpb.CreateRequest) (err error) {
	cp.UserId, err = uuid.Parse(pb.UserId)
	if err != nil {
		return err
	}

	cp.HouseId, err = uuid.Parse(pb.HouseId)
	if err != nil {
		return err
	}

	return nil
}

type UpdateAccount struct {
	Id         uuid.UUID
	NewBalance *int64
}

func NewUpdateAccount(id uuid.UUID) *UpdateAccount {
	return &UpdateAccount{
		Id: id,
	}
}

func (dto *UpdateAccount) WithNewBalance(balance int64) *UpdateAccount {
	dto.NewBalance = new(int64)
	*dto.NewBalance = balance
	return dto
}
