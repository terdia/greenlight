package services

import (
	"errors"

	"github.com/terdia/greenlight/infrastructures/dto"
	"github.com/terdia/greenlight/internal/data"
	"github.com/terdia/greenlight/internal/validator"
	"github.com/terdia/greenlight/src/users/entities"
	"github.com/terdia/greenlight/src/users/repositories"
)

type UserValidationErrors map[string]string

type UserService interface {
	Create(request dto.CreateUserRequest) (*entities.User, UserValidationErrors, error)
}

type userService struct {
	repo            repositories.UserRepository
	passHashService PasswordHashService
}

func NewUserService(repo repositories.UserRepository, passHashService PasswordHashService) UserService {
	return &userService{repo: repo, passHashService: passHashService}
}

func (srv *userService) Create(request dto.CreateUserRequest) (*entities.User, UserValidationErrors, error) {

	password := entities.Password{PlainText: &request.Password}

	err := srv.passHashService.Hash(&password)
	if err != nil {
		return nil, nil, err
	}

	user := &entities.User{
		Name:      request.Name,
		Email:     request.Email,
		Password:  password,
		Activated: false,
	}

	v := validator.New()
	user.ValidateRequest(v)

	if !v.Valid() {
		return nil, v.Errors, nil
	}

	duplicateEmail, err := srv.repo.GetByEmail(request.Email)
	if err != nil {
		if !errors.Is(err, data.ErrRecordNotFound) {
			return nil, nil, err
		}
	}

	if duplicateEmail != nil {
		v.AddError("email", data.ErrDuplicateEmail.Error())
		return nil, v.Errors, nil
	}

	err = srv.repo.Insert(user)
	if err != nil {
		return nil, nil, err
	}

	return user, nil, nil
}
