package data

import (
	"time"

	"github.com/terdia/greenlight/internal/custom_type"
	"github.com/terdia/greenlight/internal/validator"
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

func (movie *Movie) Validate(v *validator.Validator) {

	v.Check(movie.Title != "", "title", "must be provided")
	v.Check(len(movie.Title) <= 500, "title", "must not be more than 500 bytes long")
	v.Check(movie.Year != 0, "year", "must be provided")
	v.Check(movie.Year >= 1888, "year", "must be greater than 1888")
	v.Check(movie.Year <= int32(time.Now().Year()), "year", "must not be in the future")
	v.Check(movie.Runtime != 0, "runtime", "must be provided")
	v.Check(movie.Runtime > 0, "runtime", "must be a positive integer")
	v.Check(movie.Genres != nil, "genres", "must be provided")
	v.Check(len(movie.Genres) >= 1, "genres", "must contain at least 1 genre")
	v.Check(len(movie.Genres) <= 5, "genres", "must not contain more than 5 genres")
	v.Check(validator.UniqueStringSlice(movie.Genres), "genres", "must not contain duplicate values")
}
