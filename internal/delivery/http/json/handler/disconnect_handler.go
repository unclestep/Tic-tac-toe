package handler

import (
	"encoding/json"
	"net/http"

	"tictactoe/internal/application/port"
	j "tictactoe/internal/delivery/http/json"
	"tictactoe/internal/delivery/http/json/dto"
	"tictactoe/internal/delivery/http/json/mapper"
)

type DisconnectHandler struct {
	uc port.DisconnectUseCase
	em *j.ErrorMapper
}

func NewDisconnectHandler(uc port.DisconnectUseCase, em *j.ErrorMapper) *DisconnectHandler {
	return &DisconnectHandler{
		uc: uc,
		em: em,
	}
}

// @Summary     Disconnect a Player
// @Tags        game
// @Security    BasicAuth
// @Param  session_id  path  string  true  "Session ID"
// @Param       body       body     dto.DisconnectRequest true "Player's data"
// @Success     200        {object} dto.SessionResponse
// @Failure     400        {object} dto.ErrorResponse
// @Failure     500        {object} dto.ErrorResponse
// @Router      /game/{session_id}/disconnect [post]
func (h *DisconnectHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req dto.DisconnectRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		j.WriteError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	cmd := mapper.ToDisconnectCommand(&req)

	session, err := h.uc.Execute(r.Context(), cmd)
	if err != nil {
		j.WriteError(w, err.Error(), h.em.Status(err))
		return
	}

	j.WriteJSON(w, mapper.ToSessionResponse(session), http.StatusOK)
}
