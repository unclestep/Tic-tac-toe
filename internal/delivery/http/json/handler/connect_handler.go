package handler

import (
	"encoding/json"
	"net/http"

	"tictactoe/internal/application/port"
	jsonDelivery "tictactoe/internal/delivery/http/json"
	"tictactoe/internal/delivery/http/json/dto"
	"tictactoe/internal/delivery/http/json/mapper"
)

type ConnectHandler struct {
	uc port.ConnectUseCase
	em *jsonDelivery.ErrorMapper
}

func NewConnectHandler(uc port.ConnectUseCase, em *jsonDelivery.ErrorMapper) *ConnectHandler {
	return &ConnectHandler{
		uc: uc,
		em: em,
	}
}

// @Summary     Connect a Player
// @Tags        game
// @Security    BasicAuth
// @Param  session_uuid  path  string  true  "Session UUID"
// @Param       body       body     dto.ConnectRequest true "Player's data"
// @Success     200        {object} dto.SessionResponse
// @Failure     400        {object} dto.ErrorResponse
// @Failure     500        {object} dto.ErrorResponse
// @Router      /game/{session_id}/connect [post]
func (h *ConnectHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req dto.ConnectRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonDelivery.WriteError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	cmd := mapper.ToConnectCommand(&req)

	session, err := h.uc.Execute(r.Context(), cmd)
	if err != nil {
		jsonDelivery.WriteError(w, err.Error(), h.em.Status(err))
		return
	}

	jsonDelivery.WriteJSON(w, mapper.ToSessionResponse(session), http.StatusOK)
}
