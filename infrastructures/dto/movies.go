package dto

import (
	"github.com/terdia/greenlight/internal/custom_type"
)

type CreateMovieRequest struct {
	Title   string              `json:"title"`
	Year    int32               `json:"year"`
	Runtime custom_type.Runtime `json:"runtime"`
	Genres  []string            `json:"genres"`
}

type MovieResponse struct {
	ID      custom_type.ID      `json:"id"`
	Title   string              `json:"title"`
	Year    int32               `json:"year,omitempty"`
	Runtime custom_type.Runtime `json:"runtime,omitempty"`
	Genres  []string            `json:"genres,omitempty"`
	Version int32               `json:"version"`
}