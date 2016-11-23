package router

import (
	"github.com/austindoeswork/music_b/listener"
	"github.com/austindoeswork/music_b/messagehandler"
)

type MessageRouter struct {
	routes map[string]messagehandler.MessageHandler
}

func NewMessageRouter() *MessageRouter {
	return &MessageRouter{
		make(map[string]messagehandler.MessageHandler),
	}
}

func (r *MessageRouter) AddRoute(command string, h messagehandler.MessageHandler) {
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
