package service

import (
	"database/sql"

	"github.com/terdia/greenlight/internal/repository"
)

type Services struct {
	MovieService MovieServiceInterface
}

func NewServices(db *sql.DB) Services {
	return Services{
		MovieService: NewMovieService(repository.NewMovieRepoitory(db)),
	}
}
