package commander

import (
	"fmt"
	"log"
	// "html"
	"html/template"
	"net/http"
	// "html"
	// "net/http"
	// "golang.org/x/net/websocket"
	//TODO vendor gorilla
	// or maybe switch to golang.org? idk
	"github.com/gorilla/websocket"
)

// play(party)
// pause(party)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Player holds the connection with the audio player
type Player struct {
	//maybe keep this info in the cache
	party string //how to populate (or do we need to?)
	conn  *websocket.Conn
}

func (p *Player) listenWS() {
	defer func() {
		p.conn.Close()
	}()
	for {
		//pause
		//join
		//ok ? TODO figure this out (ack, sync)
		//status ?
		//getsongs
		//deletesong
		//
		mt, msg, err := p.conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", msg)
		err = p.conn.WriteMessage(mt, msg)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

type Commander struct {
}

func New() *Commander {
	return &Commander{}
}

func home(w http.ResponseWriter, r *http.Request) {
	homeTemplate.Execute(w, "ws://"+"austinsland.ottoq.com:8888"+"/echo")
}

func (c *Commander) Listen() {
	http.HandleFunc("/test", home)
	http.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("ws Request received.")
		serveWS(w, r)
	})

	fmt.Println("starting commander.")
	err := http.ListenAndServe(":8888", nil)
	if err != nil {
		panic("Commander: " + err.Error())
	}
}

func serveWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		panic("Commander: " + err.Error())
	}
	p := &Player{conn: conn}
	p.listenWS()
}

// type PStream struct {
// }

// type WSHandler func(*Conn)

// func checkOrigin(config *Config, req *http.Request) (err error) {
// config.Origin, err = Origin(config, req)
// if err == nil && config.Origin == nil {
// return fmt.Errorf("null origin")
// }
// return err
// }

// // ServeHTTP implements the http.Handler interface for a WebSocket
// func (h Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
// s := Server{Handler: h, Handshake: checkOrigin}
// s.serveWebSocket(w, req)
// }

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
	<p><input id="input" type="text" value=".fuckjon">
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
