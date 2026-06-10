package handler

import (
	"encoding/json"
	"net/http"

	"tictactoe/internal/application/port"
	jsonDelivery "tictactoe/internal/delivery/http/json"
	"tictactoe/internal/delivery/http/json/dto"
	"tictactoe/internal/delivery/http/json/mapper"
)

type CreateHandler struct {
	uc port.CreateUseCase
	em *jsonDelivery.ErrorMapper
}

func NewCreateHandler(uc port.CreateUseCase, em *jsonDelivery.ErrorMapper) *CreateHandler {
	return &CreateHandler{
		uc: uc,
		em: em,
	}
}

// @Summary     Create a Game
// @Tags        game
// @Security    BasicAuth
// @Param       body       body     dto.CreateRequest true "Game params"
// @Success     200        {object} dto.SessionResponse
// @Failure     400        {object} dto.ErrorResponse
// @Failure     500        {object} dto.ErrorResponse
// @Router      /game/create [post]
func (h *CreateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonDelivery.WriteError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	cmd := mapper.ToCreateCommand(&req)

	session, err := h.uc.Execute(r.Context(), cmd)
	if err != nil {
		jsonDelivery.WriteError(w, err.Error(), h.em.Status(err))
		return
	}

	jsonDelivery.WriteJSON(w, mapper.ToSessionResponse(session), http.StatusOK)
}
