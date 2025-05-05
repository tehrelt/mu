package models

import (
	"errors"
	"time"

	"github.com/tehrelt/mu/billing-service/pkg/pb/billingpb"
)

type PaymentStatus string

var (
	PaymentStatusNil      PaymentStatus = ""
	PaymentStatusPending  PaymentStatus = "pending"
	PaymentStatusSuccess  PaymentStatus = "paid"
	PaymentStatusCanceled PaymentStatus = "canceled"
)

func (s PaymentStatus) IsValid() bool {
	return s == PaymentStatusPending || s == PaymentStatusSuccess || s == PaymentStatusCanceled
}

func (s PaymentStatus) FromProto(pb billingpb.PaymentStatus) (PaymentStatus, error) {
	switch pb {
	case billingpb.PaymentStatus_pending:
		return PaymentStatusPending, nil
	case billingpb.PaymentStatus_success:
		return PaymentStatusSuccess, nil
	case billingpb.PaymentStatus_canceled:
		return PaymentStatusCanceled, nil
	default:
		return PaymentStatusNil, errors.New("invalid payment status")
	}
}

func (s PaymentStatus) ToProto() billingpb.PaymentStatus {
	switch s {
	case PaymentStatusCanceled:
		return billingpb.PaymentStatus_canceled
	case PaymentStatusSuccess:
		return billingpb.PaymentStatus_success
	default:
		return billingpb.PaymentStatus_pending
	}
}

type Payment struct {
	Id        string        `db:"id"`
	AccountId string        `db:"account_id"`
	Amount    int64         `db:"amount"`
	Status    PaymentStatus `db:"status"`
	Message   string        `db:"message"`
	CreatedAt time.Time     `db:"created_at"`
	UpdatedAt *time.Time    `db:"updated_at"`
}

func (p *Payment) ToProto() *billingpb.Payment {
	pb := &billingpb.Payment{
		Id:        p.Id,
		AccountId: p.AccountId,
		Amount:    p.Amount,
		Message:   p.Message,
		Status:    p.Status.ToProto(),
		CreatedAt: p.CreatedAt.Unix(),
	}

	if p.UpdatedAt != nil {
		pb.UpdatedAt = p.UpdatedAt.Unix()
	}

	return pb
}
