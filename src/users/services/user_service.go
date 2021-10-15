package services

import (
	"errors"

	"github.com/terdia/greenlight/infrastructures/dto"
	"github.com/terdia/greenlight/internal/data"
	"github.com/terdia/greenlight/internal/mailer"
	"github.com/terdia/greenlight/internal/validator"
	"github.com/terdia/greenlight/src/users/entities"
	"github.com/terdia/greenlight/src/users/repositories"
)

type UserValidationErrors map[string]string

type UserService interface {
	Create(request dto.CreateUserRequest) (*entities.User, UserValidationErrors, error)
	SendMail(recipient, templateFile string, data interface{}) error
	ActivateUser(request dto.ActivateUserRequest) (*entities.User, UserValidationErrors, error)
}

type userService struct {
	repo            repositories.UserRepository
	passHashService PasswordHashService
	mailer          mailer.Mailer
}

func NewUserService(repo repositories.UserRepository, passHashService PasswordHashService, mailer mailer.Mailer) UserService {
	return &userService{repo: repo, passHashService: passHashService, mailer: mailer}
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

func (srv *userService) SendMail(recipient, templateFile string, data interface{}) error {

	return srv.mailer.Send(recipient, templateFile, data)
}

func (srv *userService) ActivateUser(request dto.ActivateUserRequest) (*entities.User, UserValidationErrors, error) {

	v := validator.New()

	if validateActivateUserRequest(v, request); !v.Valid() {
		return nil, v.Errors, nil
	}

	user, err := srv.repo.GetForToken(request.TokenPlaintext, data.TokenScopeActivation)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			v.AddError("token", "invalid or expired token")
			return nil, v.Errors, nil
		default:
			return nil, nil, err
		}
	}

	user.Activated = true

	err = srv.repo.Update(user)
	if err != nil {
		return nil, nil, err
	}

	return user, nil, nil
}

func validateActivateUserRequest(v *validator.Validator, r dto.ActivateUserRequest) {
	v.Check(r.TokenPlaintext != "", "token", "must be provided")
	v.Check(len(r.TokenPlaintext) == 26, "token", "must be 26 bytes long")
}
