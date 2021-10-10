package repositories

import "github.com/terdia/greenlight/src/users/entities"

type UserRepository interface {
	Insert(user *entities.User) error
	Update(user *entities.User) error
	GetByEmail(email string) (*entities.User, error)
	//GetAll(dto.ListMovieRequest) ([]*entities.User, data.Metadata, error)
}
