package dto

import (
	"github.com/google/uuid"
	"github.com/tehrelt/mu/auth-service/internal/models"
)

type Fio struct {
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	MiddleName string `json:"middle_name"`
}

type Passport struct {
	Number int `json:"number"`
	Series int `json:"series"`
}

type PersonalData struct {
	Phone    string   `json:"phone"`
	Passport Passport `json:"passport"`
	Snils    string   `json:"snils"`
}

type RegisterUser struct {
	UserId   uuid.UUID     `json:"user_id"`
	Password string        `json:"password"`
	Roles    []models.Role `json:"roles"`
}

type LoginUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateUser struct {
	Fio
	PersonalData
	Email string `json:"email"`
}

type UserClaims struct {
	Id string `json:"id"`
}

type Profile struct {
	Id    uuid.UUID `json:"id"`
	Email string    `json:"email"`
	Fio
	Roles []models.Role `json:"roles"`
}
