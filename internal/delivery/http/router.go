package http

import (
	"net/http"

	"tictactoe/internal/delivery/http/json/middleware"

	httpSwagger "github.com/swaggo/http-swagger"
)

type Router struct {
	mux *http.ServeMux
}

func NewRouter(create http.Handler, connect http.Handler, start http.Handler, move http.Handler, disconnect http.Handler, signUp http.Handler, signIn http.Handler, auth *middleware.UserAuthenticator) *Router {
	mux := http.NewServeMux()

	// Swagger UI
	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	// API
	mux.Handle("POST /game/create", auth.WithAuth(create))
	mux.Handle("POST /game/{uuid}/connect", auth.WithAuth(connect))
	mux.Handle("POST /game/{uuid}/start", auth.WithAuth(start))
	mux.Handle("POST /game/{uuid}/move", auth.WithAuth(move))
	mux.Handle("POST /game/{uuid}/disconnect", auth.WithAuth(disconnect))

	mux.Handle("POST /auth/signup", signUp)
	mux.Handle("POST /auth/signin", signIn)

	return &Router{mux: mux}
}

func (r *Router) Handler() http.Handler {
	return r.mux
}
