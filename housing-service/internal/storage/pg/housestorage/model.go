package housestorage

import (
	"time"

	"github.com/google/uuid"
	"github.com/tehrelt/mu/housing-service/internal/models"
)

type houseDbModel struct {
	ID                string     `db:"id"`
	Address           string     `db:"address"`
	RoomsQuantity     int        `db:"rooms_quantity"`
	ResidentsQuantity int        `db:"residents_quantity"`
	CreatedAt         time.Time  `db:"created_at"`
	UpdatedAt         *time.Time `db:"updated_at"`
}

func (h *houseDbModel) ToHouse() (*models.House, error) {
	id, err := uuid.Parse(h.ID)
	if err != nil {
		return nil, err
	}

	return &models.House{
		Id:                id,
		Address:           h.Address,
		RoomsQuantity:     h.RoomsQuantity,
		ResidentsQuantity: h.ResidentsQuantity,
		CreatedAt:         h.CreatedAt,
		UpdatedAt:         h.UpdatedAt,
	}, nil
}
