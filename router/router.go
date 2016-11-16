package router

import (
	"github.com/austindoeswork/music_b/handler"
	"github.com/austindoeswork/music_b/listener"
)

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

//TODO make seperate msg and cmd handlers
// type CommandRouter struct {
// routes map[string]handler.CommandHandler
// }

// func NewCommandRouter() *CommandRouter {
// return &CommandRouter{
// make(map[string]handler.CommandHandler),
// }
// }
