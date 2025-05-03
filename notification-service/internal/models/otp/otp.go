package otp

import (
	"math/rand"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type OTP string

const Nil = ""

func generateChar() string {
	return string(rune(rand.Intn(26) + 'a'))
}

func generateWord(len int) string {
	chars := make([]string, len)
	for i := range chars {
		chars[i] = generateChar()
	}
	return strings.Join(chars, "")
}

func New() OTP {
	return OTP(generateWord(6))
}

func (otp OTP) Hash(cost int) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(otp), cost)
	if err != nil {
		return nil, err
	}

	return hashedPassword, nil
}
