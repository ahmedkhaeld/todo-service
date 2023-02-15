package main

import (
	"encoding/json"
	"fmt"
	todo "github.com/ahmedkhaeld/cli-todo"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// validateID is a helper function to validate the ID in the path
// and return the ID as an integer and an error if any
func validateID(path string, list *todo.List) (int, error) {
	i := strings.Split(path, "/")
	path = i[len(i)-1]
	id, err := strconv.Atoi(path)
	if err != nil {
		return 0, fmt.Errorf("%w: Invalid ID: %s", ErrInvalidData, err)
	}

	if id < 1 {
		return 0, fmt.Errorf("%w, Invalid ID: Less than one", ErrInvalidData)
	}

	if id > len(*list) {
		return id, fmt.Errorf("%w: ID %d not found", ErrNotFound, id)
	}

	return id, nil
}

// outText is a helper function to write a string to the response
func outText(w http.ResponseWriter, r *http.Request, status int, content string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(status)
	w.Write([]byte(content))
}

// outJSON is a helper function to write a JSON response to the client
func outJSON(w http.ResponseWriter, r *http.Request, status int, data *todoResponse) {
	body, err := json.Marshal(data)
	if err != nil {
		outError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	w.Write(body)
}

// outError is a helper function to write an error to the response and log it
func outError(w http.ResponseWriter, r *http.Request, status int, msg string) {
	log.Printf("%s %s: Error: %d %s", r.URL, r.Method, status, msg)
	http.Error(w, http.StatusText(status), status)
}
