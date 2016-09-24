package cache

import (
	"strings"

	"github.com/wardn/uuid"
)

type Party struct {
	originalName string
	songs        []*song
}

func NewParty(name string) *Party {
	return &Party{
		name,
		[]*song{},
	}
}

func (p *Party) AppendSong(songName, requester string) error {
	p.songs = append(p.songs, newSong(songName, requester))
	return nil
}

func (p *Party) PrependSong(songName, requester string) error {
	p.songs = append([]*song{newSong(songName, requester)}, p.songs...)
	return nil
}

func (p *Party) GetSongList() []string {
	var songList []string
	for _, song := range p.songs {
		songList = append(songList, getInitials(song.Requester)+" - "+song.Name)
	}
	return songList
}

func (p *Party) OriginalName() string {
	return p.originalName
}

//TODO make immutable >:(
type song struct {
	ID        string `json: "BitchID"`
	Requester string `json: "Requester"`
	Name      string `json: "Name"`
}

func newSong(name, requester string) *song {
	return &song{
		uuid.New(),
		requester,
		name,
	}
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
