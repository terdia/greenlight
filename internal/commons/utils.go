package commons

import (
	"log"
	"net/http"

	"github.com/terdia/greenlight/internal/custom_type"
)

type SharedUtil interface {
	LogError(r *http.Request, err error)
	ServerErrorResponse(rw http.ResponseWriter, r *http.Request, err error)
	NotFoundResponse(rw http.ResponseWriter, r *http.Request)
	EditConflictResponse(rw http.ResponseWriter, r *http.Request)
	MethodNotAllowedResponse(w http.ResponseWriter, r *http.Request)
	BadRequestResponse(w http.ResponseWriter, r *http.Request, err error)
	FailedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string)
	WriteJson(rw http.ResponseWriter, status int, envelop ResponseObject, headers http.Header) error
	ReadJson(rw http.ResponseWriter, r *http.Request, dst interface{}) error
	ErrorResponse(rw http.ResponseWriter, r *http.Request, status int, envelop ResponseObject)
	ExtractIdParamFromContext(r *http.Request) (int64, error)
}

type sharedUtils struct {
	logger *log.Logger
}

func NewUtil(log *log.Logger) SharedUtil {
	return &sharedUtils{logger: log}
}

//based on https://github.com/omniti-labs/jsend
type ResponseObject struct {
	StatusMsg custom_type.StatusMessage `json:"status"` //(success|fail|error)
	Message   string                    `json:"message,omitempty"`
	Data      interface{}               `json:"data,omitempty"`
}

func (r *ResponseObject) setStatus(status custom_type.StatusMessage) {
	r.StatusMsg = status
}
