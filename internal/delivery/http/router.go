package http

import (
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
)

type Router struct {
	mux *http.ServeMux
}

func NewRouter(create http.Handler, connect http.Handler, start http.Handler, move http.Handler, disconnect http.Handler) *Router {
	mux := http.NewServeMux()

	// Swagger UI
	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	// API
	mux.Handle("POST /game/create", create)
	mux.Handle("POST /game/{uuid}/connect", connect)
	mux.Handle("POST /game/{uuid}/start", start)
	mux.Handle("POST /game/{uuid}/move", move)
	mux.Handle("POST /game/{uuid}/disconnect", disconnect)

	return &Router{mux: mux}
}

func (r *Router) Handler() http.Handler {
	return r.mux
}
