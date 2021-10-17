package main

import (
	"context"
	"net/http"

	"github.com/terdia/greenlight/src/users/entities"
)

type contextKey string

const (
	userContextKey = contextKey("user")
)

func (app *application) contextSetUser(r *http.Request, user *entities.User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)

	return r.WithContext(ctx)
}

func (app *application) contextGetUser(r *http.Request) *entities.User {
	user, ok := r.Context().Value(userContextKey).(*entities.User)

	if !ok {
		panic("missing user value in context")
	}

	return user
}
