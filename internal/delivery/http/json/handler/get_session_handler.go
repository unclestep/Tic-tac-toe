package handler

import (
	"net/http"

	"tictactoe/internal/application/port"
	"tictactoe/internal/delivery/http/json"
	"tictactoe/internal/delivery/http/json/mapper"
)

type GetSession struct {
	sessionRepo port.SessionRepo
	em          *json.ErrorMapper
}

func NewGetSessionHandler(sessionRepo port.SessionRepo, em *json.ErrorMapper) *GetSession {
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
	if err != nil || len(sessions) != 1 {
		json.WriteError(w, err.Error(), h.em.Status(err))
		return
	}
	session := sessions[0]

	json.WriteJSON(w, mapper.ToSessionResponse(session), http.StatusOK)
}
