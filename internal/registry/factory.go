package registry

import (
	"database/sql"

	"github.com/terdia/greenlight/internal/commons"
	"github.com/terdia/greenlight/src/movies/services"
)

type Registry struct {
	Services *Services
}

type Services struct {
	MovieService services.MovieService
	SharedUtil   commons.SharedUtil
}

func NewRegistry(db *sql.DB, services *Services) Registry {
	return Registry{
		Services: services,
	}
}

func NewServices(movieService services.MovieService, sharedUtil commons.SharedUtil) *Services {
	return &Services{
		MovieService: movieService,
		SharedUtil:   sharedUtil,
	}
}
