package dto

import "github.com/google/uuid"

type Range struct {
	Min int64
	Max int64
}

func (r *Range) Nil() bool {
	return r.Min == 0 && r.Max == 0
}

type AccountFilters struct {
	UserId uuid.UUID
}

func NewAccountFilter() *AccountFilters {
	return &AccountFilters{}
}

func (f *AccountFilters) SetUserId(userId uuid.UUID) *AccountFilters {
	f.UserId = userId
	return f
}
