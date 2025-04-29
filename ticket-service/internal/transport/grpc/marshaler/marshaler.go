package marshaler

import (
	"github.com/tehrelt/mu/ticket-service/internal/models"
	"github.com/tehrelt/mu/ticket-service/pkg/pb/ticketpb"
)

type ispayload interface {
	isTicket_Payload()
}

func MarshalTicket(ticket models.Ticket) *ticketpb.Ticket {

	res := &ticketpb.Ticket{
		Header: marshalHeader(ticket.Header()),
	}

	switch ticket.Header().TicketType {
	case models.TicketTypeAccount:
		res.Payload = marshalAccountPayload(ticket.(*models.NewAccountTicket))
	case models.TicketTypeConnectService:
		res.Payload = marshalConnectServicePayload(ticket.(*models.ConnectServiceTicket))
	default:
		return nil
	}

	return res
}

func marshalHeader(header *models.TicketHeader) *ticketpb.TicketHeader {
	return &ticketpb.TicketHeader{
		Id:        header.Id,
		Type:      header.TicketType.ToProto(),
		Status:    header.Status.ToProto(),
		CreatedBy: header.CreatedBy,
	}
}

func marshalAccountPayload(ticket *models.NewAccountTicket) *ticketpb.Ticket_Account {
	return &ticketpb.Ticket_Account{
		Account: &ticketpb.TicketAccount{
			HouseAdress: ticket.Address,
		},
	}
}

func marshalConnectServicePayload(ticket *models.ConnectServiceTicket) *ticketpb.Ticket_ConnectService {
	return &ticketpb.Ticket_ConnectService{
		ConnectService: &ticketpb.TicketConnectService{
			AccountId: ticket.AccountId,
			ServiceId: ticket.ServiceId,
		},
	}
}

func UnmarshalTicketStatus(src ticketpb.TicketStatus) models.TicketStatus {
	switch src {
	case ticketpb.TicketStatus_TicketStatusPending:
		return models.TicketStatusPending
	case ticketpb.TicketStatus_TicketStatusApproved:
		return models.TicketStatusApproved
	case ticketpb.TicketStatus_TicketStatusRejected:
		return models.TicketStatusRejected
	default:
		return models.TicketStatusNil
	}
}
