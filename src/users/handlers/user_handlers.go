package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/terdia/greenlight/infrastructures/dto"
	"github.com/terdia/greenlight/internal/commons"
	"github.com/terdia/greenlight/internal/custom_type"
	"github.com/terdia/greenlight/internal/data"
	"github.com/terdia/greenlight/src/users/entities"
	"github.com/terdia/greenlight/src/users/services"
)

type UserHandler interface {
	CreateUser(rw http.ResponseWriter, r *http.Request)
	ActivateUser(rw http.ResponseWriter, r *http.Request)
}

type userHandler struct {
	sharedUtil   commons.SharedUtil
	service      services.UserService
	tokenService services.TokenService
}

func NewUserHandler(
	utils commons.SharedUtil,
	srv services.UserService,
	tokenService services.TokenService,
) UserHandler {
	return &userHandler{sharedUtil: utils, service: srv, tokenService: tokenService}
}

// CreateUser ... Create user
// @Summary Create new user
// @Description create a new user with given details
// @Tags Users
// @Param body body dto.CreateUserRequest true "create user"
// @Success 200 {object} commons.ResponseObject{data=dto.SingleUserResponse}
// @Header 200 {string} Location "/v1/users/QbPy4B7a2Lw1Kg7ogoEWj9k3NGMRVY"
// @Failure 422 {object} commons.ResponseObject{data=dto.ValidationError} "status: fail"
// @Failure 400,500 {object} commons.ResponseObject "e.g. status: error, message: the error reason"
// @Router /users [post]
func (handler *userHandler) CreateUser(rw http.ResponseWriter, r *http.Request) {

	request := dto.CreateUserRequest{}
	utils := handler.sharedUtil

	err := utils.ReadJson(rw, r, &request)
	if err != nil {
		utils.BadRequestResponse(rw, r, err)

		return
	}

	user, validationErrors, err := handler.service.Create(request)
	if validationErrors != nil {
		utils.FailedValidationResponse(rw, r, validationErrors)

		return
	}

	if err != nil {
		utils.ServerErrorResponse(rw, r, err)

		return
	}

	idString, _ := custom_type.EncodeId(int(user.ID))
	token, err := handler.tokenService.CreateNew(user.ID, 3*24*time.Hour, data.TokenScopeActivation)
	if err != nil {
		utils.ServerErrorResponse(rw, r, err)

		return
	}

	// send welocme email using background process
	utils.Background(func() {

		templateData := struct {
			ID    string
			Token string
		}{
			ID:    idString,
			Token: token.Plaintext,
		}

		err = handler.service.SendMail(user.Email, "user_welcome.tmpl", templateData)
		if err != nil {
			utils.LogErrorWithContext(err, map[string]string{
				"task":            "email sending gorountine",
				"userId":          idString,
				"activationToken": token.Plaintext,
			})
		}
	})

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/users/%s", idString))

	err = handler.sharedUtil.WriteJson(rw, http.StatusCreated, commons.ResponseObject{
		StatusMsg: custom_type.Success,
		Data: dto.SingleUserResponse{
			User: getUserResponse(user),
		},
	}, headers)
	if err != nil {
		handler.sharedUtil.ServerErrorResponse(rw, r, err)

		return
	}

}

// ActivateUser ... Activate user account
// @Summary Activate user account
// @Description activate the account of a user using the given token
// @Tags Users
// @Param body body dto.ActivateUserRequest true "activate user"
// @Success 200 {object} commons.ResponseObject{data=dto.SingleUserResponse}
// @Failure 409 {object} commons.ResponseObject "e.g. status: error, message: unable to update the record due to an edit conflict, please try again"
// @Failure 422 {object} commons.ResponseObject{data=dto.ValidationError} "status: fail"
// @Failure 400,500 {object} commons.ResponseObject "e.g. status: error, message: the error reason"
// @Router /users/activated [put]
func (handler *userHandler) ActivateUser(rw http.ResponseWriter, r *http.Request) {
	request := dto.ActivateUserRequest{}
	utils := handler.sharedUtil

	err := utils.ReadJson(rw, r, &request)
	if err != nil {
		utils.BadRequestResponse(rw, r, err)

		return
	}

	user, validationErrors, err := handler.service.ActivateUser(request)
	if validationErrors != nil {
		utils.FailedValidationResponse(rw, r, validationErrors)

		return
	}

	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			utils.EditConflictResponse(rw, r)
		default:
			utils.ServerErrorResponse(rw, r, err)
		}
		return
	}

	// run background task to delete tokens
	utils.Background(func() {
		plainTokenText := request.TokenPlaintext
		err := handler.tokenService.DeleteByUserIdAndScope(user.ID, data.TokenScopeActivation)

		if err != nil {
			encodedUserId, err := custom_type.EncodeId(int(user.ID))
			utils.LogErrorWithContext(err, map[string]string{
				"task":            "user activation, delete token for user by scope",
				"userId":          encodedUserId,
				"activationToken": plainTokenText,
			})
		}

	})

	err = handler.sharedUtil.WriteJson(rw, http.StatusOK, commons.ResponseObject{
		StatusMsg: custom_type.Success,
		Data: dto.SingleUserResponse{
			User: getUserResponse(user),
		},
	}, nil)
	if err != nil {
		handler.sharedUtil.ServerErrorResponse(rw, r, err)

		return
	}

}

func getUserResponse(user *entities.User) dto.UserResponse {
	return dto.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Activated: user.Activated,
		CreatedAt: user.CreatedAt,
		Version:   user.Version,
	}
}
