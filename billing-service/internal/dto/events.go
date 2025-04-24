package dto

import (
	"time"

	"github.com/tehrelt/mu/billing-service/internal/models"
)

type EventStatusChanged struct {
	AccountId string               `json:"accountId"`
	PaymentId string               `json:"paymentId"`
	NewStatus models.PaymentStatus `json:"newStatus"`
	Timestamp time.Time            `json:"timestamp"`
}
