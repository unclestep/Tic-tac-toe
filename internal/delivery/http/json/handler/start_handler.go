package handler

import (
	"encoding/json"
	"net/http"

	"tictactoe/internal/application/port"
	"tictactoe/internal/delivery/http/json/dto"
	"tictactoe/internal/delivery/http/json/mapper"
)

type StartHandler struct {
	uc port.StartUseCase
}

func NewStartHandler(uc port.StartUseCase) *StartHandler {
	return &StartHandler{
		uc: uc,
	}
}

// @Summary     Start a Game
// @Tags        game
// @Param  session_id  path  string  true  "Session ID"
// @Param       body       body     dto.StartRequest true "Player's data"
// @Success     200        {object} dto.SessionResponse
// @Failure     400        {object} dto.ErrorResponse
// @Failure     500        {object} dto.ErrorResponse
// @Router      /game/{session_id}/start [post]
func (h *StartHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req dto.StartRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	cmd := mapper.ToStartCommand(&req)

	session, err := h.uc.Execute(r.Context(), cmd)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, mapper.ToSessionResponse(session), http.StatusOK)
}
