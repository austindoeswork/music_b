package router

import (
	"github.com/austindoeswork/music_b/handler"
	"github.com/austindoeswork/music_b/listener"
)

type Router struct {
	routes map[string]handler.Handler
}

func New() *Router {
	return &Router{
		make(map[string]handler.Handler),
	}
}

func (r *Router) AddRoute(command string, h handler.Handler) {
	r.routes[command] = h
}

func (r *Router) Route(msg listener.Message) {

	cmd := msg.Command()

	if _, ok := r.routes[cmd]; ok {
		go r.routes[cmd].Handle(msg)
	}
}
