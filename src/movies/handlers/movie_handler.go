package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/terdia/greenlight/infrastructures/dto"
	"github.com/terdia/greenlight/internal/commons"
	"github.com/terdia/greenlight/internal/custom_type"
	"github.com/terdia/greenlight/internal/data"
	"github.com/terdia/greenlight/internal/validator"
	"github.com/terdia/greenlight/src/movies/entities"
	"github.com/terdia/greenlight/src/movies/services"
)

type MovieHandle interface {
	CreateMovie(rw http.ResponseWriter, r *http.Request)
	ShowMovie(rw http.ResponseWriter, r *http.Request)
	UpdateMovie(rw http.ResponseWriter, r *http.Request)
	DeleteMovie(rw http.ResponseWriter, r *http.Request)
	ListMovie(rw http.ResponseWriter, r *http.Request)
}

type movieHandler struct {
	sharedUtil commons.SharedUtil
	service    services.MovieService
}

func NewMovieHandler(util commons.SharedUtil, service services.MovieService) MovieHandle {
	return &movieHandler{
		sharedUtil: util,
		service:    service,
	}
}

// CreateMovie ... Create movie
// @Summary Create new movie
// @Description create a new movie with given details
// @Tags Movies
// @Param body body dto.MovieRequest true "Update movie request"
// @Param Authorization header string true "Authorization: Bearer XXSGGSSHHSSJSJSSS"
// @Success 200 {object} commons.ResponseObject{data=dto.SingleMovieResponse}
// @Header 200 {string} Location "/v1/movies/QbPy4B7a2Lw1Kg7ogoEWj9k3NGMRVY"
// @Failure 422 {object} commons.ResponseObject{data=dto.ValidationError} "status: fail"
// @Failure 400,401,403,500 {object} commons.ResponseObject "e.g. status: error, message: the error reason"
// @Router /movies [post]
func (handler *movieHandler) CreateMovie(rw http.ResponseWriter, r *http.Request) {
	var input dto.MovieRequest

	err := handler.sharedUtil.ReadJson(rw, r, &input)
	if err != nil {
		handler.sharedUtil.BadRequestResponse(rw, r, err)

		return
	}

	movie := &entities.Movie{
		Title:   *input.Title,
		Year:    *input.Year,
		Runtime: *input.Runtime,
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
		Data: dto.SingleMovieResponse{
			Movie: getMovieResponse(movie),
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

// ShowMovie ... Show movie
// @Summary Show movie details by id
// @Description show details of a given movie
// @Tags Movies
// @Param id path string false "Id of the movie to show"
// @Param Authorization header string true "Authorization: Bearer XXSGGSSHHSSJSJSSS"
// @Success 200 {object} commons.ResponseObject{data=dto.SingleMovieResponse}
// @Failure 400,401,403,500 {object} commons.ResponseObject "e.g. status: error, message: the error reason"
// @Router /movies/{id} [get]
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
		Data: dto.SingleMovieResponse{
			Movie: getMovieResponse(movie),
		},
	}

	err = handler.sharedUtil.WriteJson(rw, http.StatusOK, result, nil)
	if err != nil {
		handler.sharedUtil.ServerErrorResponse(rw, r, err)

		return
	}
}

// UpdateMovie ... Update movie
// @Summary Update a given movie
// @Description update movie with given details
// @Tags Movies
// @Param id path string true "Id of the movie to update"
// @Param body body dto.MovieRequest false "Update movie request"
// @Param Authorization header string true "Authorization: Bearer XXSGGSSHHSSJSJSSS"
// @Success 200 {object} commons.ResponseObject{data=dto.SingleMovieResponse}
// @Header 200 {string} Location "/v1/movies/QbPy4B7a2Lw1Kg7ogoEWj9k3NGMRVY"
// @Failure 409 {object} commons.ResponseObject "e.g. status: error, message: unable to update the record due to an edit conflict, please try again"
// @Failure 422 {object} commons.ResponseObject{data=dto.ValidationError} "status: fail"
// @Failure 400,401,403,404,500 {object} commons.ResponseObject "e.g. status: error, message: the error reason"
// @Router /movies/{id} [patch]
func (handler *movieHandler) UpdateMovie(rw http.ResponseWriter, r *http.Request) {

	id, err := handler.sharedUtil.ExtractIdParamFromContext(r)
	if err != nil {
		handler.sharedUtil.NotFoundResponse(rw, r)

		return
	}

	var input dto.MovieRequest
	err = handler.sharedUtil.ReadJson(rw, r, &input)
	if err != nil {
		handler.sharedUtil.BadRequestResponse(rw, r, err)

		return
	}

	movie, validationErrors, err := handler.service.Update(id, input)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			handler.sharedUtil.NotFoundResponse(rw, r)
		case errors.Is(err, data.ErrEditConflict):
			handler.sharedUtil.EditConflictResponse(rw, r)
		default:
			handler.sharedUtil.ServerErrorResponse(rw, r, err)
		}
		return
	}

	if validationErrors != nil {
		handler.sharedUtil.FailedValidationResponse(rw, r, validationErrors)

		return
	}

	result := commons.ResponseObject{
		StatusMsg: custom_type.Success,
		Data: dto.SingleMovieResponse{
			Movie: getMovieResponse(movie),
		},
	}

	err = handler.sharedUtil.WriteJson(rw, http.StatusOK, result, nil)
	if err != nil {
		handler.sharedUtil.ServerErrorResponse(rw, r, err)

		return
	}
}

