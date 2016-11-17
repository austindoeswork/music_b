// Websocket Helpers
var ws; // the websocket itself
var createSuccess = false;
var gotFirst = false;
var qLength = "0";
var mbInfo = {
  "roomName": "",
  "id": ""
};

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

function nextSong() {
  var msg = {
    Command: "next",
    Body: [mbInfo.id]
  };

  ws.send(JSON.stringify(msg));
}

function joinRoom() {
  var msg = {
    Command: "join",
    Body: [mbInfo.roomName]
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

function rejoinRoom() {
  var msg = {
    Command: "rejoin",
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

function getQueueLength() {
  var msg = {
    Command: "length",
    Body: [mbInfo.roomName]
  };

  ws.send(JSON.stringify(msg));
}

function parseResponse(r) {
  res = JSON.parse(r);

  if (res.Command == "skip") {
    AudioEndedHandler();
    return "skipped";
  } else if (res.Command == "length") {
    if (res.Body[0] != "FAIL" && res.Body.length != 0) {
      qLength = res.Body[0];
    }
    return res.Body[0];
  } else if (res.Command == "get" || res.Command == "next") {
    if (res.Body[0] != "FAIL" && res.Body.length != 0) {
      playQueue = [];
      playQueue.push("http://austindoes.work/song/" + res.Body[0]);
      return "got";
    } else {
      return "FAIL";
    }
  } else if (res.Command == "create") {
    if (res.Body[0] != "FAIL" && res.Body.length != 0) {
      createSuccess = true;
      var partyText = "you're in \"" + res.Body[0] + "\"";
      document.getElementById("partyName").innerHTML = partyText;
    }
    return res.Body[0];
  } else if (res.Command == "join" || res.Command == "rejoin") {
    if (res.Body[0] != "FAIL" && res.Body.length != 0) {
      createSuccess = true;
      var partyText = "you're hosting \"" + res.Body[0] + "\"";
      document.getElementById("partyName").innerHTML = partyText;
    }
    return res.Body[0];
  } else if (res.Command == "id") {
    if (res.Body[0] != "FAIL" && res.Body.length != 0) {
      currentSongname = res.Body[0];
      OnSongNameChange();
    }
    return res.Body[0];
  } else if (res.Command == "terminate") {
    // uh oh
  }
}
