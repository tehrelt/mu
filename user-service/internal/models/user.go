package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id           uuid.UUID
	FirstName    string
	LastName     string
	MiddleName   string
	Email        string
	PersonalData PersonalData
	CreatedAt    time.Time
	UpdatedAt    *time.Time
}

type Passport struct {
	Series int
	Number int
}

type PersonalData struct {
	Phone    string
	Passport Passport
	Snils    string
}

type CreateUser struct {
	FirstName    string
	LastName     string
	MiddleName   string
	Email        string
	PersonalData PersonalData
}

type UserFilters struct {
	Limit  uint64
	Offset uint64
}
