package dto

import "time"

type Payment struct {
	Id        string     `json:"id"`
	Status    string     `json:"status"`
	Amount    float64    `json:"amount"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
}
