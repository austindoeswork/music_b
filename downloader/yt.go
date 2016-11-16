package downloader

import (
	"errors"
	// "fmt"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"golang.org/x/net/html"

	"github.com/rylio/ytdl"
)

var (
	ytSearchUrl = "https://www.youtube.com/results?search_query="
	ytWatchUrl  = "https://www.youtube.com/watch?v="
)

//TODO make this clean itself up somehow

type YTDownloader struct {
	dir string
}

func NewYTDownloader(dir string) *YTDownloader {
	return &YTDownloader{
		dir,
	}
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