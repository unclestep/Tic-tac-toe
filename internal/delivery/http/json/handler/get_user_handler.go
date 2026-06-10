package handler

import (
	"net/http"

	"tictactoe/internal/application/port"
	"tictactoe/internal/delivery/http/json"
	"tictactoe/internal/delivery/http/json/mapper"
)

type GetUser struct {
	userRepo port.UserRepo
	em       *json.ErrorMapper
}

func NewGetUserHandler(userRepo port.UserRepo, em *json.ErrorMapper) *GetUser {
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

	user, err := h.userRepo.Get(r.Context(), uuid)
	if err != nil {
		json.WriteError(w, err.Error(), h.em.Status(err))
		return
	}

	json.WriteJSON(w, mapper.ToUserReposnse(user), http.StatusOK)
}
