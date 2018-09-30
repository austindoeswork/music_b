package cache

import (
	"fmt"
	"os"

	"github.com/austindoeswork/music_b/youtube"
)

const (
	AudioDir = "/tmp/audio/"
)

// CACHE
type Cache struct {
	Parties map[string]*Party
}

func New() (*Cache, error) {
	return &Cache{
		Parties: map[string]*Party{},
	}, nil
}

func (c *Cache) GetParty(name string) (*Party, error) {
	if len(name) < 1 {
		return nil, fmt.Errorf("Invalid party name")
	}
	party, exist := c.Parties[name]
	if !exist {
		party = &Party{
			name:   name,
			songs:  []*Song{},
			status: Abandoned,
		}
		c.Parties[name] = party
	}

	return party, nil
}

// PARTY
type PartyStatus int

const (
	Connected PartyStatus = 0
	Abandoned PartyStatus = 1
)

type Party struct {
	name   string
	songs  []*Song
	status PartyStatus
}

// name
func (p *Party) Name() string {
	return p.name
}

// status
func (p *Party) GetStatus() PartyStatus {
	return p.status
}
func (p *Party) UpdateStatus(status PartyStatus) {
	p.status = status
}

// songs
func (p *Party) AddSong(query string) error {
	song, err := NewSong(query)
	if err != nil {
		return err
	}

	p.songs = append(p.songs, song)
	return nil
}
func (p *Party) PopSong() error {
	if len(p.songs) == 0 {
		return fmt.Errorf("No songs to pop")
	}
	p.songs = append(p.songs[:0], p.songs[1:]...)
	return nil
}
func (p *Party) Songs() []*Song {
	return p.songs
}

// SONG
type Song struct {
	ID       string `json:"id"`
	Query    string `json:"query"`
	Filepath string `json:"-"`
}

func NewSong(query string) (*Song, error) {
	hits := youtube.Query(query, 5)
	if _, err := os.Stat(AudioDir); os.IsNotExist(err) {
		os.Mkdir(AudioDir, 0777)
	}

	var err error
	for _, id := range hits {
		filepath := fmt.Sprintf("%s%s.mp3", AudioDir, id)
		err = youtube.Download(id, filepath)
		if err == nil {
			return &Song{
				ID:       id,
				Query:    query,
				Filepath: filepath,
			}, nil
		}
	}
	return nil, err
}
