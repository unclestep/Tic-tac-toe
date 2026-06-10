package handler

import (
	"encoding/json"
	"net/http"

	"tictactoe/internal/application/usecase"
	jsonDelivery "tictactoe/internal/delivery/http/json"
	"tictactoe/internal/delivery/http/json/dto"
	"tictactoe/internal/delivery/http/json/mapper"
)

type SignUpHandler struct {
	uc usecase.SignUpUseCase
	em *jsonDelivery.ErrorMapper
}

func NewSignUpHandler(uc usecase.SignUpUseCase, em *jsonDelivery.ErrorMapper) *SignUpHandler {
	return &SignUpHandler{
		uc: uc,
		em: em,
	}
}

// @Summary      Sign up
// @Description  Registers a new user with login and password, returns created user UUID
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body dto.SignUpRequest true "User credentials"
// @Success      201     {object} dto.UserResponse
// @Failure      400     {object} dto.ErrorResponse "Invalid request body"
// @Failure      409     {object} dto.ErrorResponse "Login already taken"
// @Router       /auth/signup [post]
func (h *SignUpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req dto.SignUpRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonDelivery.WriteError(w, "invalud request body", http.StatusBadRequest)
		return
	}

	cmd := mapper.ToSignUpCommand(&req)

	user, err := h.uc.Execute(r.Context(), cmd)
	if err != nil {
		jsonDelivery.WriteError(w, err.Error(), h.em.Status(err))
		return
	}

	jsonDelivery.WriteJSON(w, mapper.ToUserReposnse(user), http.StatusOK)
}
