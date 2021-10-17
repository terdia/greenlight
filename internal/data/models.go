package data

import (
	"errors"
)

var (
	ErrRecordNotFound     = errors.New("models: record not found")
	ErrEditConflict       = errors.New("models: edit conflict")
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	ErrDuplicateEmail     = errors.New("models: a user with this email address already exists")
)

const (
	TokenScopeActivation     = "activation"
	TokenScopeAuthentication = "authentication"
)
