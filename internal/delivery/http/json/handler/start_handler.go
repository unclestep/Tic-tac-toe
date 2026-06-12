package handler

import (
	"encoding/json"
	"net/http"

	"tictactoe/internal/application/port"
	j "tictactoe/internal/delivery/http/json"
	"tictactoe/internal/delivery/http/json/dto"
	"tictactoe/internal/delivery/http/json/mapper"
)

type StartHandler struct {
	uc port.StartUseCase
	em *j.ErrorMapper
}

func NewStartHandler(uc port.StartUseCase, em *j.ErrorMapper) *StartHandler {
	return &StartHandler{
		uc: uc,
		em: em,
	}
}

// @Summary     Start a Game
// @Tags        game
// @Security    BasicAuth
// @Param  session_id  path  string  true  "Session ID"
// @Param       body       body     dto.StartRequest true "Player's data"
// @Success     200        {object} dto.SessionResponse
// @Failure     400        {object} dto.ErrorResponse
// @Failure     500        {object} dto.ErrorResponse
// @Router      /game/{session_id}/start [post]
func (h *StartHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req dto.StartRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		j.WriteError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	cmd := mapper.ToStartCommand(&req)

	session, err := h.uc.Execute(r.Context(), cmd)
	if err != nil {
		j.WriteError(w, err.Error(), h.em.Status(err))
		return
	}

	j.WriteJSON(w, mapper.ToSessionResponse(session), http.StatusOK)
}
