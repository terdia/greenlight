package mock

import (
	"github.com/terdia/greenlight/internal/custom_type"
	"github.com/terdia/greenlight/src/users/entities"
	tr "github.com/terdia/greenlight/src/users/repositories"
)

type tokenRepositoryMock struct{}

func NewTokenRepositoryMock() tr.TokenRepository {
	return &tokenRepositoryMock{}
}

func (repo *tokenRepositoryMock) Create(token *entities.Token) error {

	return nil
}

func (repo *tokenRepositoryMock) DeleteAllForUserByScope(scope string, userID custom_type.ID) error {

	return nil
}
