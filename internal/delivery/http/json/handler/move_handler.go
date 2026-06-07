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
	em *ErrorMapper
}

func NewMakeMoveHandler(uc port.MakeMoveUseCase, em *ErrorMapper) *MakeMoveHandler {
	return &MakeMoveHandler{
		uc: uc,
		em: em,
	}
}

// @Summary     Make a Move
// @Tags        game
// @Param  session_id  path  string  true  "Session ID"
// @Param       body       body     dto.MakeMoveRequest true "Move data"
// @Success     200        {object} dto.SessionResponse
// @Failure     400        {object} dto.ErrorResponse
// @Failure     500        {object} dto.ErrorResponse
// @Router      /game/{session_id}/move [post]
func (h *MakeMoveHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req dto.MakeMoveRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	cmd := mapper.ToMakeMoveCommand(&req)

	session, err := h.uc.Execute(r.Context(), cmd)
	if err != nil {
		writeError(w, err.Error(), h.em.Status(err))
		return
	}

	writeJSON(w, mapper.ToSessionResponse(session), http.StatusOK)
}
