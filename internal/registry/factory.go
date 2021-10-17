package registry

import (
	"database/sql"
	"sync"

	"github.com/terdia/greenlight/infrastructures/logger"
	"github.com/terdia/greenlight/infrastructures/persistence/postgres/repository"
	"github.com/terdia/greenlight/internal/commons"
	"github.com/terdia/greenlight/internal/mailer"
	"github.com/terdia/greenlight/src/movies/handlers"
	"github.com/terdia/greenlight/src/movies/services"
	user_handler "github.com/terdia/greenlight/src/users/handlers"
	user_repository "github.com/terdia/greenlight/src/users/repositories"
	user_services "github.com/terdia/greenlight/src/users/services"
)

type Registry struct {
	Services *Services
	Handlers *Handlers
}

type Services struct {
	SharedUtil     commons.SharedUtil
	UserService    user_services.UserService
	UserRepository user_repository.UserRepository
}

type Handlers struct {
	MovieHandler handlers.MovieHandle
	UserHandler  user_handler.UserHandler
}

//todo clean up, split into domains and aggregate here
func NewRegistry(db *sql.DB, logger *logger.Logger, mailer mailer.Mailer, wg *sync.WaitGroup) Registry {

	userRepository := repository.NewUserRepoitory(db)

	utils := commons.NewUtil(logger, wg)
	movieService := services.NewMovieService(repository.NewMovieRepoitory(db))

	tokenService := user_services.NewTokenService(
		repository.NewTokenRepository(db),
	)

	userService := user_services.NewUserService(
		userRepository,
		user_services.NewPasswordService(),
		mailer,
		tokenService,
	)

	services := newServices(utils, userService, userRepository)

	movieHandler := handlers.NewMovieHandler(utils, movieService)
	userHandler := user_handler.NewUserHandler(utils, userService, tokenService)

	handlers := newHandlers(movieHandler, userHandler)

	return Registry{
		Services: services,
		Handlers: handlers,
	}
}

func newServices(
	sharedUtil commons.SharedUtil,
	userService user_services.UserService,
	userRepository user_repository.UserRepository,
) *Services {
	return &Services{
		SharedUtil:     sharedUtil,
		UserService:    userService,
		UserRepository: userRepository,
	}
}

func newHandlers(movieHandler handlers.MovieHandle, userHandler user_handler.UserHandler) *Handlers {
	return &Handlers{
		MovieHandler: movieHandler,
		UserHandler:  userHandler,
	}
}
