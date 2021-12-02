package api

import (
	"encoding/json"
	"net/http"
)

type envelope map[string]interface{}

func writeJSON(w http.ResponseWriter, status int, data envelope) error {
	// attempt to encode data into JSON
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	// append a newline for nicer terminal output
	js = append(js, '\n')

	// set content type, set status, and write the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}
