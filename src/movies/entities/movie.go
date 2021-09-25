package entities

import (
	"time"

	"github.com/terdia/greenlight/internal/custom_type"
)

type Movie struct {
	ID        custom_type.ID
	Title     string
	Year      int32
	Runtime   custom_type.Runtime
	Genres    []string
	Version   int32
	CreatedAt time.Time
}
