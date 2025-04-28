package models

import (
	"github.com/tehrelt/mu/ticket-service/pkg/pb/ticketpb"
)

type TicketType string

const (
	TicketTypeNil            TicketType = ""
	TicketTypeAccount        TicketType = "account"
	TicketTypeConnectService TicketType = "connect_service"
)

func (t *TicketType) ToProto() ticketpb.TicketType {
	switch *t {
	case TicketTypeAccount:
		return ticketpb.TicketType_TicketTypeAccount
	case TicketTypeConnectService:
		return ticketpb.TicketType_TicketTypeConnectService
	default:
		return ticketpb.TicketType_TicketTypeUnknown
	}
}
