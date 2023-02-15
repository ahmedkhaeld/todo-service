package main

import (
	"net/http"
	"sync"
)

func newMux(todoFile string) http.Handler {
	mux := http.NewServeMux()

	mu := &sync.Mutex{}

	mux.HandleFunc("/", root)

	t := todoRouter(todoFile, mu)
	//deal with the trailing slash as well as the root path
	mux.Handle("/todo", http.StripPrefix("/todo", t))
	mux.Handle("/todo/", http.StripPrefix("/todo", t))

	return mux
}
