package handler

import (
	"encoding/json"
	"net/http"

	"tictactoe/internal/application/usecase"
	j "tictactoe/internal/delivery/http/json"
	"tictactoe/internal/delivery/http/json/dto"
	"tictactoe/internal/delivery/http/json/mapper"
)

type SignInHandler struct {
	uc usecase.SignInUseCase
	em *j.ErrorMapper
}

func NewSignInHandler(uc usecase.SignInUseCase, em *j.ErrorMapper) *SignInHandler {
	return &SignInHandler{
		uc: uc,
		em: em,
	}
}

// @Summary      Sign in
// @Description  Authenticates user by login and password, returns user UUID
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body dto.SignInRequest true "User credentials"
// @Success      200     {object} dto.UserResponse
// @Failure      400     {object} dto.ErrorResponse    "Invalid request body"
// @Failure      401     {object} dto.ErrorResponse    "Invalid credentials"
// @Router       /auth/signin [post]
func (h *SignInHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req dto.SignInRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		j.WriteError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	cmd := mapper.ToSignInCommand(&req)

	user, err := h.uc.Execute(r.Context(), cmd)
	if err != nil {
		j.WriteError(w, err.Error(), h.em.Status(err))
		return
	}

	j.WriteJSON(w, mapper.ToUserReposnse(user), http.StatusOK)
}
