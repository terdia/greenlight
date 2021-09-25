package commons

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/terdia/greenlight/internal/custom_type"
)

func (util *sharedUtils) ExtractIdParamFromContext(r *http.Request) (int64, error) {

	stringId := chi.URLParam(r, "id")
	if stringId == "" {
		return 0, errors.New("invalid parameter")
	}

	decodedId, err := custom_type.DecodeId(stringId)
	if err != nil {
		return 0, errors.New("invalid Id")
	}

	return int64(decodedId[0]), nil

}
