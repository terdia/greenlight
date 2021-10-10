package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/terdia/greenlight/internal/data"
	"github.com/terdia/greenlight/src/users/entities"
	"github.com/terdia/greenlight/src/users/repositories"
)

const (
	queryTimeout = 3 * time.Second
)

type userRepository struct {
	sql.DB
}

func NewUserRepoitory(db *sql.DB) repositories.UserRepository {
	return &userRepository{*db}
}

func (repo *userRepository) Insert(user *entities.User) error {
	query := `
			INSERT INTO users (name, email, password_hash, activated)
			VALUES ($1, $2, $3, $4)
			RETURNING id, created_at, version`

	queryParams := []interface{}{user.Name, user.Email, user.Password.Hash, user.Activated}

	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	err := repo.QueryRowContext(ctx, query, queryParams...).Scan(&user.ID, &user.CreatedAt, &user.Version)
	if err != nil {
		return err
	}

	return nil
}

func (repo *userRepository) GetByEmail(email string) (*entities.User, error) {

	query := `
			SELECT id, created_at, name, email, password_hash, activated, version
			FROM users
			WHERE email = $1`

	var user entities.User

	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	err := repo.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Name,
		&user.Email,
		&user.Password.Hash,
		&user.Activated,
		&user.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, data.ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (repo *userRepository) Update(user *entities.User) error {
	query := `
			UPDATE users
			SET name = $1, email = $2, password_hash = $3, activated = $4, version = version + 1
			WHERE id = $5 AND version = $6
			RETURNING version`

	args := []interface{}{user.Name, user.Email, user.Password.Hash, user.Activated, user.ID, user.Version}

	// Execute the SQL query. If no matching row could be found, we know the user
	// version has changed (or the record has been deleted) and we return our custom
	// ErrEditConflict error.
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)

	defer cancel()

	err := repo.DB.QueryRowContext(ctx, query, args...).Scan(&user.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return data.ErrDuplicateEmail
		case errors.Is(err, sql.ErrNoRows):
			return data.ErrEditConflict
		default:
			return err
		}
	}

	return nil
}
