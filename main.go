package main

import (
	"log"

	"github.com/austindoeswork/music_b/config"
	"github.com/austindoeswork/music_b/server"
)

func main() {
	var err error
	c, err := config.Find()
	if err != nil {
		log.Fatalf("Fatal: Config error: %s\nExample:\n%s\n", err.Error(), config.Example())
	}

	s, err := server.New(c)
	if err != nil {
		log.Fatalf("Fatal: Server error: %s\n", err.Error())
	}

	err = s.Start()
	if err != nil {
		log.Fatal("Server Error: ", err)
	}
}
