package main

import (
	"fmt"
	"net/http"
)

func (app *application) healthcheckHandler(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(rw, "status: available")
	fmt.Fprintf(rw, "environment: %s\n", app.config.env)
	fmt.Fprintf(rw, "version: %s\n", version)
}
