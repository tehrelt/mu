package dto

type Range struct {
	Min int64
	Max int64
}

func (r *Range) Nil() bool {
	return r.Min == 0 && r.Max == 0
}

type AccountFilters struct {
	UserId string
}

func NewAccountFilter() *AccountFilters {
	return &AccountFilters{}
}

func (f *AccountFilters) SetUserId(userId string) *AccountFilters {
	f.UserId = userId
	return f
}
