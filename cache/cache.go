package cache

import (
	// "fmt"
	"errors"
	"regexp"
	// "strconv"
	"strings"
	"sync"
	"time"
	// "github.com/austindoeswork/music_b/commander"
)

type Cache struct {
	mux     *sync.Mutex
	parties map[string]*Party //name -> party
	threads map[string]string //thread -> party
	songs   map[string]*Song  //songID -> song
	players map[string]string //party -> player
}

func New() *Cache {
	return &Cache{
		&sync.Mutex{},
		make(map[string]*Party),
		make(map[string]string),
		make(map[string]*Song),
		make(map[string]string),
	}
}

//TODO clean this up
//CACHE FUNC=====================================================

func (c *Cache) AddPlayer(partyName string, playerID string) error {
	encodedName := encodeName(partyName)
	// c.mux.Lock()
	if _, ok := c.parties[encodedName]; !ok {
		// c.mux.Unlock()
		return NoPartyError{encodedName}
	}
	c.players[encodedName] = playerID
	// c.mux.Unlock()
	return nil
}

func (c *Cache) DeletePlayer(partyName string) error {
	encodedName := encodeName(partyName)
	// c.mux.Lock()
	if _, ok := c.players[encodedName]; !ok {
		// c.mux.Unlock()
		return NoPartyError{encodedName}
	}
	delete(c.players, encodedName)
	// c.mux.Unlock()
	return nil
}

func (c *Cache) GetPlayer(partyName string) (string, error) {
	encodedName := encodeName(partyName)
	// c.mux.Lock()
	if _, ok := c.players[encodedName]; !ok {
		// c.mux.Unlock()
		return "", NoPartyError{encodedName}
	}
	playerID := c.players[encodedName]
	// c.mux.Unlock()
	return playerID, nil
}

func (c *Cache) GetEncodedName(partyName string) string {
	return encodeName(partyName)
}

func (c *Cache) MakeParty(partyName string) (string, error) {
	if len(partyName) <= 1 {
		return "", errors.New("partyName too short")
	}
	var encodedName string
	encodedName = encodeName(partyName)

	if _, ok := c.parties[encodedName]; ok {
		return "", PartyExistsError{encodedName}
	}

	// c.mux.Lock()
	c.parties[encodedName] = NewParty(partyName, encodedName)
	// c.mux.Unlock()

	return partyName, nil
}

//TODO somehow send msg to connected threads
func (c *Cache) EndParty(partyName string) error {
	encodedName := encodeName(partyName)
	// c.mux.Lock()
	if _, ok := c.parties[encodedName]; !ok {
		// c.mux.Unlock()
		return NoPartyError{encodedName}
	}
	delete(c.parties, encodedName)

	c.DeletePlayer(encodedName)

	toDelete := []string{}
	for thread, party := range c.threads {
		if party == encodedName {
			toDelete = append(toDelete, thread)
		}
	}
	for _, thread := range toDelete {
		delete(c.threads, thread)
	}
	// c.mux.Unlock()
	return nil
}

func (c *Cache) JoinParty(partyName, threadID string) error {
	encodedName := encodeName(partyName)
	// c.mux.Lock()
	if _, ok := c.parties[encodedName]; !ok {
		// c.mux.Unlock()
		return NoPartyError{encodedName}
	}
	c.threads[threadID] = encodedName
	// c.mux.Unlock()
	return nil
}

func (c *Cache) AddSong(songID, path, title string, len time.Duration, requester string) error {
	// c.mux.Lock()
	if song, ok := c.songs[songID]; ok {
		// c.mux.Unlock()
		song.Added()
		return nil
	}

	s := &Song{
		path,
		title,
		len,
		time.Now(),
		requester,
		1,
		0,
		songID,
	}
	c.songs[songID] = s
	// c.mux.Unlock()
	return nil
}

func (c *Cache) GetSong(songID string) (*Song, error) {
	// c.mux.Lock()
	if _, ok := c.songs[songID]; !ok {
		// c.mux.Unlock()
		return nil, NoSongError{songID}
	}
	s := c.songs[songID]
	// c.mux.Unlock()
	return s, nil
}

func (c *Cache) PopSong(partyName string) error {
	encodedName := encodeName(partyName)
	// c.mux.Lock()
	if _, ok := c.parties[encodedName]; !ok {
		// c.mux.Unlock()
		return NoPartyError{encodedName}
	}
	song, _ := c.parties[encodedName].PopSong()
	if _, ok := c.songs[song]; ok {
		c.songs[song].Played()
	}
	// c.mux.Unlock()
	return nil
}

