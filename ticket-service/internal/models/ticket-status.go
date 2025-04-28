package models

import "github.com/tehrelt/mu/ticket-service/pkg/pb/ticketpb"

type TicketStatus string

const (
	TicketStatusNil      TicketStatus = ""
	TicketStatusPending  TicketStatus = "pending"
	TicketStatusApproved TicketStatus = "approved"
	TicketStatusRejected TicketStatus = "rejected"
)

func (s TicketStatus) IsValid() bool {
	return s == TicketStatusPending || s == TicketStatusApproved || s == TicketStatusRejected
}

func (s TicketStatus) ToProto() ticketpb.TicketStatus {
	switch s {
	case TicketStatusPending:
		return ticketpb.TicketStatus_TicketStatusPending
	case TicketStatusApproved:
		return ticketpb.TicketStatus_TicketStatusApproved
	case TicketStatusRejected:
		return ticketpb.TicketStatus_TicketStatusRejected
	default:
		return ticketpb.TicketStatus_TicketStatusUnknown
	}
}
