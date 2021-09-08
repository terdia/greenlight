package main

import (
	"fmt"
	"net/http"
)

func (app *application) createMovieHandler(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(rw, "create a new movie")
}

func (app *application) showMovieHandler(rw http.ResponseWriter, r *http.Request) {

	id, err := app.extractIdParamFromContext(r)
	if err != nil {
		http.NotFound(rw, r)
		return
	}

	fmt.Fprintf(rw, "show details of movie %d\n", id)
}
