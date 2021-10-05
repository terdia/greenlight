package commons

import (
	"fmt"
	"net/http"

	"github.com/terdia/greenlight/infrastructures/dto"
	"github.com/terdia/greenlight/internal/custom_type"
)

func (util *sharedUtils) LogError(r *http.Request, err error) {
	util.logger.PrintError(err, map[string]string{
		"request_method": r.Method,
		"request_url":    r.URL.String(),
	})
}

func (util *sharedUtils) ServerErrorResponse(rw http.ResponseWriter, r *http.Request, err error) {
	util.LogError(r, err)

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
