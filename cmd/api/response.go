package main

import (
	"encoding/json"
	"net/http"

	"github.com/terdia/greenlight/internal/custom_type"
)

//based on https://github.com/omniti-labs/jsend
type responseObject struct {
	StatusMsg custom_type.StatusMessage `json:"status"` //(success|fail|error)
	Message   string                    `json:"message,omitempty"`
	Data      interface{}               `json:"data,omitempty"`
}

func (r *responseObject) setStatus(status custom_type.StatusMessage) {
	r.StatusMsg = status
}

func (app *application) writeJson(rw http.ResponseWriter, status int, envelop responseObject, headers http.Header) error {

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

func (app *application) errorResponse(rw http.ResponseWriter, r *http.Request, status int, envelop responseObject) {

	envelop.setStatus(custom_type.Error)
	err := app.writeJson(rw, status, envelop, nil)
	if err != nil {
		app.logError(r, err)
		rw.WriteHeader(http.StatusInternalServerError)
	}
}
