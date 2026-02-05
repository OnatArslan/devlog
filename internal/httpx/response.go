package httpx

import (
	"encoding/json"
	"net/http"
)

type errorResponse struct {
	Error string `json:"error"`
}

func WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// Error helpers

func WriteError(w http.ResponseWriter, status int, err error) {
	if err == nil {
		WriteJSON(w, status, errorResponse{
			Error: http.StatusText(status),
		})
	}
	WriteJSON(w, status, errorResponse{
		Error: err.Error(),
	})
}
