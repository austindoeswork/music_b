package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/austindoeswork/music_b/cache"
	"github.com/austindoeswork/music_b/commander"
	"github.com/austindoeswork/music_b/config"
	"github.com/austindoeswork/music_b/downloader"
	"github.com/austindoeswork/music_b/handler"
	"github.com/austindoeswork/music_b/listener"
	"github.com/austindoeswork/music_b/router"
	"github.com/austindoeswork/music_b/server"
)

var (
	fbsessionPath = os.Getenv("HOME") + "/.music_b/session_data"
	fbemail       = os.Getenv("FBEMAIL")
	fbpass        = os.Getenv("FBPASS")

	configFlag  = flag.String("config", os.Getenv("HOME")+"/.music_b/default.json", "path to config")
	versionFlag = flag.Bool("v", false, "git commit hash")

	commithash string
)

func main() {
	flag.Parse()
	if *versionFlag {
		fmt.Println(commithash)
		return
	}
	fmt.Println("version: " + commithash)
	//config
	conf, err := config.Parse(*configFlag)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	conf.PPrint()

	//facebook
	f, err := listener.NewFBListener(fbemail, fbpass, fbsessionPath)
	if err != nil {
		fmt.Println(err.Error())
	}
	fb_c := f.Listen()
	fmt.Println("LISTENER: initialized.")

	//cache
	c := cache.New()

	//commander
	com := commander.New(c)
	com.Listen(conf.StaticDir)

	//downloader
	ytd, err := downloader.NewYTDownloader(c, conf.MusicDir)
	if err != nil {
		fmt.Println("DOWNLOADER: failed.")
		return
	}
	fmt.Println("DOWNLOADER: initialized.")

	//server
	serv := server.New(c)
	serv.Start()

	//router
	r := router.NewMessageRouter()
	addMessageRoutes(r, c, ytd, com)

	//pass messages to router
	go func() {
		for {
			select {
			case msg := <-fb_c:
				fmt.Println("fb: " + msg.Fulltext())
				go easterEggs(msg)
				r.Route(msg)
			}
		}
	}()

	//start
	fmt.Println("blastoff. " + conf.ServerPath)
	http.ListenAndServe(conf.ServerPath, nil)
}

func addMessageRoutes(r *router.MessageRouter, c *cache.Cache, d *downloader.YTDownloader, com *commander.Commander) {
	r.AddRoute(".test", handler.NewTestHandler())
	r.AddRoute(".help", handler.NewHelpHandler())
	r.AddRoute(".clear", handler.NewClearHandler(c))
	r.AddRoute(".join", handler.NewJoinPartyHandler(c))
	r.AddRoute(".parties", handler.NewGetPartiesHandler(c))
	r.AddRoute(".status", handler.NewStatusHandler(c))
	r.AddRoute(".skip", handler.NewSkipHandler(c, com))
	r.AddRoute(".add", handler.NewAddSongHandler(c, d))
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
