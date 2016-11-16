package router

import (
	"github.com/austindoeswork/music_b/handler"
	"github.com/austindoeswork/music_b/listener"
)

//TODO change to message router
//TODO add path router
type MessageRouter struct {
	routes map[string]handler.Handler
}

func NewMessageRouter() *MessageRouter {
	return &MessageRouter{
		make(map[string]handler.Handler),
	}
}

func (r *MessageRouter) AddRoute(command string, h handler.Handler) {
	r.routes[command] = h
}

func (r *MessageRouter) Route(msg listener.Message) {

	cmd := msg.Command()

	if _, ok := r.routes[cmd]; ok {
		go r.routes[cmd].Handle(msg)
	}
}
