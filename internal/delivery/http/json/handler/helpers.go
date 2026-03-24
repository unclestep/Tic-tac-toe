package handler

import (
	"encoding/json"
	"net/http"
	"tictactoe/internal/delivery/http/json/dto"
)

func writeJSON(w http.ResponseWriter, v any, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func writeError(w http.ResponseWriter, msg string, status int) {
	writeJSON(w, dto.ErrorResponse{Error: msg}, status)
}
