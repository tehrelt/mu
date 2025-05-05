package models

type TicketFilters struct {
	Status    *TicketStatus
	Type      *TicketType
	UserId    *string
	AccountId *string
	ServiceId *string
}

func NewFilters() *TicketFilters {
	return &TicketFilters{}
}

func (f *TicketFilters) SetStatus(status TicketStatus) *TicketFilters {
	f.Status = &status
	return f
}

func (f *TicketFilters) SetType(t TicketType) *TicketFilters {
	f.Type = &t
	return f
}

func (f *TicketFilters) SetUserId(id string) *TicketFilters {
	f.UserId = &id
	return f
}

func (f *TicketFilters) SetAccountId(id string) *TicketFilters {
	f.AccountId = &id
	return f
}
