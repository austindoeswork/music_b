package downloader

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/austindoeswork/music_b/cache"

	"github.com/rylio/ytdl"
	"golang.org/x/net/html"
)

var (
	ytSearchUrl = "https://www.youtube.com/results?search_query="
	ytWatchUrl  = "https://www.youtube.com/watch?v="
)

//TODO make this clean itself up somehow

type YTDownloader struct {
	c   *cache.Cache
	dir string
}

func NewYTDownloader(c *cache.Cache, dir string) (*YTDownloader, error) {
	d, err := os.Open(dir)
	if err != nil {
		return nil, err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return nil, err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return nil, err
		}
	}

	downloader := YTDownloader{
		c,
		dir,
	}

	go func() {
		for x := range time.Tick(time.Minute * 3) {
			fmt.Println("DOWNLOADER: clean up starting... " + x.String())
			err := downloader.Clean()
			if err != nil {
				log.Println(err.Error())
			}
		}
	}()

	return &downloader, nil
}

func (d *YTDownloader) Clean() error {
	songs, err := d.c.GetAllSongs()
	if err != nil {
		return err
	}
	for _, song := range songs {
		//TODO also delete old songs
		if song.AddCount() == song.PlayCount() {
			err = d.c.DeleteSong(song.ID())
			if err != nil {
				fmt.Println("DOWNLOADER: couldn't delete " + song.ID() + " from cache")
				return err
			}
			err = os.Remove(song.Path())
			if err != nil {
				fmt.Println("DOWNLOADER: couldn't delete " + song.ID() + " from filesystem")
				return err
			}
			fmt.Println("DOWNLOADER: deleted " + song.Path())
		}
	}
	return nil
}

//TODO err handling
//TODO threading and pooling
func (d *YTDownloader) FromQuery(query string) (string, string, string, time.Duration, error) {
	vids := ytQuery(query, 1)

	// for _, vid := range vids {
	// fmt.Println(vid)
	// }

	if vids == nil {
		return "", "", "", 0, errors.New("failed to get video")
	}
	vid := vids[0]
	if s, err := d.c.GetSong(vid); err == nil {
		fmt.Println("song exists")
		return vid, s.Path(), s.Title(), s.Length(), nil
	}

	v, err := ytdl.GetVideoInfo(ytWatchUrl + vid)
	if err != nil {
		return "", "", "", 0, err
	}
	// fmt.Println(v.Title)

	filePath := d.dir + vid + ".musicb"
	cmd := exec.Command("youtube-dl", "--extract-audio",
		"-o", filePath, ytWatchUrl+vid)
	err = cmd.Start()
	if err != nil {
		return "", "", "", 0, err
	}

	// log.Printf("Waiting for command to finish...")
	errCode := cmd.Wait()

	log.Printf("Command finished with error: %v", errCode)

	return vid, filePath, v.Title, v.Duration, nil
}

// HELPERS
//ytQuery returns a list of ids from a search string
//ported from Tylers python script
func ytQuery(query string, maxHits int) []string {
	query = strings.Replace(query, " ", "+", -1)
	res, err := http.Get(ytSearchUrl + query)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	z := html.NewTokenizer(res.Body)
	var vids []string
	count := 0
	//TODO scrape better
	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			// End of the document, we're done
			return vids
		case tt == html.StartTagToken:
			t := z.Token()

			if t.Data == "div" {
				found := false
				for _, a := range t.Attr {
					if a.Key == "class" {
						if a.Val == "yt-lockup yt-lockup-tile yt-lockup-video vve-check clearfix" {
							found = true
						}
					}
					if found && a.Key == "data-context-item-id" {
						vids = append(vids, a.Val)
						count++
						if count == maxHits {
							return vids
						}
					} //Wow
				} //This
			} //Is
		} //Hella
	} //Dirty

	return nil

}

// func getBetween(input, from, til string) string {
// for _, line := range strings.Split(input, "\n") {
// if strings.Contains(line, from) {
// continue
// } else {
// output := strings.SplitN(line, from, 1)[0]
// println(strings.SplitN(line, from, 1)[0])
// // output = strings.SplitN(output, til, 1)[0]
// // output = strings.SplitN(output, til, 1)[0]
// // return output
// }
// }
// return ""
// }
