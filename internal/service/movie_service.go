package service

import (
	"github.com/terdia/greenlight/internal/data"
	"github.com/terdia/greenlight/internal/repository"
)

type MovieServiceInterface interface {
	Create(movie *data.Movie) error
	GetById(id int64) (*data.Movie, error)
	Update(movie *data.Movie) error
	Delete(id int64) error
}

type movieService struct {
	repo repository.MovieRepositoryInterface
}

func NewMovieService(repo repository.MovieRepositoryInterface) *movieService {
	return &movieService{repo: repo}
}

func (srv *movieService) Create(movie *data.Movie) error {
	return srv.repo.Insert(movie)
}

func (srv *movieService) GetById(id int64) (*data.Movie, error) {
	return srv.repo.Get(id)
}

func (srv *movieService) Update(movie *data.Movie) error {
	return srv.repo.Update(movie)
}

func (srv *movieService) Delete(id int64) error {
	return srv.repo.Delete(id)
}
