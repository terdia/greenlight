package handlers

import (
	"errors"
	"net/http"

	"github.com/terdia/greenlight/infrastructures/dto"
	"github.com/terdia/greenlight/internal/commons"
	"github.com/terdia/greenlight/internal/custom_type"
	"github.com/terdia/greenlight/internal/data"
)

// GetAuthenticationToken ... Get authentication token
// @Summary Get user authentication token
// @Description Generate a new token for a user using the given credentials
// @Tags Token
// @Param body body dto.AuthTokenRequest true "auth token credentials"
// @Success 200 {object} commons.ResponseObject{data=dto.TokenResponse}
// @Failure 422 {object} commons.ResponseObject{data=dto.ValidationError} "status: fail"
// @Failure 401,500 {object} commons.ResponseObject "e.g. status: error, message: the error reason"
// @Router /tokens/authentication [post]
func (handler *userHandler) GetAuthenticationToken(rw http.ResponseWriter, r *http.Request) {

	request := dto.AuthTokenRequest{}

	utils := handler.sharedUtil

	err := utils.ReadJson(rw, r, &request)
	if err != nil {
		utils.BadRequestResponse(rw, r, err)

		return
	}

	token, validationErrors, err := handler.service.CreateAuthenticationToken(request, data.TokenScopeAuthentication)
	if validationErrors != nil {
		utils.FailedValidationResponse(rw, r, validationErrors)

		return
	}

	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			utils.InvalidCredentialsResponse(rw, r)
		case errors.Is(err, data.ErrInvalidCredentials):
			utils.InvalidCredentialsResponse(rw, r)
		default:
			utils.ServerErrorResponse(rw, r, err)
		}
		return
	}

	tokenDto := dto.Token{
		PlainText: token.Plaintext,
		Expiry:    token.Expiry,
	}

	err = handler.sharedUtil.WriteJson(rw, http.StatusOK, commons.ResponseObject{
		StatusMsg: custom_type.Success,
		Data: dto.TokenResponse{
			Token: tokenDto,
		},
	}, nil)

	if err != nil {
		handler.sharedUtil.ServerErrorResponse(rw, r, err)

		return
	}

}
