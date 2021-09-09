package main

import (
	"net/http"
)

func (app *application) healthcheckHandler(rw http.ResponseWriter, r *http.Request) {

	data := responseData{
		"status": "success",
		"data": map[string]map[string]string{
			"system_info": {
				"environment": app.config.env,
				"version":     version,
			},
		},
	}

	err := app.writeJson(rw, http.StatusOK, data, nil)
	if err != nil {
		app.logger.Println(err)
		http.Error(rw, "Server cannot process your request", http.StatusInternalServerError)

		return
	}
}
