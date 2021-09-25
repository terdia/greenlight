package repository

import (
	"database/sql"
	"errors"

	"github.com/lib/pq"

	"github.com/terdia/greenlight/internal/data"
	"github.com/terdia/greenlight/src/movies/entities"
)

type movieRepository struct {
	sql.DB
}

func NewMovieRepoitory(db *sql.DB) *movieRepository {
	return &movieRepository{*db}
}

func (repo *movieRepository) Insert(movie *entities.Movie) error {
	query := `INSERT INTO movies (title, year, runtime, genres)
			 VALUES($1, $2, $3, $4)
			 RETURNING id, created_at, version`

	queryParams := []interface{}{movie.Title, movie.Year, movie.Runtime, pq.Array(movie.Genres)}

	return repo.DB.QueryRow(query, queryParams...).Scan(&movie.ID, &movie.CreatedAt, &movie.Version)
}

func (repo *movieRepository) Get(id int64) (*entities.Movie, error) {

	if id < 1 {
		return nil, data.ErrRecordNotFound
	}

	query := `SELECT id, created_at, title, year, runtime, genres, version
			  FROM movies
			  WHERE id = $1`

	var movie entities.Movie

	err := repo.DB.QueryRow(query, id).Scan(
		&movie.ID,
		&movie.CreatedAt,
		&movie.Title,
		&movie.Year,
		&movie.Runtime,
		pq.Array(&movie.Genres),
		&movie.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, data.ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &movie, nil
}

func (repo *movieRepository) Update(movie *entities.Movie) error {
	query := `
			UPDATE movies
			SET title = $1, year = $2, runtime = $3, genres = $4, version = version + 1
			WHERE id = $5
			RETURNING version`

	args := []interface{}{movie.Title, movie.Year, movie.Runtime, pq.Array(movie.Genres), movie.ID}

	return repo.DB.QueryRow(query, args...).Scan(&movie.Version)
}

func (repo *movieRepository) Delete(id int64) error {
	if id < 1 {
		return data.ErrRecordNotFound
	}

	query := `DELETE FROM movies WHERE id = $1`

	result, err := repo.DB.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return data.ErrRecordNotFound
	}

	return nil
}
