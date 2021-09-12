package main

import (
	"net/http"

	"github.com/terdia/greenlight/internal/custom_type"
)

func (app *application) healthcheckHandler(rw http.ResponseWriter, r *http.Request) {

	data := responseObject{
		StatusMsg: custom_type.Success,
		Data: map[string]map[string]string{
			"system_info": {
				"environment": app.config.env,
				"version":     version,
			},
		},
	}

	err := app.writeJson(rw, http.StatusOK, data, nil)
	if err != nil {
		app.serverErrorResponse(rw, r, err)
	}
}
