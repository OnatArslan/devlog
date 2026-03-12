package httpx

import (
	"encoding/json"
	"net/http"
)

type errorResponse struct {
	Error string `json:"error"`
}

// WriteJSON encodes data as JSON and writes it with the given status code.
func WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// WriteError writes a JSON error response with the given status code and error message.
func WriteError(w http.ResponseWriter, status int, err error) {
	if err == nil {
		WriteJSON(w, status, errorResponse{
			Error: http.StatusText(status),
		})
		return
	}
	WriteJSON(w, status, errorResponse{
		Error: err.Error(),
	})
}
