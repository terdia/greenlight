package repositories

import (
	"github.com/terdia/greenlight/internal/custom_type"
	"github.com/terdia/greenlight/internal/data"
)

type PermissionRepository interface {
	GetAllForUser(userID custom_type.ID) (data.Permissions, error)
	AddForUser(userID custom_type.ID, codes ...string) error
}
