package listener

import (
	"github.com/austindoeswork/music_b/listener/facebook"
)

type Message interface {
	UserID() string
	UserName() string
	ThreadID() string
	IsGroup() bool
	Command() string
	HasFlag(potFlag string) bool
	Text() string
	Fulltext() string
	Respond(text string) error
	//TODO add attachments?
}

type Listener interface {
	Listen() (chan Message, error)
	//TODO ListenAndRoute(router) (error)
	Send(threadID, text string) error
}

//implementations
func NewFBListener(username, password, sessionPath string) (*facebook.Facebook, error) {
	return facebook.New(username, password, sessionPath)
}
