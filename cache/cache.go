package cache

import (
	// "fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/austindoeswork/music_b/commander"
)

type Cache struct {
	mux     *sync.Mutex
	parties map[string]*Party
	threads map[string]string
	songs   map[string]*Song
}

func New() *Cache {
	return &Cache{
		&sync.Mutex{},
		make(map[string]*Party),
		make(map[string]string),
		make(map[string]*Song),
	}
}

//CACHE FUNC=====================================================
func (c *Cache) MakeParty(partyName string) (string, error) {
	var encodedName string

	if len(partyName) == 0 {
		//TODO make sure there's no collisions
		partyName = genID()
		encodedName = partyName
	} else {
		encodedName = encodeName(partyName)
	}

	if _, ok := c.parties[encodedName]; ok {
		return "", PartyExistsError{encodedName}
	}

	c.mux.Lock()
	c.parties[encodedName] = NewParty(partyName)
	c.mux.Unlock()

	return partyName, nil
}

func (c *Cache) JoinParty(partyName, threadID string) error {
	encodedName := encodeName(partyName)
	c.mux.Lock()
	if _, ok := c.parties[encodedName]; !ok {
		c.mux.Unlock()
		return NoPartyError{encodedName}
	}
	c.threads[threadID] = encodedName
	c.mux.Unlock()
	return nil
}

func (c *Cache) AddSong(songID, path, title string, len time.Duration, requester string) error {
	c.mux.Lock()
	s := &Song{
		path,
		title,
		len,
		time.Now(),
		requester,
	}
	c.songs[songID] = s
	c.mux.Unlock()
	return nil
}

func (c *Cache) GetSong(songID string) (*Song, error) {
	c.mux.Lock()
	if _, ok := c.songs[songID]; !ok {
		c.mux.Unlock()
		return nil, NoSongError{songID}
	}
	s := c.songs[songID]
	c.mux.Unlock()
	return s, nil
}

func (c *Cache) DeleteSong(partyName, songID string) error {
	encodedName := encodeName(partyName)
	c.mux.Lock()
	if _, ok := c.parties[encodedName]; !ok {
		c.mux.Unlock()
		return NoPartyError{encodedName}
	}
	c.parties[encodedName].RemoveSongByID(songID)
	c.mux.Unlock()
	return nil
}

func (c *Cache) SetPlayer(partyName string, player *commander.Player) error {
	encodedName := encodeName(partyName)
	c.mux.Lock()
	if _, ok := c.parties[encodedName]; !ok {
		c.mux.Unlock()
		return NoPartyError{encodedName}
	}
	c.parties[encodedName].SetPlayer(player)
	c.mux.Unlock()
	return nil
}

func (c *Cache) GetPlayer(partyName string) (*commander.Player, error) {
	encodedName := encodeName(partyName)
	c.mux.Lock()
	if _, ok := c.parties[encodedName]; !ok {
		c.mux.Unlock()
		return nil, NoPartyError{encodedName}
	}
	player := c.parties[encodedName].GetPlayer()
	c.mux.Unlock()
	return player, nil
}

func (c *Cache) AppendSong(partyName, requester, song string) error {
	encodedName := encodeName(partyName)
	c.mux.Lock()
	if _, ok := c.parties[encodedName]; !ok {
		c.mux.Unlock()
		return NoPartyError{encodedName}
	}

	err := c.parties[encodedName].AppendSong(song)
	c.mux.Unlock()

	return err
}

// func (c *Cache) PrependSong(partyName, requester, song string) error {
// encodedName := encodeName(partyName)
// c.mux.Lock()
// if _, ok := c.parties[encodedName]; !ok {
// c.mux.Unlock()
// return NoPartyError{encodedName}
// }

// err := c.parties[encodedName].PrependSong(song, requester)
// c.mux.Unlock()

// return err
// }

func (c *Cache) ThreadToPartyID(threadID string) (string, error) {
	c.mux.Lock()
	partyID, ok := c.threads[threadID]
	if !ok {
		c.mux.Unlock()
		return "", NoThreadError{threadID}
	}
	c.mux.Unlock()
	return partyID, nil
}

func (c *Cache) ThreadToPartyName(threadID string) (string, error) {
	c.mux.Lock()
	encodedName, ok := c.threads[threadID]
	if !ok {
		c.mux.Unlock()
		return "", NoThreadError{threadID}
	}
	partyName := c.parties[encodedName].OriginalName()
	c.mux.Unlock()
	return partyName, nil
}

func (c *Cache) GetParties() []string {
	var pList []string
	c.mux.Lock()
	for _, p := range c.parties {
		pList = append(pList, p.OriginalName())
	}
	c.mux.Unlock()
	return pList
}

func (c *Cache) GetSongList(partyName string) ([]string, error) {
	songs := []string{}
	encodedName := encodeName(partyName)
	c.mux.Lock()
	if _, ok := c.parties[encodedName]; !ok {
		c.mux.Unlock()
		return nil, NoPartyError{encodedName}
	}

	songList := c.parties[encodedName].GetSongList()
	for _, song := range songList {
		//TODO add checking here
		songs = append(songs, c.songs[song].Title())
	}
	c.mux.Unlock()

	return songs, nil
}

func (c *Cache) GetSongsJson(partyName string) ([]byte, error) {
	encodedName := encodeName(partyName)
	c.mux.Lock()
	if _, ok := c.parties[encodedName]; !ok {
		c.mux.Unlock()
		return nil, NoPartyError{encodedName}
	}

	songs_j, err := c.parties[encodedName].GetSongsJson()
	c.mux.Unlock()

	return songs_j, err
}

//HELPY==========================================================

func genID() string {
	//TODO seed once yo
	rand.Seed(int64(time.Now().Nanosecond()))
	return strconv.Itoa(rand.Intn(8999) + 1000)
}

type NoPartyError struct {
	partyID string
}

func (n NoPartyError) Error() string {
	return "No such party: " + n.partyID
}

type PartyExistsError struct {
	partyID string
}

func (n PartyExistsError) Error() string {
	return "Party Already Exists: " + n.partyID
}

type NoThreadError struct {
	threadID string
}

func (n NoThreadError) Error() string {
	return "No such thread: " + n.threadID
}

type NoSongError struct {
	songID string
}

func (n NoSongError) Error() string {
	return "No such song: " + n.songID
}

func encodeName(name string) string {
	encodedName := strings.ToLower(name)
	re := regexp.MustCompile("\\W")
	encodedName = re.ReplaceAllString(encodedName, "")
	return encodedName
}
