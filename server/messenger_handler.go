package server

import (
	"fmt"
	"net/http"
)

func (s *Server) MessengerHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Println(formatRequest(req))

	params := req.URL.Query()
	challenge := params.Get("hub.challenge")
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, challenge)
}
