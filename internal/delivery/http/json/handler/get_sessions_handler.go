package handler

import (
	"encoding/json"
	"net/http"

	"tictactoe/internal/application/port"
	j "tictactoe/internal/delivery/http/json"
	"tictactoe/internal/delivery/http/json/dto"
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

	sessions, err := h.sessionRepo.Get(r.Context(), port.WithState(domainState))
	if err != nil {
		j.WriteError(w, err.Error(), h.em.Status(err))
		return
	}

	responses := make([]dto.SessionResponse, 0, len(sessions))
	for _, s := range sessions {
		responses = append(responses, mapper.ToSessionResponse(s))
	}

	body, err := json.Marshal(responses)
	if err != nil {
		j.WriteError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}
