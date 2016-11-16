package cache

import (
	"encoding/json"
	"strings"

	"github.com/austindoeswork/music_b/commander"
	// "github.com/wardn/uuid"
)

type Party struct {
	originalName string
	player       *commander.Player
	songs        []string
}

func NewParty(name string) *Party {
	return &Party{
		name,
		nil,
		[]string{},
	}
}

func (p *Party) SetPlayer(player *commander.Player) {
	//TODO handle nil
	//TODO how to handle hijaking?
	p.player = player
}

func (p *Party) GetPlayer() *commander.Player {
	return p.player
}

func (p *Party) RemoveSongByID(songID string) error {
	//TODO make this wayyyy more efficient
	for i, song := range p.songs {
		if song == songID {
			//delete the song
			p.songs = append(p.songs[:i], p.songs[i+1:]...)
			return nil
		}
	}
	return NoSongError{songID}
}

func (p *Party) ClearSongs() error {
	p.songs = []string{}
	return nil
}

func (p *Party) AppendSong(sid string) error {
	p.songs = append(p.songs, sid)
	return nil
}

// func (p *Party) PrependSong(songName, requester string) error {
// p.songs = append([]*song{newSong(songName, requester)}, p.songs...)
// return nil
// }

func (p *Party) GetSongList() []string {
	return p.songs
}

func (p *Party) GetSongsJson() ([]byte, error) {
	return json.Marshal(p.songs)
}

func (p *Party) OriginalName() string {
	return p.originalName
}

//helpers

func getInitials(fullName string) string {
	names := strings.Split(fullName, " ")
	if len(names) <= 0 {
		return "??"
	}

	var first string
	if len(names[0]) <= 0 {
		first = "?"
	} else {
		first = string(names[0][0])
	}
	var last string
	if len(names[len(names)-1]) <= 0 {
		last = "?"
	} else {
		last = string(names[len(names)-1][0])
	}
	return first + last
}
