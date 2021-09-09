package main

import (
	"net/http"
	"time"

	"github.com/terdia/greenlight/internal/data"
)

func (app *application) createMovieHandler(rw http.ResponseWriter, r *http.Request) {
	err := app.writeJson(rw, http.StatusOK, responseData{"status": "success", "message": "create a new movie"}, nil)
	if err != nil {
		app.logger.Println(err)
		http.Error(rw, "Server cannot process your request", http.StatusInternalServerError)

		return
	}
}

func (app *application) showMovieHandler(rw http.ResponseWriter, r *http.Request) {

	id, err := app.extractIdParamFromContext(r)
	if err != nil {
		http.NotFound(rw, r)
		return
	}

	result := responseData{
		"status": "success",
		"data": map[string]data.Movie{
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
		app.logger.Println(err)
		http.Error(rw, "Server cannot process your request", http.StatusInternalServerError)

		return
	}
}