func (c *Cache) DeleteSong(songID string) error {
	if _, ok := c.songs[songID]; !ok {
		return NoSongError{songID}
	}
	delete(c.songs, songID)
	return nil
}

func (c *Cache) DeleteSongFromParty(partyName, songID string) error {
	encodedName := encodeName(partyName)
	// c.mux.Lock()
	if _, ok := c.parties[encodedName]; !ok {
		// c.mux.Unlock()
		return NoPartyError{encodedName}
	}
	c.parties[encodedName].RemoveSongByID(songID)
	// c.mux.Unlock()
	return nil
}

// func (c *Cache) SetPlayer(partyName string, player *commander.Player) error {
// encodedName := encodeName(partyName)
// c.mux.Lock()
// if _, ok := c.parties[encodedName]; !ok {
// c.mux.Unlock()
// return NoPartyError{encodedName}
// }
// c.parties[encodedName].SetPlayer(player)
// c.mux.Unlock()
// return nil
// }

// func (c *Cache) GetPlayer(partyName string) (*commander.Player, error) {
// encodedName := encodeName(partyName)
// c.mux.Lock()
// if _, ok := c.parties[encodedName]; !ok {
// c.mux.Unlock()
// return nil, NoPartyError{encodedName}
// }
// player := c.parties[encodedName].GetPlayer()
// c.mux.Unlock()
// return player, nil
// }

func (c *Cache) AppendSong(partyName, requester, song string) error {
	encodedName := encodeName(partyName)
	// c.mux.Lock()
	if _, ok := c.parties[encodedName]; !ok {
		// c.mux.Unlock()
		return NoPartyError{encodedName}
	}

	err := c.parties[encodedName].AppendSong(song)
	// c.mux.Unlock()

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
	// c.mux.Lock()
	partyID, ok := c.threads[threadID]
	if !ok {
		// c.mux.Unlock()
		return "", NoThreadError{threadID}
	}
	// c.mux.Unlock()
	return partyID, nil
}

func (c *Cache) ThreadToPartyName(threadID string) (string, error) {
	// c.mux.Lock()
	encodedName, ok := c.threads[threadID]
	if !ok {
		// c.mux.Unlock()
		return "", NoThreadError{threadID}
	}
	partyName := c.parties[encodedName].OriginalName()
	// c.mux.Unlock()
	return partyName, nil
}

func (c *Cache) GetParties() []string {
	var pList []string
	// c.mux.Lock()
	for _, p := range c.parties {
		pList = append(pList, p.OriginalName())
	}
	// c.mux.Unlock()
	return pList
}

func (c *Cache) GetSongList(partyName string) ([]string, error) {
	songs := []string{}
	encodedName := encodeName(partyName)
	// c.mux.Lock()
	if _, ok := c.parties[encodedName]; !ok {
		// c.mux.Unlock()
		return nil, NoPartyError{encodedName}
	}

	songList := c.parties[encodedName].GetSongList()
	for _, song := range songList {
		//TODO add checking here
		songs = append(songs, c.songs[song].Title())
	}
	// c.mux.Unlock()

	return songs, nil
}

//TODO better naming on these functions
func (c *Cache) GetAllSongs() ([]*Song, error) {
	slist := []*Song{}
	for _, song := range c.songs {
		slist = append(slist, song)
	}
	return slist, nil
}

func (c *Cache) GetSongs(partyName string, count int) ([]string, error) {
	encodedName := encodeName(partyName)
	// c.mux.Lock()
	if _, ok := c.parties[encodedName]; !ok {
		// c.mux.Unlock()
		return nil, NoPartyError{encodedName}
	}

	songs := c.parties[encodedName].GetSongList()
	if count <= 0 {
		return songs, nil
	}
	if len(songs) >= count {
		songs = songs[:count]
	}
	// c.mux.Unlock()

	return songs, nil
}

func (c *Cache) ClearSongs(partyName string) error {
	encodedName := encodeName(partyName)
	// c.mux.Lock()
	if _, ok := c.parties[encodedName]; !ok {
		// c.mux.Unlock()
		return NoPartyError{encodedName}
	}

	err := c.parties[encodedName].ClearSongs()
	// c.mux.Unlock()

	return err
}

//HELPY==========================================================

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
