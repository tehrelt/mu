package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/tehrelt/mu/rate-service/pkg/pb/ratepb"
)

type ServiceType string

const (
	ServiceTypeUnknown     ServiceType = "unknown"
	ServiceTypeWaterSupply ServiceType = "water_supply"
	ServiceTypeHeating     ServiceType = "heating"
	ServiceTypePowerSupply ServiceType = "power_supply"
	ServiceTypeGasSupply   ServiceType = "gas_supply"
)

func ServiceTypeFromProto(p ratepb.ServiceType) (s ServiceType) {
	switch p {
	case ratepb.ServiceType_WATER_SUPPLY:
		s = ServiceTypeWaterSupply
	case ratepb.ServiceType_HEATING:
		s = ServiceTypeHeating
	case ratepb.ServiceType_POWER_SUPPLY:
		s = ServiceTypePowerSupply
	case ratepb.ServiceType_GAS_SUPPLY:
		s = ServiceTypeGasSupply
	default:
		s = ServiceTypeUnknown
	}

	return s
}

func (s ServiceType) ToProto() ratepb.ServiceType {
	switch s {
	case ServiceTypeWaterSupply:
		return ratepb.ServiceType_WATER_SUPPLY
	case ServiceTypeHeating:
		return ratepb.ServiceType_HEATING
	case ServiceTypePowerSupply:
		return ratepb.ServiceType_POWER_SUPPLY
	case ServiceTypeGasSupply:
		return ratepb.ServiceType_GAS_SUPPLY
	default:
		return ratepb.ServiceType_UNKNOWN
	}
}

type Service struct {
	Id          string      `db:"id"`
	Name        string      `db:"s_name"`
	MeasureUnit string      `db:"measure_unit"`
	Type        ServiceType `db:"s_type"`
	Rate        int64       `db:"rate"`
	CreatedAt   time.Time   `db:"created_at"`
	UpdatedAt   *time.Time  `db:"updated_at"`
}

type CreateService struct {
	Name        string      `db:"name"`
	MeasureUnit string      `db:"measure_unit"`
	Rate        int64       `db:"rate"`
	Type        ServiceType `db:"s_type""`
}

type UpdateServiceRate struct {
	Id   uuid.UUID `db:"id"`
	Rate int64     `db:"rate"`
}

type EventRateChanged struct {
	Id        uuid.UUID `json:"id"`
	OldRate   int64     `json:"oldRate"`
	NewRate   int64     `json:"newRate"`
	Timestamp time.Time `json:"timestamp"`
}

type RateFilters struct {
	Type *ServiceType
}

func NewRateFilters() *RateFilters {
	return &RateFilters{}
}

func (f *RateFilters) WithType(t ServiceType) *RateFilters {
	f.Type = &t
	return f
}
