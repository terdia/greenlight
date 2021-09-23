package repository

import (
	"database/sql"

	"github.com/lib/pq"

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
	query := `INSERT INTO movies (title, year, runtime, genres)
			 VALUES($1, $2, $3, $4)
			 RETURNING id, created_at, version`

	queryParams := []interface{}{movie.Title, movie.Year, movie.Runtime, pq.Array(movie.Genres)}

	return repo.DB.QueryRow(query, queryParams...).Scan(&movie.ID, &movie.CreatedAt, &movie.Version)
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
