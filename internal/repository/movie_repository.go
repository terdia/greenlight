package repository

import (
	"database/sql"

	"github.com/terdia/greenlight/internal/data"
)

type MovieRepositoryInterface interface {
	Insert(movie *data.Movie) error
	Get(id int64) (*data.Movie, error)
	Update(movie *data.Movie) error
	Delete(id int64) error
}

type movieRepository struct {
	sql.DB
}

func NewMovieRepoitory(db *sql.DB) *movieRepository {
	return &movieRepository{*db}
}

func (repo *movieRepository) Insert(movie *data.Movie) error {
	return nil
}

func (repo *movieRepository) Get(id int64) (*data.Movie, error) {
	return nil, nil
}

func (repo *movieRepository) Update(movie *data.Movie) error {
	return nil
}

func (repo *movieRepository) Delete(id int64) error {
	return nil
}
