package server

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/austindoeswork/music_b/config"
)

type Server struct {
	c *config.Config
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

func (s *Server) MessengerHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Println(formatRequest(req))

	params := req.URL.Query()
	challenge := params.Get("hub.challenge")
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, challenge)
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
	http.HandleFunc("/messenger", s.MessengerHandler)
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir(c.StaticDir))))

	// start server
	log.Printf("blastoff @ %s\n", c.ServerAddress)
	if c.Secure {
		err = http.ListenAndServeTLS(c.ServerAddress, c.SSLCert, c.SSLKey, nil)
	} else {
		err = http.ListenAndServe(c.ServerAddress, nil)
	}

	return err
}
