package mock

import (
	"sync"

	"github.com/terdia/greenlight/infrastructures/logger"
	"github.com/terdia/greenlight/internal/commons"
	"github.com/terdia/greenlight/internal/mailer"
	"github.com/terdia/greenlight/internal/registry"
	"github.com/terdia/greenlight/src/movies/handlers"
	"github.com/terdia/greenlight/src/movies/services"
	user_handler "github.com/terdia/greenlight/src/users/handlers"
	user_repository "github.com/terdia/greenlight/src/users/repositories"
	user_services "github.com/terdia/greenlight/src/users/services"
)

type Handlers struct {
	MovieHandler handlers.MovieHandle
	UserHandler  user_handler.UserHandler
}

//todo clean up, split into domains and aggregate here
func NewRegistry(logger *logger.Logger, mailer mailer.Mailer, wg *sync.WaitGroup, movieCount int) registry.Registry {

	userRepository := NewUserRepoitoryMock()
	permissionRepository := NewPermissionRepositoryMock()

	utils := commons.NewUtil(logger, wg)
	movieService := services.NewMovieService(NewMovieRepoitoryMock(movieCount))

	tokenService := user_services.NewTokenService(NewTokenRepositoryMock())

	userService := user_services.NewUserService(
		userRepository,
		user_services.NewPasswordService(),
		mailer,
		tokenService,
	)

	services := newServices(utils, userService, userRepository, permissionRepository)

	movieHandler := handlers.NewMovieHandler(utils, movieService)
	userHandler := user_handler.NewUserHandler(utils, userService, tokenService, permissionRepository)

	handlers := newHandlers(movieHandler, userHandler)

	return registry.Registry{
		Services: services,
		Handlers: handlers,
	}
}

func newServices(
	sharedUtil commons.SharedUtil,
	userService user_services.UserService,
	userRepository user_repository.UserRepository,
	permissionRepository user_repository.PermissionRepository,
) *registry.Services {
	return &registry.Services{
		SharedUtil:           sharedUtil,
		UserService:          userService,
		UserRepository:       userRepository,
		PermissionRepository: permissionRepository,
	}
}

func newHandlers(
	movieHandler handlers.MovieHandle,
	userHandler user_handler.UserHandler,
) *registry.Handlers {
	return &registry.Handlers{
		MovieHandler: movieHandler,
		UserHandler:  userHandler,
	}
}
