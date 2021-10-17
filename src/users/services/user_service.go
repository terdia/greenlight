package services

import (
	"errors"
	"time"

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
	CreateAuthenticationToken(request dto.AuthTokenRequest, scope string) (*entities.Token, UserValidationErrors, error)
}

type userService struct {
	repo            repositories.UserRepository
	passHashService PasswordHashService
	mailer          mailer.Mailer
	tokenService    TokenService
}

func NewUserService(
	repo repositories.UserRepository,
	passHashService PasswordHashService,
	mailer mailer.Mailer,
	tokenService TokenService,
) UserService {
	return &userService{
		repo:            repo,
		passHashService: passHashService,
		mailer:          mailer,
		tokenService:    tokenService,
	}
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

	if validateTokenRequest(v, request.TokenPlaintext); !v.Valid() {
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

func (srv *userService) CreateAuthenticationToken(
	request dto.AuthTokenRequest,
	scope string,
) (*entities.Token, UserValidationErrors, error) {

	v := validator.New()

	entities.ValidateEmail(v, request.Email)
	entities.ValidatePasswordPlaintext(v, request.Password)

	if !v.Valid() {
		return nil, v.Errors, nil
	}

	user, err := srv.repo.GetByEmail(request.Email)
	if err != nil {
		return nil, nil, err
	}

	matchPassword, err := srv.passHashService.Verify(user.Password.Hash, request.Password)
	if err != nil {
		return nil, nil, err
	}

	if !matchPassword {
		return nil, nil, data.ErrInvalidCredentials
	}

	token, err := srv.tokenService.CreateNew(user.ID, 24*time.Hour, data.TokenScopeAuthentication)

	return token, nil, err
}

func validateTokenRequest(v *validator.Validator, plainText string) {
	v.Check(plainText != "", "token", "must be provided")
	v.Check(len(plainText) == 26, "token", "must be 26 bytes long")
}
