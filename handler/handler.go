package handler

import (
	"fmt"
	"strings"

	"github.com/austindoeswork/music_b/cache"
	"github.com/austindoeswork/music_b/downloader"
	"github.com/austindoeswork/music_b/listener"
)

type Handler interface {
	Handle(msg listener.Message)
}

func NewTestHandler() *TestHandler {
	return &TestHandler{}
}
func NewHelpHandler() *HelpHandler {
	return &HelpHandler{}
}
func NewJoinPartyHandler(c *cache.Cache) *JoinPartyHandler {
	return &JoinPartyHandler{c}
}
func NewGetPartiesHandler(c *cache.Cache) *GetPartiesHandler {
	return &GetPartiesHandler{c}
}
func NewAddSongHandler(c *cache.Cache, d *downloader.YTDownloader) *AddSongHandler {
	return &AddSongHandler{c, d}
}
func NewStatusHandler(c *cache.Cache) *StatusHandler {
	return &StatusHandler{c}
}

//IMPLEMENTATIONS =============================================================

//TESTHANDLER ====================
type TestHandler struct {
}

func (h *TestHandler) Handle(msg listener.Message) {
	msg.Respond(fmt.Sprintf("usr: %s\ncmd: %s\nrest: %s\n", msg.UserName(), msg.Command(), msg.Text()))
}

//HELPHANDLER ====================
type HelpHandler struct {
}

func (h *HelpHandler) Handle(msg listener.Message) {
	msg.Respond(commandList)
}

//JOINPARTY ====================
type JoinPartyHandler struct {
	c *cache.Cache
}

func (h *JoinPartyHandler) Handle(msg listener.Message) {
	err := h.c.JoinParty(msg.Text(), msg.ThreadID())
	if err != nil {
		msg.Respond("error joining party, sorry :(")
		fmt.Println(err.Error())
		return
	}
	msg.Respond("Joined Party: " + msg.Text())
	return
}

//ADDSONG ======================
type AddSongHandler struct {
	c *cache.Cache
	d *downloader.YTDownloader
}

func (h *AddSongHandler) Handle(msg listener.Message) {
	partyID, err := h.c.ThreadToPartyID(msg.ThreadID())
	if err != nil {
		msg.Respond("please join a party bb <3")
		return
	}
	if len(msg.Text()) < 5 {
		msg.Respond("could u plz gimme a whole song name or sumtin")
		return
	}
	// if msg.HasFlag("-n") {
	// err := h.c.PrependSong(partyID, msg.UserName(), msg.Text())
	// if err != nil {
	// fmt.Println("bf309eda-51a6-4614-81b2-48e643ac9f9d")
	// msg.Respond("something went wrong bb :(")
	// return
	// }
	// msg.Respond("prepended.")
	// return
	// }
	go func() {
		vid, filePath, title, duration, err := h.d.FromQuery(msg.Text())
		if err != nil {
			msg.Respond("error downloading: " + err.Error())
		} else {
			h.c.AddSong(vid, filePath, title, duration, msg.UserName())
			msg.Respond(title + " succesfully stolen.")
			err := h.c.AppendSong(partyID, msg.UserName(), vid)
			if err != nil {
				fmt.Println("f5575f0b-5af0-4a1d-9feb-4f9266fbb3b0")
				msg.Respond("something went wrong adding song to queue :(")
				return
			}

		}
	}()
	msg.Respond("downloading...")
}

//GETPARTIES ====================
type GetPartiesHandler struct {
	c *cache.Cache
}

func (h *GetPartiesHandler) Handle(msg listener.Message) {
	response := "Parties Happenin: \n"

	pList := h.c.GetParties()
	response += strings.Join(pList, "\n")

	msg.Respond(response)
	return
}

//STATUS ====================
type StatusHandler struct {
	c *cache.Cache
}

func (h *StatusHandler) Handle(msg listener.Message) {
	partyName, err := h.c.ThreadToPartyName(msg.ThreadID())
	if err != nil {
		msg.Respond("please join a party bb <3")
		return
	}

	songs, err := h.c.GetSongList(partyName)
	if err != nil {
		fmt.Println("4489eeed-f4d5-4ed5-9e9a-3a0d4897cdf1")
	}

	response := ""

	if msg.HasFlag("-q") {
		if len(songs) == 0 {
			response += "no songs in queue"
		} else {
			response += partyName + ": \n" + strings.Join(songs, "\n")
		}
	} else if msg.HasFlag("-j") {
		songs_j, err := h.c.GetSongsJson(partyName)
		if err != nil {
			fmt.Println("1a351459-5533-4e8f-99f1-4514c8022739")
		}
		response += string(songs_j)
	} else {
		if len(songs) == 0 {
			response += "nothing's playing dawg"
		} else {
			response += "Now playing: \n" + songs[0]
		}
	}

	msg.Respond(response)
	return
}

const (
	commandList = `.help: display this list
.play <song query>:	add a song to queue
.parties: list all parties
.join <party name>: start controlling a party
.status: what's happening?
	-q queue
	-j json for some reason
//.clear: fuck people in the ass
//.pause: be that guy who killed the jams
//.resume: resume paused music
//.skip: fuck this song
//.whoami: who the fuck am I
`
)
