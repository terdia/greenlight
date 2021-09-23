package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/terdia/greenlight/internal/custom_type"
	"github.com/terdia/greenlight/internal/data"
)

func (app *application) createMovieHandler(rw http.ResponseWriter, r *http.Request) {
	var input data.CreateMovieRequest

	err := app.readJson(rw, r, &input)
	if err != nil {
		app.badRequestResponse(rw, r, err)

		return
	}

	movie := &data.Movie{
		Title:   input.Title,
		Year:    input.Year,
		Runtime: input.Runtime,
		Genres:  input.Genres,
	}

	validationErrors, err := app.services.MovieService.Create(movie)
	if validationErrors != nil {
		app.failedValidationResponse(rw, r, validationErrors)

		return
	}
	if err != nil {
		app.serverErrorResponse(rw, r, err)

		return
	}

	result := responseObject{
		StatusMsg: custom_type.Success,
		Data: map[string]data.Movie{
			"movie": *movie,
		},
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/movies/%d", movie.ID))

	err = app.writeJson(rw, http.StatusCreated, result, headers)
	if err != nil {
		app.serverErrorResponse(rw, r, err)

		return
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
