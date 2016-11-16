package main

import (
	"flag"
	"fmt"

	"github.com/austindoeswork/music_b/downloader"
)

var (
	queryFlag = flag.String("q", "", "yt query")
)

func main() {
	flag.Parse()
	fmt.Println("hello test")

	yt := downloader.NewYTDownloader("~/Music/Illegal/")
	yt.FromQuery(*queryFlag)
}
