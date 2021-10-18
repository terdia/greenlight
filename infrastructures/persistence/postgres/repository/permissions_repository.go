package repository

import (
	"context"
	"database/sql"

	"github.com/lib/pq"

	"github.com/terdia/greenlight/internal/custom_type"
	"github.com/terdia/greenlight/internal/data"
	"github.com/terdia/greenlight/src/users/repositories"
)

type permissionRepository struct {
	*sql.DB
}

func NewPermissionRepository(db *sql.DB) repositories.PermissionRepository {
	return &permissionRepository{DB: db}
}

func (p *permissionRepository) GetAllForUser(userID custom_type.ID) (data.Permissions, error) {

	query := `
			SELECT permissions.code
			FROM permissions
			INNER JOIN users_permissions ON users_permissions.permission_id = permissions.id
			INNER JOIN users ON users_permissions.user_id = users.id
			WHERE users.id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), QueryTimeout)
	defer cancel()

	rows, err := p.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var permissions data.Permissions

	for rows.Next() {
		var permission string

		err := rows.Scan(&permission)
		if err != nil {
			return nil, err
		}

		permissions = append(permissions, permission)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return permissions, nil
}

func (p *permissionRepository) AddForUser(userID custom_type.ID, codes ...string) error {
	query := `
		INSERT INTO users_permissions
		SELECT $1, permissions.id FROM permissions WHERE permissions.code = ANY($2)`

	ctx, cancel := context.WithTimeout(context.Background(), QueryTimeout)
	defer cancel()

	_, err := p.DB.ExecContext(ctx, query, userID, pq.Array(codes))

	return err
}
