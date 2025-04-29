package dto

import "github.com/tehrelt/mu/gateway/pkg/pb/ratepb"

func ParseServiceType(s string) ratepb.ServiceType {
	switch s {
	case ratepb.ServiceType_GAS_SUPPLY.String():
		return ratepb.ServiceType_GAS_SUPPLY
	case ratepb.ServiceType_HEATING.String():
		return ratepb.ServiceType_HEATING
	case ratepb.ServiceType_POWER_SUPPLY.String():
		return ratepb.ServiceType_POWER_SUPPLY
	case ratepb.ServiceType_WATER_SUPPLY.String():
		return ratepb.ServiceType_WATER_SUPPLY
	default:
		return ratepb.ServiceType_UNKNOWN
	}
}
