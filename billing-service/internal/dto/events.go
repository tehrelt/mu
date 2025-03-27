package dto

import "github.com/tehrelt/mu/billing-service/internal/models"

type EventStatusChanged struct {
	AccountId string               `json:"accountId"`
	PaymentId string               `json:"paymentId"`
	NewStatus models.PaymentStatus `json:"newStatus"`
}
