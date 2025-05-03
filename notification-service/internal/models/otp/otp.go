package otp

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type OTP string

const Nil = ""

func New() OTP {
	return OTP(uuid.New().String())
}

func (otp OTP) Hash(cost int) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(otp), cost)
	if err != nil {
		return nil, err
	}

	return hashedPassword, nil
}
