package dto

import (
	"time"

	"github.com/terdia/greenlight/internal/custom_type"
)

type CreateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SingleUserResponse struct {
	User UserResponse `json:"user"`
}

type UserResponse struct {
	ID        custom_type.ID `json:"id"`
	Name      string         `json:"name"`
	Email     string         `json:"email"`
	CreatedAt time.Time      `json:"created_at"`
	Version   int            `json:"version"`
}
