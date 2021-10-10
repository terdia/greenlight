package dto

import (
	"github.com/terdia/greenlight/internal/custom_type"
	"github.com/terdia/greenlight/internal/data"
)

type MovieRequest struct {
	Title   *string              `json:"title"`   //title for the movie, max length 500
	Year    *int32               `json:"year"`    // published year e.g. 2021, must not be in the future
	Runtime *custom_type.Runtime `json:"runtime"` // e.g 98 mins
	Genres  []string             `json:"genres"`  // unique genres e.g action,adventure... maximum 5 genres
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
