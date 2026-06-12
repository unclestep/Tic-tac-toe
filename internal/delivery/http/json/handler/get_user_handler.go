package handler

import (
	"net/http"

	"tictactoe/internal/application/port"
	j "tictactoe/internal/delivery/http/json"
	"tictactoe/internal/delivery/http/json/mapper"
)

type GetUser struct {
	userRepo port.UserRepo
	em       *j.ErrorMapper
}

func NewGetUserHandler(userRepo port.UserRepo, em *j.ErrorMapper) *GetUser {
	return &GetUser{
		userRepo: userRepo,
		em:       em,
	}
}

// @Summary     Get user information
// @Tags        user
// @Security    BasicAuth
// @Param       uuid        path     string  true  "User UUID"
// @Success     200        {object} dto.UserResponse
// @Failure     400        {object} dto.ErrorResponse
// @Failure     500        {object} dto.ErrorResponse
// @Router      /user/{uuid} [get]
func (h *GetUser) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	uuid := r.PathValue("uuid")

	users, err := h.userRepo.Get(r.Context(), port.WithUUID(uuid))
	if err != nil {
		j.WriteError(w, err.Error(), h.em.Status(err))
		return
	}
	if len(users) == 0 {
		j.WriteError(w, port.ErrUserNotFound.Error(), h.em.Status(err))
		return
	}

	j.WriteJSON(w, mapper.ToUserResponse(users[0]), http.StatusOK)
}
