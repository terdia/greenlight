package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {

	router := chi.NewRouter()

	utils := app.registry.Services.SharedUtil

	router.NotFound(utils.NotFoundResponse)
	router.MethodNotAllowed(utils.MethodNotAllowedResponse)

	router.Get("/v1/healthcheck", app.healthcheckHandler)

	//Domain routes
	movieHandler := app.registry.Handlers.MovieHandler

	router.Post("/v1/movies", movieHandler.CreateMovie)
	router.Get("/v1/movies/{id}", movieHandler.ShowMovie)
	router.Put("/v1/movies/{id}", movieHandler.UpdateMovie)
	router.Delete("/v1/movies/{id}", movieHandler.DeleteMovie)

	return router
}
