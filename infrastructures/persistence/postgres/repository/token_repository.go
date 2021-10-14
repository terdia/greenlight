package repository

import (
	"context"
	"database/sql"

	"github.com/terdia/greenlight/src/users/entities"
	"github.com/terdia/greenlight/src/users/repositories"
)

type tokenRepository struct {
	*sql.DB
}

func NewTokenRepository(db *sql.DB) repositories.TokenRepository {
	return &tokenRepository{db}
}

func (repo *tokenRepository) Create(token *entities.Token) error {

	query := `
			INSERT INTO tokens (hash, user_id, expiry, scope)
			VALUES ($1, $2, $3, $4)`

	args := []interface{}{token.Hash, token.UserId, token.Expiry, token.Scope}

	ctx, cancel := context.WithTimeout(context.Background(), QueryTimeout)
	defer cancel()

	_, err := repo.DB.ExecContext(ctx, query, args...)

	return err
}

func (repo *tokenRepository) DeleteAllForUserByScope(scope string, userID int64) error {

	query := `
			DELETE FROM tokens
			WHERE scope = $1 AND user_id = $2`

	ctx, cancel := context.WithTimeout(context.Background(), QueryTimeout)
	defer cancel()

	_, err := repo.DB.ExecContext(ctx, query, scope, userID)

	return err
}
