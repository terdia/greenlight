package entities

import (
	"time"

	"github.com/terdia/greenlight/internal/custom_type"
	"github.com/terdia/greenlight/internal/validator"
)

type Token struct {
	Plaintext string
	Hash      []byte
	UserId    custom_type.ID
	Expiry    time.Time
	Scope     string
}

func (t Token) ValidateTokenPlaintext(v *validator.Validator) {
	v.Check(t.Plaintext != "", "token", "must be provided")
	v.Check(len(t.Plaintext) == 26, "token", "must be 26 bytes long")
}
