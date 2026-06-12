package handler

import (
	"net/http"

	"tictactoe/internal/application/port"
	j "tictactoe/internal/delivery/http/json"
	"tictactoe/internal/delivery/http/json/mapper"
)

type GetSession struct {
	sessionRepo port.SessionRepo
	em          *j.ErrorMapper
}

func NewGetSessionHandler(sessionRepo port.SessionRepo, em *j.ErrorMapper) *GetSession {
	return &GetSession{
		sessionRepo: sessionRepo,
		em:          em,
	}
}

// @Summary     Get session information
// @Tags        game
// @Security    BasicAuth
// @Param       uuid        path     string  true  "Session UUID"
// @Success     200        {object} dto.SessionResponse
// @Failure     400        {object} dto.ErrorResponse
// @Failure     500        {object} dto.ErrorResponse
// @Router      /game/{uuid} [get]
func (h *GetSession) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	UUID := r.PathValue("uuid")

	sessions, err := h.sessionRepo.Get(r.Context(), port.WithUUID(UUID))
	if err != nil {
		j.WriteError(w, err.Error(), h.em.Status(err))
		return
	}
	if len(sessions) == 0 {
		j.WriteError(w, "session not found", http.StatusNotFound)
		return
	}

	j.WriteJSON(w, mapper.ToSessionResponse(sessions[0]), http.StatusOK)
}
