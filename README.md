# music_b

### install:

- `cd $GOPATH/src/github.com/austindoeswork/`
- `git clone https://github.com/austindoeswork/music_b.git`

### dependencies:

go packages

- `go get`

youtube-dl

- `sudo curl -L https://yt-dl.org/downloads/latest/youtube-dl -o /usr/local/bin/youtube-dl`
- `sudo chmod a+rx /usr/local/bin/youtube-dl`

### setting up:

get ssl cert/key for https (only if running https locally)

- run these commands
- `openssl genrsa -out server.key 2048`
- `openssl ecparam -genkey -name secp384r1 -out server.key`
- `openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650`
- should create
	- `server.crt` and `server.key`



create a config

- from root dir (`github.com/austindoeswork/music_b/`) 
	- touch `config.json`
- example `config.json`

```{
{
  "ServerAddress": ":8100",
  "Secure": false,
  "SSLCert": "server.crt",
  "SSLKey": "server.key"
}
```

run the server

- `go run main.go` or `go build .;./music_b`

### deploying:

remote_address: `austindoes.work`

- ensure server is not running
- `./bin/deploy`
- `ssh austindoes.work`
- `$user@austindoeswork> sudo ~/music_b`

