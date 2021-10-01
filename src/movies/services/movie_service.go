package services

import (
	"time"

	"github.com/terdia/greenlight/infrastructures/dto"
	"github.com/terdia/greenlight/internal/validator"
	"github.com/terdia/greenlight/src/movies/entities"
	"github.com/terdia/greenlight/src/movies/repositories"
)

type CreateMovieValidationErrors map[string]string

type MovieService interface {
	Create(movie *entities.Movie) (CreateMovieValidationErrors, error)
	GetById(id int64) (*entities.Movie, error)
	Update(id int64, request dto.MovieRequest) (*entities.Movie, CreateMovieValidationErrors, error)
	Delete(id int64) error
	List(listMovieRequest dto.ListMovieRequest) ([]*entities.Movie, error)
}

type movieService struct {
	repo repositories.MovieRepository
}

func NewMovieService(repo repositories.MovieRepository) MovieService {
	return &movieService{repo: repo}
}

func (srv *movieService) Create(movie *entities.Movie) (CreateMovieValidationErrors, error) {
	v := validator.New()

	if validateMovie(v, movie); !v.Valid() {
		return v.Errors, nil
	}

	return nil, srv.repo.Insert(movie)
}

func (srv *movieService) GetById(id int64) (*entities.Movie, error) {
	return srv.repo.Get(id)
}

func (srv *movieService) Update(id int64, request dto.MovieRequest) (*entities.Movie, CreateMovieValidationErrors, error) {

	movie, err := srv.GetById(id)
	if err != nil {
		return nil, nil, err
	}

	if request.Title != nil {
		movie.Title = *request.Title
	}

	if request.Year != nil {
		movie.Year = *request.Year
	}

	if request.Runtime != nil {
		movie.Runtime = *request.Runtime
	}

	if request.Genres != nil {
		movie.Genres = request.Genres
	}

	v := validator.New()

	if validateMovie(v, movie); !v.Valid() {
		return nil, v.Errors, nil
	}

	return movie, nil, srv.repo.Update(movie)
}

func (srv *movieService) Delete(id int64) error {
	return srv.repo.Delete(id)
}

func (srv *movieService) List(listMovieRequest dto.ListMovieRequest) ([]*entities.Movie, error) {
	return srv.repo.GetAll(listMovieRequest)
}

func validateMovie(v *validator.Validator, movie *entities.Movie) {

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
