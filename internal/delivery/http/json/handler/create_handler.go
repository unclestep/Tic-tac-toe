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
	em *ErrorMapper
}

func NewCreateHandler(uc port.CreateUseCase, em *ErrorMapper) *CreateHandler {
	return &CreateHandler{
		uc: uc,
		em: em,
	}
}

// @Summary     Create a Game
// @Tags        game
// @Param       body       body     dto.CreateRequest true "Game params"
// @Success     200        {object} dto.SessionResponse
// @Failure     400        {object} dto.ErrorResponse
// @Failure     500        {object} dto.ErrorResponse
// @Router      /game/create [post]
func (h *CreateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	cmd := mapper.ToCreateCommand(&req)

	session, err := h.uc.Execute(r.Context(), cmd)
	if err != nil {
		writeError(w, err.Error(), h.em.Status(err))
		return
	}

	writeJSON(w, mapper.ToSessionResponse(session), http.StatusOK)
}
