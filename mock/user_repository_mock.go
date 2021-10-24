package mock

import (
	"github.com/terdia/greenlight/src/users/entities"
	"github.com/terdia/greenlight/src/users/repositories"
)

type userRepositoryMock struct{}

func NewUserRepoitoryMock() repositories.UserRepository {
	return &userRepositoryMock{}
}

func (repo *userRepositoryMock) Insert(user *entities.User) error {

	return nil
}

func (repo *userRepositoryMock) GetByEmail(email string) (*entities.User, error) {

	var user entities.User

	return &user, nil
}

func (repo *userRepositoryMock) Update(user *entities.User) error {

	return nil
}

func (repo *userRepositoryMock) GetForToken(tokenPlainText, scope string) (*entities.User, error) {

	var user entities.User

	return &user, nil

}
