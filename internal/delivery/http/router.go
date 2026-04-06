package http

import (
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
)

type Router struct {
	create     http.Handler
	connect    http.Handler
	start      http.Handler
	move       http.Handler
	disconnect http.Handler
}

func NewRouter(create http.Handler, connect http.Handler, start http.Handler, move http.Handler, disconnect http.Handler) *Router {
	return &Router{
		create:     create,
		connect:    connect,
		start:      start,
		move:       move,
		disconnect: disconnect,
	}
}

func (r *Router) Handler() http.Handler {
	mux := http.NewServeMux()

	// Swagger UI
	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	// API
	mux.Handle("POST /game/create", r.create)
	mux.Handle("POST /game/{uuid}/connect", r.connect)
	mux.Handle("POST /game/{uuid}/start", r.start)
	mux.Handle("POST /game/{uuid}/move", r.move)
	mux.Handle("POST /game/{uuid}/disconnect", r.disconnect)

	return mux
}
