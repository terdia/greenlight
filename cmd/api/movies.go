package main

import (
	"net/http"
	"time"

	"github.com/terdia/greenlight/internal/custom_type"
	"github.com/terdia/greenlight/internal/data"
	"github.com/terdia/greenlight/internal/validator"
)

func (app *application) createMovieHandler(rw http.ResponseWriter, r *http.Request) {
	var input struct {
		Title   string              `json:"title"`
		Year    int32               `json:"year"`
		Runtime custom_type.Runtime `json:"runtime"`
		Genres  []string            `json:"genres"`
	}

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

	v := validator.New()

	if movie.Validate(v); !v.Valid() {
		app.failedValidationResponse(rw, r, v.Errors)

		return
	}

	result := responseObject{
		StatusMsg: custom_type.Success,
		Data: map[string]data.Movie{
			"movie": *movie,
		},
	}

	err = app.writeJson(rw, http.StatusOK, result, nil)
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
