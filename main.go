package main

import (
	"log"

	"github.com/austindoeswork/music_b/cache"
	"github.com/austindoeswork/music_b/config"
	"github.com/austindoeswork/music_b/server"
)

func main() {
	var err error
	conf, err := config.Find()
	if err != nil {
		log.Fatalf("Fatal: Config error: %s\nExample:\n%s\n", err.Error(), config.Example())
	}

	cache, err := cache.New()
	if err != nil {
		log.Fatalf("Fatal: Cache error: %s\n", err.Error())
	}

	s, err := server.New(conf, cache)
	if err != nil {
		log.Fatalf("Fatal: Server error: %s\n", err.Error())
	}

	err = s.Start()
	if err != nil {
		log.Fatal("Server Error: ", err)
	}
}
