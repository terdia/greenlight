package entities

import (
	"time"

	"github.com/terdia/greenlight/internal/custom_type"
)

type Token struct {
	Plaintext string
	Hash      []byte
	UserId    custom_type.ID
	Expiry    time.Time
	Scope     string
}
