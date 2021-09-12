package main

import (
	"net/http"
	"time"

	"github.com/terdia/greenlight/internal/custom_type"
	"github.com/terdia/greenlight/internal/data"
)

func (app *application) createMovieHandler(rw http.ResponseWriter, r *http.Request) {
	err := app.writeJson(rw, http.StatusOK, responseObject{StatusMsg: custom_type.Success, Message: "create a new movie"}, nil)
	if err != nil {
		app.serverErrorResponse(rw, r, err)
	}
}

func (app *application) showMovieHandler(rw http.ResponseWriter, r *http.Request) {

	id, err := app.extractIdParamFromContext(r)
	if err != nil {
		app.notFoundResponse(rw, r)

		return
	}

	result := responseObject{
		StatusMsg: custom_type.Success,
		Data: map[string]data.Movie{
			"movie": {
				ID:        id,
				Title:     "Casablanca",
				Runtime:   102,
				Genres:    []string{"drama", "romance", "war"},
				Version:   1,
				CreatedAt: time.Now(),
			},
		},
	}

	err = app.writeJson(rw, http.StatusOK, result, nil)
	if err != nil {
		app.serverErrorResponse(rw, r, err)

		return
	}
}
