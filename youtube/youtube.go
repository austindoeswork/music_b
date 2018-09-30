package youtube

import (
	"net/http"
	"os/exec"
	"strings"

	"golang.org/x/net/html"
	// "github.com/rylio/ytdl"
)

var (
	ytSearchUrl = "https://www.youtube.com/results?search_query="
	ytWatchUrl  = "https://www.youtube.com/watch?v="
)

func Download(ytID, filePath string) error {
	var err error
	audioUrl := ytWatchUrl + ytID

	cmd := exec.Command(
		"youtube-dl",
		"--extract-audio",
		"--audio-format", "mp3",
		"-o", filePath,
		audioUrl)

	err = cmd.Start()
	if err != nil {
		return err
	}
	err = cmd.Wait()
	return err
}

func Query(query string, maxHits int) []string {
	query = strings.Replace(query, " ", "+", -1)
	res, err := http.Get(ytSearchUrl + query)
	if err != nil {
		panic(err)
		return nil
	}
	defer res.Body.Close()

	tokenizer := html.NewTokenizer(res.Body)
	var vids []string
	count := 0
	for {
		tt := tokenizer.Next()
		switch {
		case tt == html.ErrorToken:
			return vids // End of the document, we're done
		case tt == html.StartTagToken:
			t := tokenizer.Token()
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
					}
				} // for range t.Attr
			} // if t.Data == "div"
		} // switch
	}

	return nil
}
