package utils

import (
	"encoding/json"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	encoder.Encode(data)
}

func WriteError(w http.ResponseWriter, status int, message string) {
	response := map[string]string{
		"error": message,
	}

	WriteJSON(w, status, response)
}
