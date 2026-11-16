package api

import (
	"encoding/json/v2"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	err := json.MarshalWrite(w, data)
	if err != nil {
		http.Error(w, "Error writing response", http.StatusInternalServerError)
	}
}
