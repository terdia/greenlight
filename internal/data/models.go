package data

import (
	"errors"
)

var (
	ErrRecordNotFound     = errors.New("record not found")
	ErrEditConflict       = errors.New("edit conflict")
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	ErrDuplicateEmail     = errors.New("a user with this email address already exists")
)
