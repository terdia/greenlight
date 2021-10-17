package dto

import (
	"time"
)

type TokenResponse struct {
	Token Token `json:"authentication_token"`
}

type Token struct {
	PlainText string    `json:"token"`
	Expiry    time.Time `json:"expiry"`
}

type AuthTokenRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
