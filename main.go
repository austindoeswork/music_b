package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/austindoeswork/music_b/cache"
	"github.com/austindoeswork/music_b/commander"
	"github.com/austindoeswork/music_b/downloader"
	"github.com/austindoeswork/music_b/handler"
	"github.com/austindoeswork/music_b/listener"
	"github.com/austindoeswork/music_b/router"
	"github.com/austindoeswork/music_b/server"
)

var (
	fbsessionPath = os.Getenv("HOME") + "/.music_bitch_ui/session_data"
	fbemail       = os.Getenv("FBEMAIL")
	fbpass        = os.Getenv("FBPASS")

	audioDirFlag = flag.String("audioDir", "/Users/austin/Music/Illegal/", "path to audioDir")
	//TODO use a config
	comFlag = flag.String("com", ":80", "port for com")
	serFlag = flag.String("server", ":6969", "port for music server")
)

func main() {
	flag.Parse()

	//facebook
	f, err := listener.NewFBListener(fbemail, fbpass, fbsessionPath)
	if err != nil {
		fmt.Println(err.Error())
	}
	fb_c := f.Listen()

	//cache
	c := cache.New()

	//make test party
	pID, err := c.MakeParty("Sausage Fest")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("Party Created: ", pID)

	//commander
	com := commander.New(c)
	go com.Listen(*comFlag)

	//downloader
	ytd := downloader.NewYTDownloader(*audioDirFlag)

	//server
	serv := server.New(c)
	go serv.Start(*serFlag)

	//router
	r := router.NewMessageRouter()
	addMessageRoutes(r, c, ytd, com)

	//pass messages to router
	for {
		select {
		case msg := <-fb_c:
			fmt.Println("fb: " + msg.Fulltext())
			go easterEggs(msg)
			r.Route(msg)
		}
	}

}

func addMessageRoutes(r *router.MessageRouter, c *cache.Cache, d *downloader.YTDownloader, com *commander.Commander) {
	r.AddRoute(".test", handler.NewTestHandler())
	r.AddRoute(".help", handler.NewHelpHandler())
	r.AddRoute(".clear", handler.NewClearHandler(c))
	r.AddRoute(".join", handler.NewJoinPartyHandler(c))
	r.AddRoute(".parties", handler.NewGetPartiesHandler(c))
	r.AddRoute(".status", handler.NewStatusHandler(c))
	r.AddRoute(".skip", handler.NewSkipHandler(c, com))
	r.AddRoute(".play", handler.NewAddSongHandler(c, d))
}

func easterEggs(msg listener.Message) {
	words := strings.Split(msg.Fulltext(), " ")

	//check for fucks and your
	for i := 0; i < len(words)-1; i++ {
		if words[i] == "fuck" {
			msg.Respond("yo fuck " + words[i+1] + " tho")
		}
		if words[i] == "your" {
			msg.Respond("*you're")
		}
	}
}
