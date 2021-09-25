package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/terdia/greenlight/src/movies/handlers"
)

func (app *application) routes() http.Handler {

	router := chi.NewRouter()

	utils := app.registry.Services.SharedUtil

	router.NotFound(utils.NotFoundResponse)
	router.MethodNotAllowed(utils.MethodNotAllowedResponse)

	router.Get("/v1/healthcheck", app.healthcheckHandler)

	movieHandler := handlers.NewMovieHandler(utils, app.registry.Services.MovieService)

	router.Post("/v1/movies", movieHandler.CreateMovie)
	router.Get("/v1/movies/{id}", movieHandler.ShowMovie)

	return router
}
