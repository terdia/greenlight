package main

import (
	"fmt"
	"net/http"
)

func (app *application) healthcheckHandler(rw http.ResponseWriter, r *http.Request) {

	js := `{"status": "available", "environment": %q, "version": %q}`
	js = fmt.Sprintf(js, app.config.env, version)

	rw.Header().Set("Content-Type", "application/json")

	rw.Write([]byte(js))
}
