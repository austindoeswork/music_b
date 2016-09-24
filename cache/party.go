package cache

import "github.com/wardn/uuid"

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
		songList = append(songList, song.Name)
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
