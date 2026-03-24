package http

import (
	"net/http"
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

	mux.HandleFunc("POST /game/create", r.create.ServeHTTP)
	mux.HandleFunc("POST /game/{uuid}", r.connect.ServeHTTP)
	mux.HandleFunc("POST /game/{uuid}/start", r.start.ServeHTTP)
	mux.HandleFunc("POST /game/{uuid}/move", r.move.ServeHTTP)
	mux.HandleFunc("POST /game/{uuid}", r.disconnect.ServeHTTP)

	return mux
}
