package dto

import "github.com/tehrelt/mu/gateway/pkg/pb/ticketpb"

type Ticket interface {
	isTicket()
}

type TicketHeader struct {
	Id           string `json:"id"`
	TicketType   string `json:"ticketType"`
	TicketStatus string `json:"ticketStatus"`
	CreatedBy    string `json:"createdBy"`
}

type ticketNewAccount struct {
	TicketHeader
	Address string `json:"address"`
}

func (t *ticketNewAccount) isTicket() {}

type ticketConnectService struct {
	TicketHeader
	ServiceId string `json:"serviceId"`
	AccountId string `json:"accountId"`
}

func (t *ticketConnectService) isTicket() {}

func marshalHeader(src *ticketpb.TicketHeader) TicketHeader {
	return TicketHeader{
		Id:           src.Id,
		TicketType:   src.Type.String(),
		TicketStatus: src.Status.String(),
		CreatedBy:    src.CreatedBy,
	}
}

func MarshalTicket(header *ticketpb.TicketHeader, payload any) Ticket {

	h := marshalHeader(header)

	switch src := payload.(type) {
	case *ticketpb.Ticket_Account:
		return &ticketNewAccount{
			TicketHeader: h,
			Address:      src.Account.HouseAdress,
		}
	case *ticketpb.Ticket_ConnectService:
		return &ticketConnectService{
			TicketHeader: h,
			ServiceId:    src.ConnectService.ServiceId,
			AccountId:    src.ConnectService.AccountId,
		}
	default:
		return nil
	}
}
