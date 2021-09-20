package main

import (
	"fmt"
	"net/http"

	"github.com/terdia/greenlight/internal/custom_type"
)

func (app *application) logError(r *http.Request, err error) {
	app.logger.Println(err)
}

func (app *application) serverErrorResponse(rw http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)

	message := "the server encountered a problem and could not process your request"
	app.errorResponse(rw, r, http.StatusInternalServerError, responseObject{
		Message: message,
	})
}

func (app *application) notFoundResponse(rw http.ResponseWriter, r *http.Request) {

	app.errorResponse(rw, r, http.StatusNotFound, responseObject{
		Message: "the requested resource could not be found",
	})
}

func (app *application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	app.errorResponse(w, r, http.StatusMethodNotAllowed, responseObject{
		Message: message,
	})
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.errorResponse(w, r, http.StatusBadRequest, responseObject{
		Message: err.Error(),
	})
}

func (app *application) failedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	app.errorResponse(w, r, http.StatusUnprocessableEntity, responseObject{
		StatusMsg: custom_type.Fail,
		Data: map[string]map[string]string{
			"errors": errors,
		},
	})
}
