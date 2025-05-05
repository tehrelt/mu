package events

type IncomingBalanceChanged struct {
	AccountId  string `json:"accountId"`
	NewBalance int64  `json:"newBalance"`
	OldBalance int64  `json:"oldBalance"`
	Reason     string `json:"reason"`
}

type BalanceChanged struct {
	EventHeader
	Address    string `json:"address"`
	NewBalance int64  `json:"newBalance"`
	OldBalance int64  `json:"oldBalance"`
	Reason     string `json:"reason"`
}
