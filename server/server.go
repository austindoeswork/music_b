package server

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/austindoeswork/music_b/cache"
	"github.com/austindoeswork/music_b/config"

	"github.com/gorilla/websocket"
)

// SERVER /////////////////////////////////////////////////////////////////////

type Server struct {
	conf  *config.Config
	cache *cache.Cache
}

func New(conf *config.Config, cache *cache.Cache) (*Server, error) {
	s := &Server{
		conf:  conf,
		cache: cache,
	}
	return s, nil
}

func (s *Server) Start() error {
	var err error
	conf := s.conf

	// assign handlers
	http.HandleFunc("/ws", s.WebsocketHandler)
	http.HandleFunc("/song/", s.ServeSong)

	http.HandleFunc("/test", s.TestHandler)
	http.HandleFunc("/messenger", s.MessengerHandler)
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir(conf.StaticDir))))

	// start server
	log.Printf("blastoff @ %s\n", conf.ServerAddress)
	if conf.Secure {
		err = http.ListenAndServeTLS(conf.ServerAddress, conf.SSLCert, conf.SSLKey, nil)
	} else {
		err = http.ListenAndServe(conf.ServerAddress, nil)
	}

	return err
}

// HELPERS ////////////////////////////////////////////////////////////////////

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// formatRequest generates ascii representation of a request
func formatRequest(r *http.Request) string {
	// Create return string
	var request []string
	// Add the request string
	url := fmt.Sprintf("%v %v %v", r.Method, r.URL, r.Proto)
	request = append(request, url)
	// Add the host
	request = append(request, fmt.Sprintf("Host: %v", r.Host))
	// Loop through headers
	for name, headers := range r.Header {
		name = strings.ToLower(name)
		for _, h := range headers {
			request = append(request, fmt.Sprintf("%v: %v", name, h))
		}
	}

	// If this is a POST, add post data
	if r.Method == "POST" {
		r.ParseForm()
		request = append(request, "\n")
		request = append(request, r.Form.Encode())
	}
	// Return the request as a string
	return strings.Join(request, "\n")
}
