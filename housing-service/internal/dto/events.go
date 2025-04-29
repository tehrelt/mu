package dto

type EventServiceConnected struct {
	HouseId   string `json:"houseId"`
	ServiceId string `json:"serviceId"`
}

type EventServiceConnect struct {
	HouseId   string `json:"houseId"`
	ServiceId string `json:"serviceId"`
}
