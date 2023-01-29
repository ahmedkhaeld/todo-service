package main

import (
	"errors"
	"net/http"
)

var content = "serve the api content here"

var (
	ErrNotFound    = errors.New("not found")
	ErrInvalidData = errors.New("invalid data")
)

func root(w http.ResponseWriter, r *http.Request) {
	//check if the client requested the root path
	if r.URL.Path != "/" {
		//if not, return a 404 not found error
		outError(w, r, http.StatusNotFound, "not found")
		return
	}

	outText(w, r, http.StatusOK, content)
}