// DeleteMovie ... Delete a given movie
// @Summary Delete a given movie
// @Description delete a given movie by Id
// @Tags Movies
// @Param id path string false "Id of the movie to delete"
// @Param Authorization header string true "Authorization: Bearer XXSGGSSHHSSJSJSSS"
// @Success 200 {object} commons.ResponseObject
// @Failure 401,403,404,500 {object} commons.ResponseObject "e.g. status: error, message: the error reason"
// @Router /movies/{id} [delete]
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

// ListMovie ... Get all movies
// @Summary Get all movies
// @Description get all movies
// @Tags Movies
// @Param title query string false "full text search by movie title"
// @Param genres query string false "command seperated list e.g. crime,drama"
// @Param page query integer false "page number"  default(1) minimum(1) maximum(10000000)
// @Param page_size query integer false "page size" default(10) minimum(1) maximum(100)
// @Param sort query string false "add - to sort in descing order" Enums(id, title, year, runtime, -id, -title, -year, -runtime) default(id)
// @Param Authorization header string true "Authorization: Bearer XXSGGSSHHSSJSJSSS"
// @Success 200 {object} commons.ResponseObject{data=dto.ListMovieResponse}
// @Failure 422 {object} commons.ResponseObject{data=dto.ValidationError} "status: fail"
// @Failure 401,403,500 {object} commons.ResponseObject "e.g. status: error, message: the error reason"
// @Router /movies [get]
func (handler *movieHandler) ListMovie(rw http.ResponseWriter, r *http.Request) {
	util := handler.sharedUtil
	v := validator.New()

	qs := r.URL.Query()

	filters := data.Filters{
		Page:         util.ReadInt(qs, "page", 1, v),
		PageSize:     util.ReadInt(qs, "page_size", 10, v),
		Sort:         util.ReadString(qs, "sort", "id"),
		SortSafelist: []string{"id", "title", "year", "runtime", "-id", "-title", "-year", "-runtime"},
	}

	filters.ValidateFilters(v)
	if !v.Valid() {
		util.FailedValidationResponse(rw, r, v.Errors)
		return
	}

	listMoviesRequest := dto.ListMovieRequest{
		Title:   util.ReadString(qs, "title", ""),
		Genres:  util.ReadCSV(qs, "genres", []string{}),
		Filters: filters,
	}

	movies, metadata, err := handler.service.List(listMoviesRequest)
	if err != nil {
		util.ServerErrorResponse(rw, r, err)
		return
	}

	moviesDto := []dto.MovieResponse{}
	for _, movie := range movies {
		moviesDto = append(moviesDto, getMovieResponse(movie))
	}

	err = handler.sharedUtil.WriteJson(rw, http.StatusOK, commons.ResponseObject{
		StatusMsg: custom_type.Success,
		Data: dto.ListMovieResponse{
			Metadata: metadata,
			Movies:   moviesDto,
		},
	}, nil)

	if err != nil {
		handler.sharedUtil.ServerErrorResponse(rw, r, err)

		return
	}
}

func getMovieResponse(movie *entities.Movie) dto.MovieResponse {
	return dto.MovieResponse{
		ID:      movie.ID,
		Title:   movie.Title,
		Year:    movie.Year,
		Runtime: movie.Runtime,
		Genres:  movie.Genres,
		Version: movie.Version,
	}
}
