package main

import (
	"encoding/json"
	todo "github.com/ahmedkhaeld/cli-todo"
	"time"
)

// todoResponse wraps the list of to-do items
// with an exported field Results
type todoResponse struct {
	Results todo.List `json:"result"`
}

// MarshalJSON is a custom JSON marshaller for the todoResponse type
// It adds the date and total_results fields to the JSON response
//
// Note: JSON customization in Go, you can change the spelling of fields,
// remap them to other fields, or even omit them from the JSON output.
func (r *todoResponse) MarshalJSON() ([]byte, error) {
	resp := struct {
		Results      todo.List `json:"result"`
		Date         int64     `json:"date"`
		TotalResults int       `json:"total_results"`
	}{
		Results:      r.Results,
		Date:         time.Now().Unix(),
		TotalResults: len(r.Results),
	}
	return json.Marshal(resp)
}
