package handler

import (
	"net/http"

	"tictactoe/internal/application/port"
	j "tictactoe/internal/delivery/http/json"
	"tictactoe/internal/delivery/http/json/mapper"
)

type GetSessions struct {
	sessionRepo port.SessionRepo
	em          *j.ErrorMapper
}

func NewGetSessionsHandler(sessionRepo port.SessionRepo, em *j.ErrorMapper) *GetSessions {
	return &GetSessions{
		sessionRepo: sessionRepo,
		em:          em,
	}
}

// @Summary     Get sessions by state
// @Tags        game
// @Security    BasicAuth
// @Param       state       query    string  false "Session state filter (e.g. Playing)"
// @Success     200        {array}  dto.SessionResponse
// @Failure     400        {object} dto.ErrorResponse
// @Failure     500        {object} dto.ErrorResponse
// @Router      /game [get]
func (h *GetSessions) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rawState := r.URL.Query().Get("state")

	domainState, err := mapper.StringToState(rawState)
	if err != nil {
		j.WriteError(w, err.Error(), http.StatusBadRequest)
		return
	}

	var opt port.SessionGetOpt
	if rawState != "" {
		opt = port.WithState(domainState)
	}
	sessions, err := h.sessionRepo.Get(r.Context(), opt)
	if err != nil {
		j.WriteError(w, err.Error(), h.em.Status(err))
		return
	}

	j.WriteJSON(w, sessions, http.StatusOK)
}
