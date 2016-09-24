package facebook

// The email and password used is taken from the environment variables
// FBEMAIL and FBPASS.

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/1lann/messenger"
)

type FBMessage struct {
	userID   string
	userName string
	threadID string
	command  string
	flags    map[string]bool
	text     string
	fulltext string
	isGroup  bool
	fb       *Facebook
}

func NewMessage(userID, userName, threadID, fulltext string, isGroup bool, fb *Facebook) *FBMessage {
	cmd, flgs, rest := parse(fulltext)
	return &FBMessage{
		userID,
		userName,
		threadID,
		cmd,
		flgs,
		rest,
		fulltext,
		isGroup,
		fb,
	}
}

func parse(fulltext string) (string, map[string]bool, string) {
	if len(fulltext) == 0 {
		return "", nil, ""
	}

	var command string
	flags := make(map[string]bool)
	var remainingText string

	words := strings.Split(fulltext, " ")

	if words[0][0] == '.' {
		if len(words[0]) == 1 && len(words) >= 2 {
			command = strings.Join(words[:2], "")
			words = words[2:]
		} else {
			command = words[0]
			words = words[1:]
		}
	} else {
		return "", nil, fulltext
	}

	words_i := 0
	for ; words_i < len(words); words_i++ {
		potFlag := words[words_i]
		if len(potFlag) >= 1 && potFlag[0] == '-' {
			flags[potFlag] = true
		} else {
			break
		}
	}
	remainingText = strings.Join(words[words_i:], " ")

	return strings.ToLower(command), flags, remainingText

}

func (m *FBMessage) Respond(text string) error {
	return m.fb.Send(m.ThreadID(), m.IsGroup(), text)
}
func (m *FBMessage) Command() string {
	return m.command
}
func (m *FBMessage) HasFlag(potFlag string) bool {
	_, ok := m.flags[potFlag]
	return ok
}
func (m *FBMessage) Text() string {
	return m.text
}
func (m *FBMessage) Fulltext() string {
	return m.fulltext
}
func (m *FBMessage) UserID() string {
	return m.userID
}
func (m *FBMessage) UserName() string {
	return m.userName
}
func (m *FBMessage) ThreadID() string {
	return m.threadID
}
func (m *FBMessage) IsGroup() bool {
	return m.isGroup
}

type Facebook struct {
	sesh        *messenger.Session
	sessionPath string
}

func New(username, password, sessionPath string) (*Facebook, error) {

	sesh := messenger.NewSession()
	sessionData, err := ioutil.ReadFile(sessionPath)

	if os.IsNotExist(err) {
		fmt.Println("FB: No session file, logging in with user and passwd")
		err = sesh.Login(username, password)
		if err != nil {
			return nil, err
		}
	} else {
		err := sesh.RestoreSession(sessionData)
		if err != nil {
			fmt.Println("FB: Failed to restore session, logging in with user and passwd")
			err = sesh.Login(username, password)
			if err != nil {
				return nil, err
			}
		}
	}

	err = sesh.ConnectToChat()
	if err != nil {
		return nil, err
	}

	fb := &Facebook{
		sesh,
		sessionPath,
	}

	//preserve the sesh yo
	go func() {
		ticker := time.Tick(time.Minute)
		for range ticker {
			fb.saveSession()
		}
	}()

	return fb, nil

}

func (f *Facebook) Listen() <-chan *FBMessage {
	c := make(chan *FBMessage)
	f.sesh.OnMessage(func(msg *messenger.Message) {
		var userName string
		userProfile, err := f.sesh.UserProfileInfo(msg.FromUserID)
		if err != nil {
			userName = ""
		} else {
			userName = userProfile.Name
		}
		m := NewMessage(msg.FromUserID, userName, msg.Thread.ThreadID, msg.Body, msg.Thread.IsGroup, f)
		c <- m
	})
	go f.sesh.Listen()
	return c
}

func (f *Facebook) Send(threadID string, isGroup bool, text string) error {

	t := messenger.Thread{
		threadID,
		isGroup,
	}

	m := f.sesh.NewMessageWithThread(t)
	m.Body = text

	f.sesh.SendMessage(m)

	return nil
}

func (f *Facebook) saveSession() {
	data, err := f.sesh.DumpSession()
	if err != nil {
		fmt.Println("FB: Failed to save session:", err.Error())
		return
	}

	//TODO smarter opening files
	err = ioutil.WriteFile(f.sessionPath, data, 0644)
	if err != nil {
		fmt.Println("FB: Failed to write session to file:" + err.Error())
	}
}
