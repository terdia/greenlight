package repository

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"time"

	"github.com/terdia/greenlight/internal/data"
	"github.com/terdia/greenlight/src/users/entities"
	"github.com/terdia/greenlight/src/users/repositories"
)

const (
	QueryTimeout = 3 * time.Second
)

type userRepository struct {
	*sql.DB
}

func NewUserRepoitory(db *sql.DB) repositories.UserRepository {
	return &userRepository{db}
}

func (repo *userRepository) Insert(user *entities.User) error {
	query := `
			INSERT INTO users (name, email, password_hash, activated)
			VALUES ($1, $2, $3, $4)
			RETURNING id, created_at, version`

	queryParams := []interface{}{user.Name, user.Email, user.Password.Hash, user.Activated}

	ctx, cancel := context.WithTimeout(context.Background(), QueryTimeout)
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

	ctx, cancel := context.WithTimeout(context.Background(), QueryTimeout)
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
	ctx, cancel := context.WithTimeout(context.Background(), QueryTimeout)

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

func (repo *userRepository) GetForToken(tokenPlainText, scope string) (*entities.User, error) {

	hash := sha256.Sum256([]byte(tokenPlainText))

	query := `
			SELECT users.id, users.created_at, users.name, users.email, 
			users.password_hash, users.activated, users.version
			FROM users
			INNER JOIN tokens
			ON users.id = tokens.user_id
			WHERE tokens.hash = $1
			AND tokens.scope = $2
			AND tokens.expiry > $3`

	args := []interface{}{hash[:], scope, time.Now()}

	ctx, cancel := context.WithTimeout(context.Background(), QueryTimeout)
	defer cancel()

	var user entities.User

	err := repo.DB.QueryRowContext(ctx, query, args...).Scan(
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
