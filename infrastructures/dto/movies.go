package dto

import (
	"github.com/terdia/greenlight/internal/custom_type"
	"github.com/terdia/greenlight/internal/data"
)

type MovieRequest struct {
	Title   *string              `json:"title"`
	Year    *int32               `json:"year"`
	Runtime *custom_type.Runtime `json:"runtime"`
	Genres  []string             `json:"genres"`
}

type SingleMovieResponse struct {
	Movie MovieResponse `json:"movie"`
}

type MovieResponse struct {
	ID      custom_type.ID      `json:"id"`
	Title   string              `json:"title"`
	Year    int32               `json:"year,omitempty"`
	Runtime custom_type.Runtime `json:"runtime,omitempty"`
	Genres  []string            `json:"genres,omitempty"`
	Version int32               `json:"version"`
}

type ListMovieRequest struct {
	Title   string
	Genres  []string
	Filters data.Filters
}

type ListMovieResponse struct {
	Metadata data.Metadata   `json:"metadata"`
	Movies   []MovieResponse `json:"movies"`
}

type ValidationError struct {
	Errors map[string]string `json:"errors"`
}
