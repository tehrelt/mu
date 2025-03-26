package models

import (
	"time"

	"github.com/google/uuid"
)

type Service struct {
	Id          string     `db:"id"`
	Name        string     `db:"s_name"`
	MeasureUnit string     `db:"measure_unit"`
	Rate        int64      `db:"rate"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   *time.Time `db:"updated_at"`
}

type CreateService struct {
	Name        string `db:"name"`
	MeasureUnit string `db:"measure_unit"`
	Rate        int64  `db:"rate"`
}

type UpdateServiceRate struct {
	Id   uuid.UUID `db:"id"`
	Rate int64     `db:"rate"`
}
