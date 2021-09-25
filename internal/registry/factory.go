package registry

import (
	"database/sql"
	"log"

	"github.com/terdia/greenlight/infrastructures/persistence/postgres/repository"
	"github.com/terdia/greenlight/internal/commons"
	"github.com/terdia/greenlight/src/movies/handlers"
	"github.com/terdia/greenlight/src/movies/services"
)

type Registry struct {
	Services *Services
	Handlers *Handlers
}

type Services struct {
	SharedUtil commons.SharedUtil
}

type Handlers struct {
	MovieHandler handlers.MovieHandle
}

func NewRegistry(db *sql.DB, logger *log.Logger) Registry {

	utils := commons.NewUtil(logger)
	movieService := services.NewMovieService(repository.NewMovieRepoitory(db))

	services := newServices(utils)

	movieHandler := handlers.NewMovieHandler(utils, movieService)

	handlers := newHandlers(movieHandler)

	return Registry{
		Services: services,
		Handlers: handlers,
	}
}

func newServices(sharedUtil commons.SharedUtil) *Services {
	return &Services{
		SharedUtil: sharedUtil,
	}
}

func newHandlers(movieHandler handlers.MovieHandle) *Handlers {
	return &Handlers{
		MovieHandler: movieHandler,
	}
}
