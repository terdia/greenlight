package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/terdia/greenlight/docs"
)

func (app *application) routes() http.Handler {

	router := chi.NewRouter()

	utils := app.registry.Services.SharedUtil

	router.NotFound(utils.NotFoundResponse)
	router.MethodNotAllowed(utils.MethodNotAllowedResponse)

	router.Use(app.recoverPanic, app.rateLimit)

	router.Get("/v1/healthcheck", app.healthcheckHandler)

	//Domain routes
	movieHandler := app.registry.Handlers.MovieHandler

	router.Get("/v1/movies", movieHandler.ListMovie)
	router.Get("/v1/movies/{id}", movieHandler.ShowMovie)
	router.Post("/v1/movies", movieHandler.CreateMovie)
	router.Patch("/v1/movies/{id}", movieHandler.UpdateMovie)
	router.Delete("/v1/movies/{id}", movieHandler.DeleteMovie)

	// swagger API documentation UI
	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:4000/docs/swagger.json"),
	))

	fileServer := http.FileServer(http.Dir("./docs/"))
	router.Handle("/docs/*", http.StripPrefix("/docs", fileServer))

	return router
}
