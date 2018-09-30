package server

import (
	"fmt"
	"io"
	"net/http"
)

func (s *Server) TestHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("music_b: hello\n"))
	fmt.Fprintf(w, "music_b: hello\n")
	io.WriteString(w, "music_b: hello\n")
}
