package handler

import (
	"encoding/json"
	"net/http"
	"tictactoe/internal/application/port"
	"tictactoe/internal/delivery/http/json/dto"
	"tictactoe/internal/delivery/http/json/mapper"
)

type MakeMoveHandler struct {
	uc port.MakeMoveUseCase
}

func NewMakeMoveHandler(uc port.MakeMoveUseCase) *MakeMoveHandler {
	return &MakeMoveHandler{
		uc: uc,
	}
}

func (h *MakeMoveHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req dto.MakeMoveRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	cmd := mapper.ToMakeMoveCommand(&req)

	session, err := h.uc.Execute(r.Context(), cmd)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, mapper.ToSessionResponse(session), http.StatusOK)
}
