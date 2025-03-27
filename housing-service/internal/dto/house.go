package dto

import "github.com/google/uuid"

type ConnectService struct {
	HouseId   uuid.UUID
	ServiceId uuid.UUID
}

type CreateHouse struct {
	Address     string
	ResidentQty int
	RoomsQty    int
}
