package commander

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/austindoeswork/music_b/cache"
	"github.com/gorilla/websocket"
	"github.com/wardn/uuid"
)

//TODO make this not a fucking switch case!!!!!!
//TODO better err handling

//TODO better implementation for this

type Player struct {
	c     *cache.Cache
	id    string
	party string
	conn  *websocket.Conn
}

type PlayerCommand struct {
	Command string
	Body    []string
}

// { "Command": "asdfasdf", "Body":["asdfasdf asdd"]}
func (p *Player) listenWS() {
	defer func() {
		p.c.EndParty(p.party)
		p.conn.Close()
	}()
	for {
		mt, msg, err := p.conn.ReadMessage()
		if err != nil {
			log.Println("error reading ws msg (" + p.party + ")")
			break
		}
		var cmd PlayerCommand
		err = json.Unmarshal(msg, &cmd)
		if err != nil {
			log.Println("error unmarshalling command (" + p.party + ")")
			continue
		}
		fmt.Println("ws: (" + p.party + ") " + cmd.Command)
		switch cmd.Command {
		case "join":
			if len(cmd.Body) > 0 {
				partyName := cmd.Body[0]
				encodedName := p.c.GetEncodedName(partyName)
				_, err := p.c.MakeParty(partyName)

				if err != nil {
					p.party = encodedName
					p.respond("res", "HIJACKING")
				} else {
					p.party = encodedName
					p.respond("res", partyName)
				}
				err = p.c.AddPlayer(encodedName, p.id)
				if err != nil {
					log.Println("error adding player to cache")
				}
				continue
			} else {
				p.respond("res", "FAIL")
				continue
			}
		case "get":
			songs, err := p.c.GetSongs(p.party, 2)
			if err != nil {
				p.respond("res", "FAIL")
				continue
			} else {
				p.respond("res", songs...)
				continue
			}
		case "next":
			err := p.c.PopSong(p.party)
			if err != nil {
				p.respond("res", "FAIL")
				continue
			}
			songs, err := p.c.GetSongs(p.party, 2)
			if err != nil {
				p.respond("res", "FAIL")
				continue
			} else {
				p.respond("res", songs...)
				continue
			}
		default:
		}

		err = p.conn.WriteMessage(mt, msg)
		if err != nil {
			//TODO retry instead of break
			log.Println("error writing to ws, closing connection (" + p.party + ")")
			break
		}
	}
}

// func (p *Player) Skip() error {
// err := p.respond("skip")
// return err
// }

// func (p *Player) Pause() error {
// err := p.respond("pause")
// return err
// }

// func (p *Player) Resume() error {
// err := p.respond("resume")
// return err
// }

func (p *Player) respond(cmd string, body ...string) error {
	res := PlayerCommand{
		cmd,
		body,
	}
	err := p.conn.WriteJSON(res)
	if err != nil {
		return err
	}
	return nil
}

func (p *Player) command(cmd PlayerCommand) error {
	err := p.conn.WriteJSON(cmd)
	return err
}

type Commander struct {
	c       *cache.Cache
	players map[string]*Player
}

func New(c *cache.Cache) *Commander {
	return &Commander{
		c,
		map[string]*Player{},
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	homeTemplate.Execute(w, "ws://"+"austinsland.ottoq.com:8888"+"/echo")
}

func (c *Commander) Listen(port string) {
	http.HandleFunc("/test", home)
	http.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("ws Request received.")
		c.serveWS(w, r)
	})

	fmt.Println("Started commander @ " + port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Panic("Commander: " + err.Error())
	}
}

func (c *Commander) serveWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Panic("Commander: " + err.Error())
	}
	playerID := uuid.New()

	p := &Player{
		c.c,
		playerID,
		"",
		conn,
	}
	c.players[playerID] = p
	p.listenWS()

	fmt.Println("serve ws ended")
	delete(c.players, playerID)
}

func (c *Commander) Command(playerID string, cmd PlayerCommand) error {
	if _, ok := c.players[playerID]; !ok {
		return errors.New("No such player:" + playerID)
	}
	err := c.players[playerID].command(cmd)
	return err
}

// party --
// 		   \
// party ---- usain
// 		   /
// party --

var homeTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<head>
<meta charset="utf-8">
<style>
	body {
		font-size: 14px;
		font-family: "Courier New", Courier, monospace;
	}
	td {
		margin: 5px;
		border: 1px dashed black;
	}
</style>
<script>  
window.addEventListener("load", function(evt) {
    var output = document.getElementById("output");
    var input = document.getElementById("input");
    var ws;
    var print = function(message) {
        var d = document.createElement("div");
        d.innerHTML = message;
        output.appendChild(d);
    };
	var clear = function() {
		output.innerHTML = '';
    };
    document.getElementById("open").onclick = function(evt) {
        if (ws) {
            return false;
        }
        ws = new WebSocket("{{.}}");
        ws.onopen = function(evt) {
            print("WS CONNECTION OPENED");
        }
        ws.onclose = function(evt) {
            print("WS CONNECTION CLOSED");
            ws = null;
        }
        ws.onmessage = function(evt) {
            print("RES: " + evt.data);
        }
        ws.onerror = function(evt) {
            print("ERR: " + evt.data);
        }
        return false;
    };
    document.getElementById("send").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        print("MSG: " + input.value);
        ws.send(input.value);
        return false;
    };
    document.getElementById("close").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        ws.close();
        return false;
    };
	document.getElementById("clear").onclick = function(evt) {
        clear();
		return false;
    };

});
</script>
</head>
<body>
<h3>controller</h3>
<table>
<tr><td valign="top" width="50%">
	<p>Open - create a ws connection</p>
	<p>Close - close the ws connection</p>
	<p>Send - send a message through the ws</p>
<form>
	<button id="open">Open</button>
	<button id="close">Close</button>
	<p><input id="input" type="text" value=".fuckjon"></p>
	<p>
	<button id="send">Send</button>
	<button id="clear">Clear</button>
</form>
</td>
	<td valign="top" width="50%">
	<div id="output"></div>
</td></tr></table>
<h3>player</h3>
todo
</body>
</html>
`))

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}
