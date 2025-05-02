package dto

import (
	"github.com/google/uuid"
	"github.com/tehrelt/mu/billing-service/internal/models"
	"github.com/tehrelt/mu/billing-service/pkg/pb/billingpb"
)

type CreatePayment struct {
	AccountId uuid.UUID
	Amount    int64
	Message   string
}

func (cp *CreatePayment) FromProto(pb *billingpb.CreateRequest) (err error) {
	cp.AccountId, err = uuid.Parse(pb.AccountId)
	if err != nil {
		return
	}

	cp.Amount = pb.Amount
	cp.Message = pb.Message

	return nil
}

type UpdatePayment struct {
	Id     uuid.UUID
	Amount *int64
	Status *models.PaymentStatus
}

func NewUpdatePayment(id uuid.UUID) *UpdatePayment {
	return &UpdatePayment{
		Id: id,
	}
}

func (u *UpdatePayment) SetAmount(amount int64) *UpdatePayment {
	u.Amount = new(int64)
	*u.Amount = amount
	return u
}

func (u *UpdatePayment) SetStatus(status models.PaymentStatus) *UpdatePayment {
	u.Status = new(models.PaymentStatus)
	*u.Status = status
	return u
}
