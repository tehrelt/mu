package events

type BalanceChanged struct {
	EventHeader
	NewBalance int64 `json:"newBalance"`
}
