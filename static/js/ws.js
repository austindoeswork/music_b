// Websocket Helpers
var ws; // the websocket itself
var mbInfo = {
  "roomName": "TylerTest",
  "id": "TylerTest"
}

createWS("austindoes.work/ws", "");

function createWS(url, port) {
  var loc = "ws://" + url;
  if (port != "") {
    loc += ":" + port;
  }

  ws = new WebSocket(loc);
  mbInfo.url = loc;
}

function requestSong() {
  var msg = {
    Command: "get",
    Body: [mbInfo.id]
  };

  ws.send(JSON.stringify(msg));
}

function createRoom() {
  var msg = {
    Command: "join",
    Body: [mbInfo.roomName]
  };

  ws.send(JSON.stringify(msg));
}

function parseResponse(r) {
  res = JSON.parse(r);

  if (res.Command == "skip") {
    AudioEndedHandler();
  } else if (res.Command == "res") {
    if (res.body.length == 1) {
      return; // response after joining
    }
    // else
    for (var i = 0; i < res.Body.length; i++) {
      if (playQueue.indexOf(res.Body[i]) == -1) {
        playQueue.push("http://austindoes.work/song/" + res.Body[i]);
      }
    }
  }
}
