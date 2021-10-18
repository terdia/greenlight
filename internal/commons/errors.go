package commons

import (
	"fmt"
	"net/http"

	"github.com/terdia/greenlight/infrastructures/dto"
	"github.com/terdia/greenlight/internal/custom_type"
)

func (util *sharedUtils) LogErrorWithHttpRequestContext(r *http.Request, err error) {
	util.logger.PrintError(err, map[string]string{
		"request_method": r.Method,
		"request_url":    r.URL.String(),
	})
}

func (util *sharedUtils) LogErrorWithContext(err error, context map[string]string) {
	util.logger.PrintError(err, context)
}

func (util *sharedUtils) ServerErrorResponse(rw http.ResponseWriter, r *http.Request, err error) {
	util.LogErrorWithHttpRequestContext(r, err)

	message := "the server encountered a problem and could not process your request"
	util.ErrorResponse(rw, r, http.StatusInternalServerError, ResponseObject{
		Message: message,
	})
}

func (util *sharedUtils) NotFoundResponse(rw http.ResponseWriter, r *http.Request) {

	util.ErrorResponse(rw, r, http.StatusNotFound, ResponseObject{
		Message: "the requested resource could not be found",
	})
}

func (util *sharedUtils) MethodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	util.ErrorResponse(w, r, http.StatusMethodNotAllowed, ResponseObject{
		Message: message,
	})
}

func (util *sharedUtils) BadRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	util.ErrorResponse(w, r, http.StatusBadRequest, ResponseObject{
		Message: err.Error(),
	})
}

func (util *sharedUtils) FailedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	util.ErrorResponse(w, r, http.StatusUnprocessableEntity, ResponseObject{
		StatusMsg: custom_type.Fail,
		Data: dto.ValidationError{
			Errors: errors,
		},
	})
}

func (util *sharedUtils) EditConflictResponse(w http.ResponseWriter, r *http.Request) {
	util.ErrorResponse(w, r, http.StatusConflict, ResponseObject{
		Message: "unable to update the record due to an edit conflict, please try again",
	})
}

func (util *sharedUtils) RateLimitExceededResponse(w http.ResponseWriter, r *http.Request) {

	util.ErrorResponse(w, r, http.StatusTooManyRequests, ResponseObject{
		Message: "rate limit exceeded",
	})
}

func (util *sharedUtils) InvalidCredentialsResponse(w http.ResponseWriter, r *http.Request) {

	util.ErrorResponse(w, r, http.StatusUnauthorized, ResponseObject{
		StatusMsg: custom_type.Fail,
		Message:   "invalid authentication credentials",
	})
}

func (util *sharedUtils) InvalidAuthenticationTokenResponse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("WWW-Authenticate", "Bearer")

	util.ErrorResponse(w, r, http.StatusUnauthorized, ResponseObject{
		StatusMsg: custom_type.Fail,
		Message:   "invalid or missing token",
	})
}

func (util *sharedUtils) AuthenticationRequiredResponse(w http.ResponseWriter, r *http.Request) {
	util.ErrorResponse(w, r, http.StatusUnauthorized, ResponseObject{
		Message: "you must be authenticated to access this resource",
	})
}

func (util *sharedUtils) InactiveAccountResponse(w http.ResponseWriter, r *http.Request) {
	util.ErrorResponse(w, r, http.StatusForbidden, ResponseObject{
		Message: "your user account must be activated to access this resource",
	})
}

func (util *sharedUtils) NotPermittedRResponse(w http.ResponseWriter, r *http.Request) {
	util.ErrorResponse(w, r, http.StatusForbidden, ResponseObject{
		Message: "your user account doesn't have the necessary permissions to perform this operation",
	})
}
