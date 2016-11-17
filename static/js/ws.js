// Websocket Helpers
var ws; // the websocket itself
var createSuccess = false;
var mbInfo = {
  "roomName": "TylerTest",
  "id": "TylerTest"
}

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

function getNameFromId(id) {
  var msg = {
    Command: "id",
    Body: [id]
  };

  ws.send(JSON.stringify(msg));
}

function parseResponse(r) {
  res = JSON.parse(r);
  console.log(res);

  if (res.Command == "skip") {
    AudioEndedHandler();
    return "skipped";
  } else if (res.Command == "get") {
    if (res.Body[0] == "FAIL") {
      return "FAIL";
    }
    
    for (var i = 0; i < res.Body.length; i++) {
      if (playQueue.indexOf(res.Body[i]) == -1) {
        playQueue.push("http://austindoes.work/song/" + res.Body[i]);
      }
    }
    return "got";
  } else if (res.Command == "join") {
    if (res.Body[0] != "FAIL") {
      createSuccess = true;
      var partyText = "you're hosting \"" + res.Body[0] + "\"";
      document.getElementById("partyName").innerHTML = partyText;
    }
    return res.Body[0];
  } else if (res.Command == "id") {
    currentSongname = res.Body[0];
    OnSongNameChange();
    return res.Body[0];
  }
}
