package main

import (
	"net/http"
)

func newMux(todoFile string) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", root)

	return mux
}
