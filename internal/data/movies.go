package data

import (
	"time"

	"github.com/terdia/greenlight/internal/custom_type"
)

type Movie struct {
	ID        int64               `json:"id"`
	Title     string              `json:"title"`
	Year      int32               `json:"year,omitempty"`
	Runtime   custom_type.Runtime `json:"runtime,omitempty"`
	Genres    []string            `json:"genres,omitempty"`
	Version   int32               `json:"version"`
	CreatedAt time.Time           `json:"-"` // exclude from reponse
}

type CreateMovieRequest struct {
	Title   string              `json:"title"`
	Year    int32               `json:"year"`
	Runtime custom_type.Runtime `json:"runtime"`
	Genres  []string            `json:"genres"`
}
