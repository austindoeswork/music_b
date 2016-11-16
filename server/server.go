package server

import (
	"fmt"
	"net/http"

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
	fmt.Println("SERVER: serving " + song.Path())
	http.ServeFile(w, r, song.Path())
}

func (s *Server) Start(port string) {
	f := s.serveSong
	http.HandleFunc("/song/", f)
	fmt.Println("SERVER: started @ " + port)
	http.ListenAndServe(port, nil)
}
