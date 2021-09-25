package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/terdia/greenlight/infrastructures/dto"
	"github.com/terdia/greenlight/internal/commons"
	"github.com/terdia/greenlight/internal/custom_type"
	"github.com/terdia/greenlight/src/movies/entities"
	"github.com/terdia/greenlight/src/movies/services"
)

type MovieHandle interface {
	CreateMovie(rw http.ResponseWriter, r *http.Request)
	ShowMovie(rw http.ResponseWriter, r *http.Request)
}

type movieHandler struct {
	sharedUtil commons.SharedUtil
	service    services.MovieService
}

func NewMovieHandler(util commons.SharedUtil, service services.MovieService) *movieHandler {
	return &movieHandler{
		sharedUtil: util,
		service:    service,
	}
}

func (handler *movieHandler) CreateMovie(rw http.ResponseWriter, r *http.Request) {
	var input dto.CreateMovieRequest

	err := handler.sharedUtil.ReadJson(rw, r, &input)
	if err != nil {
		handler.sharedUtil.BadRequestResponse(rw, r, err)

		return
	}

	movie := &entities.Movie{
		Title:   input.Title,
		Year:    input.Year,
		Runtime: input.Runtime,
		Genres:  input.Genres,
	}

	validationErrors, err := handler.service.Create(movie)
	if validationErrors != nil {
		handler.sharedUtil.FailedValidationResponse(rw, r, validationErrors)

		return
	}
	if err != nil {
		handler.sharedUtil.ServerErrorResponse(rw, r, err)

		return
	}

	result := commons.ResponseObject{
		StatusMsg: custom_type.Success,
		Data: map[string]entities.Movie{
			"movie": *movie,
		},
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/movies/%d", movie.ID))

	err = handler.sharedUtil.WriteJson(rw, http.StatusCreated, result, headers)
	if err != nil {
		handler.sharedUtil.ServerErrorResponse(rw, r, err)

		return
	}
}

func (handler *movieHandler) ShowMovie(rw http.ResponseWriter, r *http.Request) {

	id, err := handler.sharedUtil.ExtractIdParamFromContext(r)
	if err != nil {
		handler.sharedUtil.NotFoundResponse(rw, r)

		return
	}

	result := commons.ResponseObject{
		StatusMsg: custom_type.Success,
		Data: map[string]entities.Movie{
			"movie": {
				ID:        custom_type.ID(id),
				Title:     "Casablanca",
				Runtime:   102,
				Genres:    []string{"drama", "romance", "war"},
				Version:   1,
				CreatedAt: time.Now(),
			},
		},
	}

	err = handler.sharedUtil.WriteJson(rw, http.StatusOK, result, nil)
	if err != nil {
		handler.sharedUtil.ServerErrorResponse(rw, r, err)

		return
	}
}
