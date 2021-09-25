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

	hashId, err := custom_type.NewIdHasher()
	if err != nil {
		panic(err)
	}

	decodedId, err := hashId.DecodeWithError(stringId)
	if err != nil {
		return 0, errors.New("invalid Id")
	}

	// Convert the int32 to a ID type and assign this to the receiver. Note that we
	// use the * operator to deference the receiver (which is a pointer to a ID
	// type) in order to set the underlying value of the pointer.
	return int64(decodedId[0]), nil

}
