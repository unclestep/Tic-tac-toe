package middleware

import (
	"context"
	"net/http"

	"tictactoe/internal/application/port"

	"tictactoe/internal/application/usecase"
	"tictactoe/internal/delivery/http/json"
)

type UserAuthenticator struct {
	signIn usecase.SignInUseCase
}

func NewUserAuthenticator(signIn usecase.SignInUseCase) *UserAuthenticator {
	return &UserAuthenticator{
		signIn: signIn,
	}
}

func (m *UserAuthenticator) WithAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		login, password, ok := r.BasicAuth()
		if !ok {
			json.WriteError(w, "no auth credentials passed", http.StatusUnauthorized)
			return
		}

		user, err := m.signIn.Execute(
			r.Context(),
			&usecase.SignInCommand{
				Login:    login,
				Password: password,
			},
		)

		newCtx := context.WithValue(r.Context(), port.UserUUIDKey, user.UUID)
		reqWithCtx := r.WithContext(newCtx)

		if err != nil {
			json.WriteError(w, err.Error(), http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, reqWithCtx)
	})
}
