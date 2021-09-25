package repositories

import (
	"github.com/terdia/greenlight/src/movies/entities"
)

type MovieRepository interface {
	Insert(movie *entities.Movie) error
	Get(id int64) (*entities.Movie, error)
	Update(movie *entities.Movie) error
	Delete(id int64) error
}