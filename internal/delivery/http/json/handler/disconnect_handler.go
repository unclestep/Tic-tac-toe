package handler

import (
	"encoding/json"
	"net/http"
	"tictactoe/internal/application/port"
	"tictactoe/internal/delivery/http/json/dto"
	"tictactoe/internal/delivery/http/json/mapper"
)

type DisconnectHandler struct {
	uc port.DisconnectUseCase
}

func NewDisconnectHandler(uc port.DisconnectUseCase) *DisconnectHandler {
	return &DisconnectHandler{
		uc: uc,
	}
}

func (h *DisconnectHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req dto.DisconnectRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	cmd := mapper.ToDisconnectCommand(&req)

	session, err := h.uc.Execute(r.Context(), cmd)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, mapper.ToSessionResponse(session), http.StatusOK)
}
