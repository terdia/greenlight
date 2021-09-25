package entities

import (
	"time"

	"github.com/terdia/greenlight/internal/custom_type"
)

type Movie struct {
	ID        custom_type.ID      `json:"id"`
	Title     string              `json:"title"`
	Year      int32               `json:"year,omitempty"`
	Runtime   custom_type.Runtime `json:"runtime,omitempty"`
	Genres    []string            `json:"genres,omitempty"`
	Version   int32               `json:"version"`
	CreatedAt time.Time           `json:"-"` // exclude from reponse
}
