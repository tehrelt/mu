package dto

import "github.com/google/uuid"

type EventPaymentStatusChanged struct {
	AccountId uuid.UUID     `json:"accountId"`
	PaymentId uuid.UUID     `json:"paymentId"`
	NewStatus PaymentStatus `json:"newStatus"`
}

type EventBalanceChanged struct {
	AccountId  string `json:"accountId"`
	NewBalance int64  `json:"newBalance"`
	OldBalance int64  `json:"oldBalance"`
	Reason     string `json:"reason"`
}
