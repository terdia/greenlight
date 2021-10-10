package handlers

import (
	"fmt"
	"net/http"

	"github.com/terdia/greenlight/infrastructures/dto"
	"github.com/terdia/greenlight/internal/commons"
	"github.com/terdia/greenlight/internal/custom_type"
	"github.com/terdia/greenlight/src/users/entities"
	"github.com/terdia/greenlight/src/users/services"
)

type UserHandler interface {
	CreateMovie(rw http.ResponseWriter, r *http.Request)
}

type userHandler struct {
	sharedUtil commons.SharedUtil
	service    services.UserService
}

func NewUserHandler(utils commons.SharedUtil, srv services.UserService) UserHandler {
	return &userHandler{sharedUtil: utils, service: srv}
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
func (handler *userHandler) CreateMovie(rw http.ResponseWriter, r *http.Request) {

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

	result := commons.ResponseObject{
		StatusMsg: custom_type.Success,
		Data: dto.SingleUserResponse{
			User: getUserResponse(user),
		},
	}

	idString, _ := custom_type.EncodeId(int(user.ID))
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/users/%s", idString))

	err = handler.sharedUtil.WriteJson(rw, http.StatusCreated, result, headers)
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
		CreatedAt: user.CreatedAt,
		Version:   user.Version,
	}
}
