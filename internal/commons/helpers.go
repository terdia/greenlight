package commons

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/terdia/greenlight/internal/custom_type"
	"github.com/terdia/greenlight/internal/validator"
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

func (util *sharedUtils) ReadString(qs url.Values, key, defaultValue string) string {

	str := qs.Get(key)
	if str == "" {
		return defaultValue
	}

	return str

}

func (util *sharedUtils) ReadInt(qs url.Values, key string, defaultValue int, v *validator.Validator) int {

	s := qs.Get(key)
	if s == "" {
		return defaultValue
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		v.AddError(key, "must be an integer value")
		return defaultValue
	}

	return i

}

func (util *sharedUtils) ReadCSV(qs url.Values, key string, defaultValue []string) []string {

	csv := qs.Get(key)
	if csv == "" {
		return defaultValue
	}

	return strings.Split(csv, ",")

}

func (util *sharedUtils) Background(fn func()) {

	util.wg.Add(1)

	go func() {

		defer util.wg.Done()

		//time.Sleep(5 * time.Second)

		defer func() {
			if err := recover(); err != nil {
				util.LogErrorWithContext(fmt.Errorf("%s", err), nil)
			}
		}()

		fn()
	}()

}
