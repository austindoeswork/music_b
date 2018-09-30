package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/austindoeswork/music_b/cache"
)

func (s *Server) ServeSong(w http.ResponseWriter, r *http.Request) {
	audioDir := cache.AudioDir
	songID := r.URL.Path[len("/song/"):]
	filepath := fmt.Sprintf("%s%s.mp3", audioDir, songID)
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		w.WriteHeader(404)
		w.Write([]byte("song not found"))
		return
	}

	http.ServeFile(w, r, filepath)
}
