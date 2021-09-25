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
