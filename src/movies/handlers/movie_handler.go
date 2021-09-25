package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/terdia/greenlight/infrastructures/dto"
	"github.com/terdia/greenlight/internal/commons"
	"github.com/terdia/greenlight/internal/custom_type"
	"github.com/terdia/greenlight/internal/data"
	"github.com/terdia/greenlight/src/movies/entities"
	"github.com/terdia/greenlight/src/movies/services"
)

type MovieHandle interface {
	CreateMovie(rw http.ResponseWriter, r *http.Request)
	ShowMovie(rw http.ResponseWriter, r *http.Request)
	UpdateMovie(rw http.ResponseWriter, r *http.Request)
	DeleteMovie(rw http.ResponseWriter, r *http.Request)
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
		Data: map[string]dto.MovieResponse{
			"movie": {
				ID:      movie.ID,
				Title:   movie.Title,
				Year:    movie.Year,
				Runtime: movie.Runtime,
				Genres:  movie.Genres,
				Version: movie.Version,
			},
		},
	}

	idString, _ := custom_type.EncodeId(int(movie.ID))
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/movies/%s", idString))

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

	movie, err := handler.service.GetById(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			handler.sharedUtil.NotFoundResponse(rw, r)
		default:
			handler.sharedUtil.ServerErrorResponse(rw, r, err)
		}
		return
	}

	result := commons.ResponseObject{
		StatusMsg: custom_type.Success,
		Data: map[string]dto.MovieResponse{
			"movie": {
				ID:      movie.ID,
				Title:   movie.Title,
				Runtime: movie.Runtime,
				Genres:  movie.Genres,
				Version: movie.Version,
			},
		},
	}

	err = handler.sharedUtil.WriteJson(rw, http.StatusOK, result, nil)
	if err != nil {
		handler.sharedUtil.ServerErrorResponse(rw, r, err)

		return
	}
}

func (handler *movieHandler) UpdateMovie(rw http.ResponseWriter, r *http.Request) {

	id, err := handler.sharedUtil.ExtractIdParamFromContext(r)
	if err != nil {
		handler.sharedUtil.NotFoundResponse(rw, r)

		return
	}

	movie, err := handler.service.GetById(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			handler.sharedUtil.NotFoundResponse(rw, r)
		default:
			handler.sharedUtil.ServerErrorResponse(rw, r, err)
		}
		return
	}

	var input dto.CreateMovieRequest
	err = handler.sharedUtil.ReadJson(rw, r, &input)
	if err != nil {
		handler.sharedUtil.BadRequestResponse(rw, r, err)

		return
	}

	movie.Title = input.Title
	movie.Year = input.Year
	movie.Runtime = input.Runtime
	movie.Genres = input.Genres

	validationErrors, err := handler.service.Update(movie)
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
		Data: map[string]dto.MovieResponse{
			"movie": {
				ID:      movie.ID,
				Title:   movie.Title,
				Runtime: movie.Runtime,
				Genres:  movie.Genres,
				Version: movie.Version,
			},
		},
	}

	err = handler.sharedUtil.WriteJson(rw, http.StatusOK, result, nil)
	if err != nil {
		handler.sharedUtil.ServerErrorResponse(rw, r, err)

		return
	}
}

func (handler *movieHandler) DeleteMovie(rw http.ResponseWriter, r *http.Request) {

	id, err := handler.sharedUtil.ExtractIdParamFromContext(r)
	if err != nil {
		handler.sharedUtil.NotFoundResponse(rw, r)

		return
	}

	err = handler.service.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			handler.sharedUtil.NotFoundResponse(rw, r)
		default:
			handler.sharedUtil.ServerErrorResponse(rw, r, err)
		}
		return
	}

	result := commons.ResponseObject{
		StatusMsg: custom_type.Success,
		Message:   "movie successfully deleted",
	}

	err = handler.sharedUtil.WriteJson(rw, http.StatusOK, result, nil)
	if err != nil {
		handler.sharedUtil.ServerErrorResponse(rw, r, err)

		return
	}
}
