package dto

type EventServiceConnected struct {
	AccountId string `json:"accountId"`
	ServiceId string `json:"serviceId"`
}

type EventServiceConnect struct {
	AccountId string `json:"accountId"`
	HouseId   string `json:"houseId"`
	ServiceId string `json:"serviceId"`
}
