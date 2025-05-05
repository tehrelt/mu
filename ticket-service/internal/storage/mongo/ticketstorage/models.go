package ticketstorage

import (
	"fmt"
	"time"

	"github.com/tehrelt/mu/ticket-service/internal/models"
	"github.com/tehrelt/mu/ticket-service/internal/storage"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type isTicket interface {
	header() *Header
}

type Header struct {
	ID        *primitive.ObjectID `bson:"_id"`
	Type      models.TicketType   `bson:"type"`
	Status    models.TicketStatus `bson:"status"`
	CreatedBy string              `bson:"created_by"`
	CreatedAt time.Time           `bson:"created_at"`
	UpdatedAt time.Time           `bson:"updated_at"`
}

func (th *Header) header() *Header {
	return th
}

type accountTicket struct {
	Header  `bson:",inline"`
	Address string `bson:"address"`
}

type connectServiceTicket struct {
	Header    `bson:",inline"`
	AccountId string `bson:"account_id"`
	ServiceId string `bson:"service_id"`
}

func marshalTicket(ticket models.Ticket) (isTicket, error) {
	header, err := marshalHeader(ticket.Header())
	if err != nil {
		return nil, err
	}

	switch ticket := ticket.(type) {
	case *models.NewAccountTicket:
		return &accountTicket{
			Header:  header,
			Address: ticket.Address,
		}, nil
	case *models.ConnectServiceTicket:
		return &connectServiceTicket{
			Header:    header,
			AccountId: ticket.AccountId,
			ServiceId: ticket.ServiceId,
		}, nil
	default:
		return nil, fmt.Errorf("unsupported ticket type: %T", ticket)
	}
}

func marshalHeader(src *models.TicketHeader) (dst Header, err error) {
	dst = Header{
		Type:      src.TicketType,
		Status:    src.Status,
		CreatedBy: src.CreatedBy,
		CreatedAt: time.Now(),
	}

	if src.Id != "" {
		*dst.ID, err = primitive.ObjectIDFromHex(src.Id)
		if err != nil {
			return Header{}, err
		}
	}

	return dst, nil
}

func unmarshalTicket(src isTicket) (models.Ticket, error) {
	switch ticket := src.(type) {
	case *accountTicket:
		return &models.NewAccountTicket{
			TicketHeader: unmarshalHeader(src.header()),
			Address:      ticket.Address,
		}, nil
	case *connectServiceTicket:
		return &models.ConnectServiceTicket{
			TicketHeader: unmarshalHeader(src.header()),
			AccountId:    ticket.AccountId,
			ServiceId:    ticket.ServiceId,
		}, nil
	default:
		return nil, fmt.Errorf("unsupported ticket type: %T", ticket)
	}
}

func unmarshalHeader(src *Header) models.TicketHeader {
	id := ""
	if src.ID != nil {
		id = src.ID.Hex()
	}

	return models.TicketHeader{
		Id:         id,
		TicketType: src.Type,
		Status:     src.Status,
		CreatedBy:  src.CreatedBy,
	}
}

func factoryTicket(ticketType models.TicketType) (isTicket, error) {
	switch ticketType {
	case models.TicketTypeAccount:
		return &accountTicket{}, nil
	case models.TicketTypeConnectService:
		return &connectServiceTicket{}, nil
	default:
		return nil, storage.ErrInvalidTicketType
	}
}
