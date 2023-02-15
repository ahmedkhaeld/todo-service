package main

import (
	"encoding/json"
	"errors"
	"fmt"
	todo "github.com/ahmedkhaeld/cli-todo"
	"net/http"
	"sync"
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

////////////////////////todoHandlers////////////////////////

// todoRouter is grouped routes for the todo service
func todoRouter(todoFile string, l sync.Locker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		list := &todo.List{}

		l.Lock()
		defer l.Unlock()
		if err := list.Get(todoFile); err != nil {
			outError(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		if r.URL.Path == "" {
			switch r.Method {
			case http.MethodGet:
				getAllHandler(w, r, list)
			case http.MethodPost:
				addHandler(w, r, list, todoFile)
			default:
				outError(w, r, http.StatusMethodNotAllowed, "method not allowed")
			}
			return // return here to avoid the next if statement
		}

		id, err := validateID(r.URL.Path, list)
		if err != nil {
			if errors.Is(err, ErrNotFound) {
				outError(w, r, http.StatusNotFound, err.Error())
				return
			}
			outError(w, r, http.StatusBadRequest, err.Error())
			return
		}

		switch r.Method {
		case http.MethodGet:
			getOneHandler(w, r, list, id)
		case http.MethodDelete:
			deleteHandler(w, r, list, id, todoFile)
		case http.MethodPatch:
			patchHandler(w, r, list, id, todoFile)
		default:
			message := "Method not supported"
			outError(w, r, http.StatusMethodNotAllowed, message)
		}

	}
}

func getAllHandler(w http.ResponseWriter, r *http.Request, list *todo.List) {
	resp := &todoResponse{
		Results: *list,
	}
	outJSON(w, r, http.StatusOK, resp)
}

func getOneHandler(w http.ResponseWriter, r *http.Request, list *todo.List, id int) {

	resp := &todoResponse{
		Results: (*list)[id-1 : id],
	}
	outJSON(w, r, http.StatusOK, resp)
}

func deleteHandler(w http.ResponseWriter, r *http.Request, list *todo.List, id int, todoFile string) {

	list.Delete(id)
	if err := list.Save(todoFile); err != nil {
		outError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	outText(w, r, http.StatusNoContent, "")
}

func patchHandler(w http.ResponseWriter, r *http.Request, list *todo.List, id int, todoFile string) {

	q := r.URL.Query()

	if _, ok := q["complete"]; !ok {
		message := "Missing query param 'complete'"
		outError(w, r, http.StatusBadRequest, message)
		return
	}

	list.Complete(id)
	if err := list.Save(todoFile); err != nil {
		outError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	outText(w, r, http.StatusNoContent, "")
}

func addHandler(w http.ResponseWriter, r *http.Request, list *todo.List, todoFile string) {

	item := struct {
		Task string `json:"task"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		message := fmt.Sprintf("Invalid JSON: %s", err)
		outError(w, r, http.StatusBadRequest, message)
		return
	}

	list.Add(item.Task)
	if err := list.Save(todoFile); err != nil {
		outError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	outText(w, r, http.StatusCreated, "")
}
