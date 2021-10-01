package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"

	"github.com/terdia/greenlight/infrastructures/dto"
	"github.com/terdia/greenlight/internal/data"
	"github.com/terdia/greenlight/src/movies/entities"
	"github.com/terdia/greenlight/src/movies/repositories"
)

type movieRepository struct {
	sql.DB
}

func NewMovieRepoitory(db *sql.DB) repositories.MovieRepository {
	return &movieRepository{*db}
}

func (repo *movieRepository) Insert(movie *entities.Movie) error {
	query := `INSERT INTO movies (title, year, runtime, genres)
			 VALUES($1, $2, $3, $4)
			 RETURNING id, created_at, version`

	queryParams := []interface{}{movie.Title, movie.Year, movie.Runtime, pq.Array(movie.Genres)}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	return repo.DB.QueryRowContext(ctx, query, queryParams...).Scan(&movie.ID, &movie.CreatedAt, &movie.Version)
}

func (repo *movieRepository) Get(id int64) (*entities.Movie, error) {

	if id < 1 {
		return nil, data.ErrRecordNotFound
	}

	query := `SELECT id, created_at, title, year, runtime, genres, version
			  FROM movies
			  WHERE id = $1`

	var movie entities.Movie

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	err := repo.DB.QueryRowContext(ctx, query, id).Scan(
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
	// To enable Optimistic Concurrency Control (data race condition during edit)
	//add version to where clause, to ensure first routine to send update request is persisted and  other routine
	// receive edit error
	query := `
			UPDATE movies
			SET title = $1, year = $2, runtime = $3, genres = $4, version = version + 1
			WHERE id = $5 AND version = $6
			RETURNING version`

	args := []interface{}{movie.Title, movie.Year, movie.Runtime, pq.Array(movie.Genres), movie.ID, movie.Version}

	// Execute the SQL query. If no matching row could be found, we know the movie
	// version has changed (or the record has been deleted) and we return our custom
	// ErrEditConflict error.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	err := repo.DB.QueryRowContext(ctx, query, args...).Scan(&movie.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return data.ErrEditConflict
		default:
			return err
		}
	}

	return nil
}

func (repo *movieRepository) Delete(id int64) error {
	if id < 1 {
		return data.ErrRecordNotFound
	}

	query := `DELETE FROM movies WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	result, err := repo.DB.ExecContext(ctx, query, id)
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

func (repo *movieRepository) GetAll(dto.ListMovieRequest) ([]*entities.Movie, error) {

	query := `
			SELECT id, created_at, title, year, runtime, genres, version
			FROM movies
			ORDER BY id`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := repo.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	movies := []*entities.Movie{}

	for rows.Next() {
		var movie entities.Movie

		err := rows.Scan(
			&movie.ID,
			&movie.CreatedAt,
			&movie.Title,
			&movie.Year,
			&movie.Runtime,
			pq.Array(&movie.Genres),
			&movie.Version,
		)

		if err != nil {
			return nil, err
		}

		movies = append(movies, &movie)
	}

	return movies, nil
}
