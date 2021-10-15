package dto

import (
	"time"

	"github.com/terdia/greenlight/internal/custom_type"
)

type CreateUserRequest struct {
	Name     string `json:"name"`     // fullname
	Email    string `json:"email"`    // unique email address
	Password string `json:"password"` // minimum 8 bytes maximum 72 bytes
}

type SingleUserResponse struct {
	User UserResponse `json:"user"`
}

type ActivateUserRequest struct {
	TokenPlaintext string `json:"token"`
}

type UserResponse struct {
	ID        custom_type.ID `json:"id"`
	Name      string         `json:"name"`
	Email     string         `json:"email"`
	Activated bool           `json:"activated"`
	CreatedAt time.Time      `json:"created_at"`
	Version   int            `json:"version"`
}
