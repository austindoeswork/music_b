package server

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/austindoeswork/music_b/config"
)

type Server struct {
	c *config.Config
}

func (s *Server) TestHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("music_b: hello\n"))
	fmt.Fprintf(w, "music_b: hello\n")
	io.WriteString(w, "music_b: hello\n")
}

func New(c *config.Config) (*Server, error) {
	s := &Server{
		c: c,
	}
	return s, nil
}

func (s *Server) Start() error {
	var err error
	c := s.c

	// assign handlers
	http.HandleFunc("/test", s.TestHandler)

	// start server
	log.Printf("blastoff @ %s\n", c.ServerAddress)
	if c.Secure {
		err = http.ListenAndServeTLS(c.ServerAddress, c.SSLCert, c.SSLKey, nil)
	} else {
		err = http.ListenAndServe(c.ServerAddress, nil)
	}

	return err
}
