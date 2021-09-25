package commons

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/terdia/greenlight/internal/custom_type"
)

func (util *sharedUtils) WriteJson(rw http.ResponseWriter, status int, envelop ResponseObject, headers http.Header) error {

	js, err := json.MarshalIndent(envelop, "", "\t")
	if err != nil {
		return err
	}

	js = append(js, '\n')

	for key, value := range headers {
		rw.Header()[key] = value
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(status)
	rw.Write(js)

	return nil
}

func (util *sharedUtils) ErrorResponse(rw http.ResponseWriter, r *http.Request, status int, envelop ResponseObject) {

	if envelop.StatusMsg == 0 {
		envelop.setStatus(custom_type.Error)
	}

	err := util.WriteJson(rw, status, envelop, nil)
	if err != nil {

		util.LogError(r, err)
		rw.WriteHeader(http.StatusInternalServerError)
	}
}

func (util *sharedUtils) ReadJson(rw http.ResponseWriter, r *http.Request, dst interface{}) error {

	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(rw, r.Body, int64(maxBytes))

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(dst)
	if err != nil {

		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {

		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)

		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")

		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)

		case errors.Is(err, io.EOF):
			return errors.New("request boby must not be empty")

			// If the JSON contains a field which cannot be mapped to the target destination
			// then Decode() will now return an error message in the format "json: unknown
			// field "<name>"". We check for this, extract the field name from the error,
			// and interpolate it into our custom error message. Note that there's an open
			// issue at https://github.com/golang/go/issues/29035 regarding turning this
			// into a distinct error type in the future.
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("request body contains unknown key %s", fieldName)

			// If the request body exceeds 1MB in size the decode will now fail with the
			// error "http: request body too large". There is an open issue about turning
			// this into a distinct error type at https://github.com/golang/go/issues/30715.
		case err.Error() == "http: request body too large":
			return fmt.Errorf("body must not be larger than %d bytes", maxBytes)

		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		default:
			return err

		}
	}

	// Call Decode() again, using a pointer to an empty anonymous struct as the
	// destination. If the request body only contained a single JSON value this will
	// return an io.EOF error. So if we get anything else, we know that there is
	// additional data in the request body and we return our own custom error message.
	err = decoder.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only contain a single JSON value")
	}

	return nil
}
