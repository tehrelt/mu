package dto

import "time"

type Payment struct {
	Id        string     `json:"id"`
	Status    string     `json:"status"`
	Amount    float64    `json:"amount"`
	Message   string     `json:"message"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
}
