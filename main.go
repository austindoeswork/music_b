package main

import (
	"fmt"
	"os"

	"github.com/austindoeswork/music_b/cache"
	"github.com/austindoeswork/music_b/handler"
	"github.com/austindoeswork/music_b/listener"
	"github.com/austindoeswork/music_b/router"
)

var (
	fbsessionPath = os.Getenv("HOME") + "/.music_bitch_ui/session_data"
	fbemail       = os.Getenv("FBEMAIL")
	fbpass        = os.Getenv("FBPASS")
)

func main() {

	//facebook
	f, err := listener.NewFBListener(fbemail, fbpass, fbsessionPath)
	if err != nil {
		fmt.Println(err.Error())
	}
	fb_c := f.Listen()

	//cache
	c := cache.New()

	//make test party
	pID, err := c.MakeParty("")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("Party Created: ", pID)

	//router
	r := router.New()
	addRoutes(r, c)

	//commander
	//TODO

	//pass messages to router
	for {
		select {
		case msg := <-fb_c:
			fmt.Println(msg.Fulltext())
			r.Route(msg)
		}
	}

}

func addRoutes(r *router.Router, c *cache.Cache) {
	r.AddRoute(".test", handler.NewTestHandler())
	r.AddRoute(".join", handler.NewJoinPartyHandler(c))
	r.AddRoute(".parties", handler.NewGetPartiesHandler(c))
	r.AddRoute(".status", handler.NewStatusHandler(c))
	r.AddRoute(".play", handler.NewAddSongHandler(c))
}

func easterEggs(msg *listener.Message) {

}
