package jsonutil

import (
	"encoding/json"
	"net/http"
)

// Let's Go Further - Chapter 3.2
func Write(w http.ResponseWriter, status int, data any) error {
	// Encode the data to JSON, returning the error if there was one.
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Append a newline to make it easier to view in terminal applications.
	js = append(js, '\n')

	// Add the "Content-Type: application/json" header, then write the
	// status code and JSON response.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}
