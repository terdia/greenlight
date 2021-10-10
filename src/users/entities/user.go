package entities

import (
	"time"

	"github.com/terdia/greenlight/internal/custom_type"
	"github.com/terdia/greenlight/internal/validator"
)

type User struct {
	ID        custom_type.ID
	CreatedAt time.Time
	Name      string
	Email     string
	Password  Password
	Activated bool
	Version   int
}

type Password struct {
	PlainText *string
	Hash      []byte
}

func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "must be a valid email address")
}

func ValidatePasswordPlaintext(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 8, "password", "must be at least 8 bytes long")
	v.Check(len(password) <= 72, "password", "must not be more than 72 bytes long")
}

func (user *User) ValidateRequest(v *validator.Validator) {
	v.Check(user.Name != "", "name", "must be provided")
	v.Check(len(user.Name) <= 500, "name", "must not be more than 500 bytes long")

	ValidateEmail(v, user.Email)

	if user.Password.PlainText != nil {
		ValidatePasswordPlaintext(v, *user.Password.PlainText)
	}

	if user.Password.Hash == nil {
		panic("missing password hash for user")
	}
}
