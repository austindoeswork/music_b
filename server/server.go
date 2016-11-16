package server

import (
	"fmt"
	"net/http"
	// "os"

	"github.com/austindoeswork/music_b/cache"
)

type Server struct {
	c *cache.Cache
}

func New(c *cache.Cache) *Server {
	return &Server{
		c,
	}
}

func (s *Server) serveSong(w http.ResponseWriter, r *http.Request) {
	songID := r.URL.Path[len("/song/"):]
	song, err := s.c.GetSong(songID)
	if err != nil {
		w.Write([]byte("no such song"))
		return
	}
	fmt.Println(song.Path())

	http.ServeFile(w, r, song.Path())
}

func (s *Server) Start() {
	f := s.serveSong
	http.HandleFunc("/song/", f)
	http.ListenAndServe(":8080", nil)
}
