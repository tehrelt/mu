package models

import "github.com/google/uuid"

type UserClaims struct {
	Id    string   `json:"id"`
	Roles []string `json:"roles"`
}

type Fio struct {
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	MiddleName string `json:"middleName"`
}

type Passport struct {
	Series int `json:"series"`
	Number int `json:"number"`
}

type PersonalData struct {
	Phone    string   `json:"phone"`
	Snils    string   `json:"snils"`
	Passport Passport `json:"passport"`
}

type RegisterUserRequest struct {
	Fio
	PersonalData
	Email    string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	Id    uuid.UUID `json:"id"`
	Email string    `json:"email"`
	Fio
	PersonalData
}

type UserRole struct {
	Id uuid.UUID `json:"id"`
}

type Credentials struct {
	UserId         uuid.UUID `json:"userId"`
	HashedPassword string    `json:"hashedPassword"`
	Roles          []Role    `json:"roles"`
}
