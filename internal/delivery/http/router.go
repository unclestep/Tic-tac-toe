package http

import (
	"net/http"

	"tictactoe/internal/delivery/http/json/middleware"

	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/fx"
)

type Router struct {
	mux *http.ServeMux
}

type RouterParams struct {
	fx.In

	Create      http.Handler `name:"create_h"`
	Connect     http.Handler `name:"connect_h"`
	Start       http.Handler `name:"start_h"`
	Move        http.Handler `name:"move_h"`
	Disconnect  http.Handler `name:"disconnect_h"`
	SignUp      http.Handler `name:"sign_up_h"`
	SignIn      http.Handler `name:"sign_in_h"`
	GetSessions http.Handler `name:"get_sessions_h"`
	GetSession  http.Handler `name:"get_session_h"`
	GetUser     http.Handler `name:"get_user_h"`
	Auth        *middleware.UserAuthenticator
}

func NewRouter(p RouterParams) *Router {
	mux := http.NewServeMux()

	// Swagger UI
	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	// API
	mux.Handle("POST /game/create", p.Auth.WithAuth(p.Create))
	mux.Handle("POST /game/{uuid}/connect", p.Auth.WithAuth(p.Connect))
	mux.Handle("POST /game/{uuid}/start", p.Auth.WithAuth(p.Start))
	mux.Handle("POST /game/{uuid}/move", p.Auth.WithAuth(p.Move))
	mux.Handle("POST /game/{uuid}/disconnect", p.Auth.WithAuth(p.Disconnect))

	mux.Handle("POST /auth/signup", p.SignUp)
	mux.Handle("POST /auth/signin", p.SignIn)

	mux.Handle("GET /game", p.Auth.WithAuth(p.GetSessions))
	mux.Handle("GET /game/{uuid}", p.Auth.WithAuth(p.GetSession))
	mux.Handle("GET /user/{uuid}", p.Auth.WithAuth(p.GetUser))

	return &Router{mux: mux}
}

func (r *Router) Handler() http.Handler {
	return r.mux
}
