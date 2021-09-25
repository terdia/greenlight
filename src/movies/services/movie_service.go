package services

import (
	"time"

	"github.com/terdia/greenlight/internal/validator"
	"github.com/terdia/greenlight/src/movies/entities"
	"github.com/terdia/greenlight/src/movies/repositories"
)

type CreateMovieValidationErrors map[string]string

type MovieService interface {
	Create(movie *entities.Movie) (CreateMovieValidationErrors, error)
	GetById(id int64) (*entities.Movie, error)
	Update(movie *entities.Movie) error
	Delete(id int64) error
}

type movieService struct {
	repo repositories.MovieRepository
}

func NewMovieService(repo repositories.MovieRepository) *movieService {
	return &movieService{repo: repo}
}

func (srv *movieService) Create(movie *entities.Movie) (CreateMovieValidationErrors, error) {
	v := validator.New()

	if validateCreateMovie(v, movie); !v.Valid() {
		return v.Errors, nil
	}

	return nil, srv.repo.Insert(movie)
}

func (srv *movieService) GetById(id int64) (*entities.Movie, error) {
	return srv.repo.Get(id)
}

func (srv *movieService) Update(movie *entities.Movie) error {
	return srv.repo.Update(movie)
}

func (srv *movieService) Delete(id int64) error {
	return srv.repo.Delete(id)
}

func validateCreateMovie(v *validator.Validator, movie *entities.Movie) {

	v.Check(movie.Title != "", "title", "must be provided")
	v.Check(len(movie.Title) <= 500, "title", "must not be more than 500 bytes long")
	v.Check(movie.Year != 0, "year", "must be provided")
	v.Check(movie.Year >= 1888, "year", "must be greater than 1888")
	v.Check(movie.Year <= int32(time.Now().Year()), "year", "must not be in the future")
	v.Check(movie.Runtime != 0, "runtime", "must be provided")
	v.Check(movie.Runtime > 0, "runtime", "must be a positive integer")
	v.Check(movie.Genres != nil, "genres", "must be provided")
	v.Check(len(movie.Genres) >= 1, "genres", "must contain at least 1 genre")
	v.Check(len(movie.Genres) <= 5, "genres", "must not contain more than 5 genres")
	v.Check(validator.UniqueStringSlice(movie.Genres), "genres", "must not contain duplicate values")
}
