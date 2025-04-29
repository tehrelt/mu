package dto

import "github.com/tehrelt/mu/gateway/pkg/pb/ticketpb"

type Ticket interface {
	isTicket()
}

type TicketHeader struct {
	Id           string `json:"id"`
	TicketType   string `json:"ticket_type"`
	TicketStatus string `json:"ticket_status"`
}

type ticketNewAccount struct {
	TicketHeader
	Address string `json:"address"`
}

func (t *ticketNewAccount) isTicket() {}

type ticketConnectService struct {
	TicketHeader
	ServiceId string `json:"service_id"`
	UserId    string `json:"user_id"`
	AccountId string `json:"account_id"`
}

func (t *ticketConnectService) isTicket() {}

func marshalHeader(src *ticketpb.TicketHeader) TicketHeader {
	return TicketHeader{
		Id:           src.Id,
		TicketType:   src.Type.String(),
		TicketStatus: src.Status.String(),
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
			UserId:       src.ConnectService.UserId,
			AccountId:    src.ConnectService.AccountId,
		}
	default:
		return nil
	}
}
