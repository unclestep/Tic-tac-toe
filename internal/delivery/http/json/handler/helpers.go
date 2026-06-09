package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"tictactoe/internal/application/port"
	"tictactoe/internal/delivery/http/json/dto"
	"tictactoe/internal/domain/model"
)

type errorEntry struct {
	err    error
	status int
}

type ErrorMapper struct {
	entries []errorEntry
}

func NewErrorMapper() *ErrorMapper {
	return &ErrorMapper{
		entries: []errorEntry{
			{port.ErrSessionNotFound, http.StatusNotFound},
			{port.ErrRulesNotFound, http.StatusNotFound},
			{model.ErrPlayerNotFound, http.StatusNotFound},
			{model.ErrSessionFull, http.StatusBadRequest},
			{model.ErrPlayerExists, http.StatusBadRequest},
			{model.ErrMarkTaken, http.StatusBadRequest},
			{model.ErrGameAlreadyOver, http.StatusBadRequest},
			{model.ErrCellNotEmpty, http.StatusBadRequest},
			{model.ErrOutOfBounds, http.StatusBadRequest},
			{port.ErrGameNotStarted, http.StatusBadRequest},
			{port.ErrGameAlreadyStarted, http.StatusBadRequest},
			{port.ErrPlayerNotBelongToSession, http.StatusConflict},
			{port.ErrPlayerCantMakeMove, http.StatusConflict},
		},
	}
}

func (m *ErrorMapper) Status(err error) int {
	for _, entry := range m.entries {
		if errors.Is(err, entry.err) {
			return entry.status
		}
	}
	return http.StatusInternalServerError
}

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
