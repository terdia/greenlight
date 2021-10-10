package registry

import (
	"database/sql"

	"github.com/terdia/greenlight/infrastructures/logger"
	"github.com/terdia/greenlight/infrastructures/persistence/postgres/repository"
	"github.com/terdia/greenlight/internal/commons"
	"github.com/terdia/greenlight/src/movies/handlers"
	"github.com/terdia/greenlight/src/movies/services"
	user_handler "github.com/terdia/greenlight/src/users/handlers"
	user_services "github.com/terdia/greenlight/src/users/services"
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
	UserHandler  user_handler.UserHandler
}

func NewRegistry(db *sql.DB, logger *logger.Logger) Registry {

	utils := commons.NewUtil(logger)
	movieService := services.NewMovieService(repository.NewMovieRepoitory(db))
	userService := user_services.NewUserService(repository.NewUserRepoitory(db), user_services.NewPasswordService())

	services := newServices(utils)

	movieHandler := handlers.NewMovieHandler(utils, movieService)
	userHandler := user_handler.NewUserHandler(utils, userService)

	handlers := newHandlers(movieHandler, userHandler)

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

func newHandlers(movieHandler handlers.MovieHandle, userHandler user_handler.UserHandler) *Handlers {
	return &Handlers{
		MovieHandler: movieHandler,
		UserHandler:  userHandler,
	}
}
