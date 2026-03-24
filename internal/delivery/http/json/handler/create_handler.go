package handler

import (
	"encoding/json"
	"net/http"
	"tictactoe/internal/application/port"
	"tictactoe/internal/delivery/http/json/dto"
	"tictactoe/internal/delivery/http/json/mapper"
)

type CreateHandler struct {
	uc port.CreateUseCase
}

func NewCreateHandler(uc port.CreateUseCase) *CreateHandler {
	return &CreateHandler{
		uc: uc,
	}
}

func (h *CreateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	cmd := mapper.ToCreateCommand(&req)

	session, err := h.uc.Execute(r.Context(), cmd)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, mapper.ToSessionResponse(session), http.StatusOK)
}
