package models

type Role string

const (
	Role_Regular Role = "regular"
	Role_Admin   Role = "admin"
)

func (r Role) Valid() bool {
	switch r {
	case Role_Regular:
		return true
	case Role_Admin:
		return true
	}

	return false
}
