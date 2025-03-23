package dto

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

type CreateUser struct {
	Fio
	PersonalData
	Email string `json:"email"`
}
