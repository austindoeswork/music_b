package cache

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Cache struct {
	mux     *sync.Mutex
	parties map[string]*Party
	threads map[string]string
}

func New() *Cache {
	return &Cache{
		&sync.Mutex{},
		make(map[string]*Party),
		make(map[string]string),
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
		fmt.Println(ok)
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

func (c *Cache) AppendSong(partyName, requester, song string) error {
	encodedName := encodeName(partyName)
	c.mux.Lock()
	if _, ok := c.parties[encodedName]; !ok {
		c.mux.Unlock()
		return NoPartyError{encodedName}
	}

	err := c.parties[encodedName].AppendSong(song, requester)
	c.mux.Unlock()

	return err
}

func (c *Cache) PrependSong(partyName, requester, song string) error {
	encodedName := encodeName(partyName)
	c.mux.Lock()
	if _, ok := c.parties[encodedName]; !ok {
		c.mux.Unlock()
		return NoPartyError{encodedName}
	}

	err := c.parties[encodedName].PrependSong(song, requester)
	c.mux.Unlock()

	return err
}

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
	encodedName := encodeName(partyName)
	c.mux.Lock()
	if _, ok := c.parties[encodedName]; !ok {
		c.mux.Unlock()
		return nil, NoPartyError{encodedName}
	}

	songList := c.parties[encodedName].GetSongList()
	c.mux.Unlock()

	return songList, nil
}

//HELPY==========================================================

func genID() string {
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

func encodeName(name string) string {
	encodedName := strings.ToLower(name)
	re := regexp.MustCompile("\\W")
	encodedName = re.ReplaceAllString(encodedName, "")
	return encodedName
}
