package models

import "github.com/google/uuid"

type EventConnectService struct {
	HouseId   uuid.UUID
	ServiceId uuid.UUID
}

func ParseEventConnectService(houseId, serviceId string) (*EventConnectService, error) {

	hId, err := uuid.Parse(houseId)
	if err != nil {
		return nil, err
	}

	sId, err := uuid.Parse(serviceId)
	if err != nil {
		return nil, err
	}

	return &EventConnectService{
		HouseId:   hId,
		ServiceId: sId,
	}, nil
}
