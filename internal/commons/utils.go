package commons

import (
	"net/http"
	"net/url"
	"sync"

	"github.com/terdia/greenlight/infrastructures/logger"
	"github.com/terdia/greenlight/internal/custom_type"
	"github.com/terdia/greenlight/internal/validator"
)

type SharedUtil interface {
	LogErrorWithHttpRequestContext(r *http.Request, err error)
	LogErrorWithContext(err error, context map[string]string)
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
	ReadString(qs url.Values, key, defaultValue string) string
	ReadInt(qs url.Values, key string, defaultValue int, v *validator.Validator) int
	ReadCSV(qs url.Values, key string, defaultValue []string) []string
	RateLimitExceededResponse(w http.ResponseWriter, r *http.Request)
	Background(fn func())
	InvalidCredentialsResponse(w http.ResponseWriter, r *http.Request)
	InvalidAuthenticationTokenResponse(w http.ResponseWriter, r *http.Request)
	InactiveAccountResponse(w http.ResponseWriter, r *http.Request)
	AuthenticationRequiredResponse(w http.ResponseWriter, r *http.Request)
	NotPermittedRResponse(w http.ResponseWriter, r *http.Request)
}

type sharedUtils struct {
	logger *logger.Logger
	wg     *sync.WaitGroup
}

func NewUtil(log *logger.Logger, wg *sync.WaitGroup) SharedUtil {
	return &sharedUtils{logger: log, wg: wg}
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
