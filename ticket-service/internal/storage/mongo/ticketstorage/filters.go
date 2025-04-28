package ticketstorage

import "github.com/tehrelt/mu/ticket-service/internal/models"

type filters struct {
	Status    *models.TicketStatus `bson:"status,omitempty"`
	Type      *models.TicketType   `bson:"type,omitempty"`
	UserId    *string              `bson:"user_id,omitempty"`
	AccountId *string              `bson:"account_id,omitempty"`
	ServiceId *string              `bson:"service_id,omitempty"`
}

func marshalFilters(f *models.TicketFilters) *filters {
	return &filters{
		Status:    f.Status,
		Type:      f.Type,
		UserId:    f.UserId,
		AccountId: f.AccountId,
		ServiceId: f.ServiceId,
	}
}
