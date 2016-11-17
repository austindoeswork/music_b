# music_b

anarchistic music player

### cache
holds important info in memory

- cache
	- parties -- associates parties to songs
	- threads -- associates listener threads to parties
	- songs -- holds song metadata
- party
	- holds info about player and songs

### commander
speaks to the music player thru a websocket

serves the static content

### downloader
downloads music from different sources

- depends on youtube-dl python script
```
sudo curl -L https://yt-dl.org/downloads/latest/youtube-dl -o /usr/local/bin/youtube-dl
```
```
sudo chmod a+rx /usr/local/bin/youtube-dl
```

### listener
listens to peoples commands through fb or otherwise

- Websocket - TODO
- Facebook - Usain Bott
- GroupMe ?
- Txt Msg ?

### router
routes messages and paths to respective function calls

### server
serves music from a directory

### TODO:
- backend
	- websocket party listener
	- slack party listener
	- proper logging
	- rejoin controller
- front end
	- music control thru player page
	- nicely show info about listeners in play.html
	- mobile debugging/app
	- ellipses animation