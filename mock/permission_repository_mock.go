package mock

import (
	"github.com/terdia/greenlight/internal/custom_type"
	"github.com/terdia/greenlight/internal/data"
	"github.com/terdia/greenlight/src/users/repositories"
)

type permissionRepositoryMock struct{}

func NewPermissionRepositoryMock() repositories.PermissionRepository {
	return &permissionRepositoryMock{}
}

func (p *permissionRepositoryMock) GetAllForUser(userID custom_type.ID) (data.Permissions, error) {

	var permissions data.Permissions

	return permissions, nil
}

func (p *permissionRepositoryMock) AddForUser(userID custom_type.ID, codes ...string) error {

	return nil
}
