package main

import (
	"expvar"
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

	router.Use(app.metrics, app.recoverPanic, app.logRequest, app.enableCors, app.rateLimit, app.authenticate)

	router.Get("/v1/healthcheck", app.healthcheckHandler)

	//Domain routes
	movieHandler := app.registry.Handlers.MovieHandler
	userHandler := app.registry.Handlers.UserHandler

	router.Route("/v1/movies", func(r chi.Router) {

		r.Post("/", app.requirePermission("movies:write", movieHandler.CreateMovie))
		r.Get("/", app.requirePermission("movies:read", movieHandler.ListMovie))

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", app.requirePermission("movies:read", movieHandler.ShowMovie))
			r.Patch("/", app.requirePermission("movies:write", movieHandler.UpdateMovie))  // PATCH v1/movies/xxxx
			r.Delete("/", app.requirePermission("movies:write", movieHandler.DeleteMovie)) // DELETE v1/movies/xxxx
		})

	})

	router.Route("/v1/users", func(r chi.Router) {
		r.Post("/", userHandler.CreateUser)
		r.Put("/activated", userHandler.ActivateUser)
	})

	router.Post("/v1/tokens/authentication", userHandler.GetAuthenticationToken)

	//router.Get("/debug/vars", app.requirePermission("movies:read", expvar.Handler().ServeHTTP))
	router.Get("/debug/vars", expvar.Handler().ServeHTTP)

	// swagger API documentation UI
	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:4000/docs/swagger.json"),
	))

	fileServer := http.FileServer(http.Dir("./docs/"))
	router.Handle("/docs/*", http.StripPrefix("/docs", fileServer))

	return router
}
