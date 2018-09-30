package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/austindoeswork/music_b/cache"
	"github.com/gorilla/websocket"
)

func (c *Server) WebsocketHandler(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	partyName := params.Get("name")
	party, err := c.cache.GetParty(partyName)
	if err != nil || party.GetStatus() == cache.Connected {
		w.WriteHeader(400)
		w.Write([]byte("Could not join party"))
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Panic("FATAL server: " + err.Error())
	}

	defer func() {
		conn.Close()
		party.UpdateStatus(cache.Abandoned)
	}()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		handleWebsocketMessage(message, conn, party)
	}
}

type WebsocketInput struct {
	Type string `json:"type"`
	Body string `json:"body"`
}

type WebsocketErrorOutput struct {
	Type string `json:"type"`
	Body string `json:"body"`
}

type WebsocketSongsOutput struct {
	Type string        `json:"type"`
	Body []*cache.Song `json:"body"`
}

func handleWebsocketMessage(msg []byte, conn *websocket.Conn, party *cache.Party) {
	var input WebsocketInput
	err := json.Unmarshal(msg, &input)
	if err != nil {
		conn.WriteJSON(WebsocketErrorOutput{
			Type: "error",
			Body: "invalid command",
		})
		return
	}

	switch input.Type {
	case "push":
		query := input.Body
		err := party.AddSong(query)
		if err != nil {
			conn.WriteJSON(WebsocketErrorOutput{
				Type: "error",
				Body: fmt.Sprintf("could not download song: '%s'", query),
			})
			return
		}
	case "pop":
		err := party.PopSong()
		if err != nil {
			conn.WriteJSON(WebsocketErrorOutput{
				Type: "error",
				Body: fmt.Sprintf("could not pop song: '%s'", err.Error()),
			})
			return
		}
	default:
		conn.WriteJSON(WebsocketErrorOutput{
			Type: "error",
			Body: "invalid command",
		})
		return
	}

	conn.WriteJSON(WebsocketSongsOutput{
		Type: "queue",
		Body: party.Songs(),
	})
}
