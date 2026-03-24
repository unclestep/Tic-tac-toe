package handler

import (
	"encoding/json"
	"net/http"
	"tictactoe/internal/application/port"
	"tictactoe/internal/delivery/http/json/dto"
	"tictactoe/internal/delivery/http/json/mapper"
)

type ConnectHandler struct {
	uc port.ConnectUseCase
}

func NewConnectHandler(uc port.ConnectUseCase) *ConnectHandler {
	return &ConnectHandler{
		uc: uc,
	}
}

func (h *ConnectHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req dto.ConnectRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	cmd := mapper.ToConnectCommand(&req)

	session, err := h.uc.Execute(r.Context(), cmd)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, mapper.ToSessionResponse(session), http.StatusOK)
}
