package main

import (
	"net/http"

	"github.com/terdia/greenlight/internal/commons"
	"github.com/terdia/greenlight/internal/custom_type"
)

func (app *application) healthcheckHandler(rw http.ResponseWriter, r *http.Request) {

	data := commons.ResponseObject{
		StatusMsg: custom_type.Success,
		Data: map[string]map[string]string{
			"system_info": {
				"environment": app.config.Env,
				"version":     app.config.Version,
			},
		},
	}

	utils := app.registry.Services.SharedUtil

	err := utils.WriteJson(rw, http.StatusOK, data, nil)
	if err != nil {
		utils.ServerErrorResponse(rw, r, err)
	}
}
