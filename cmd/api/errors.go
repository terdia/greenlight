package main

import (
	"fmt"
	"net/http"
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
