package dto

import (
	"github.com/google/uuid"
	"github.com/tehrelt/mu/billing-service/internal/models"
)

type Range struct {
	Min int64
	Max int64
}

func (r *Range) Nil() bool {
	return r.Min == 0 && r.Max == 0
}

type PaymentFilters struct {
	AccountId   uuid.UUID
	Status      models.PaymentStatus
	AmountRange Range
}

func NewPaymentFilter() *PaymentFilters {
	return &PaymentFilters{
		AccountId: uuid.Nil,
		Status:    models.PaymentStatusNil,
	}
}

func (f *PaymentFilters) SetAccountId(id uuid.UUID) *PaymentFilters {
	f.AccountId = id
	return f
}

func (f *PaymentFilters) SetStatus(status models.PaymentStatus) *PaymentFilters {
	f.Status = status
	return f
}

func (f *PaymentFilters) SetAmountMin(min int64) *PaymentFilters {
	f.AmountRange.Min = min
	return f
}

func (f *PaymentFilters) SetAmountMax(max int64) *PaymentFilters {
	f.AmountRange.Max = max
	return f
}
