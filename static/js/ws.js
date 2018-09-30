// Websocket Helpers

// init
function initWs (url, port) {
  let url = 'ws://' + url;
  if (port != '') {
    url += ':' + port;
  }

  set('ws', {
    ws: new WebSocket(url)
    createSuccess: false,
    gotFirst: false,
    qLength: '0',
  });
}

// commands we can use
function requestSong () {
  ws.send(JSON.stringify({
    Command: 'get',
    Body: [get('room.id')]
  }));
}

function nextSong () {
  ws.send(JSON.stringify({
    Command: 'next',
    Body: [get('room.id')]
  }));
}

function joinRoom () {
  ws.send(JSON.stringify({
    Command: 'join',
    Body: [get('room.name')]
  }));
}

function createRoom () {
  ws.send(JSON.stringify({
    Command: 'join',
    Body: [get('room.name')]
  }));
}

function rejoinRoom () {
  ws.send(JSON.stringify({
    Command: 'rejoin',
    Body: [get('room.name')]
  }));
}

function getNameFromId (id) {
  ws.send(JSON.stringify({
    Command: 'id',
    Body: [id]
  }));
}

function getQueueLength () {
  ws.send(JSON.stringify({
    Command: 'length',
    Body: [get('room.name')]
  }));
}

function closeParty () {
  return;
}

function parseResponse (r) {
  res = JSON.parse(r);

  if (res.Command == 'skip') {
    AudioEndedHandler();
    return 'skipped';
  } else if (res.Command == 'length') {
    if (res.Body[0] != 'FAIL' && res.Body.length != 0) {
      set('play.qLength', res.Body[0]);
    }
    return res.Body[0];
  } else if (res.Command == 'get' || res.Command == 'next') {
    if (res.Body[0] != 'FAIL' && res.Body.length != 0) {
      playQueue = [];
      set('play.queue', ['http://austindoes.work/song/' + res.Body[0]])
      return 'got';
    } else {
      return 'FAIL';
    }
  } else if (res.Command == 'create') {
    if (res.Body[0] != 'FAIL' && res.Body.length != 0) {
      set('room.createSuccess', true);
      const text = 'you\'re hosting "'
                   + res.Body[0]
                   + '" <i class="fa fa-times" onclick="OnCloseParty();" aria-hidden="true"></i>';
      get('render.control.partyName').innerHTML = text;
    }
    return res.Body[0];
  } else if (res.Command == 'join' || res.Command == 'rejoin') {
    if (res.Body[0] != 'FAIL' && res.Body.length != 0) {
      set('room.joinSuccess', true);
      get('render.control.partyName').innerHTML = 'you\'re in "' + res.Body[0] + '"';
    }
    return res.Body[0];
  } else if (res.Command == 'id') {
    if (res.Body[0] != 'FAIL' && res.Body.length != 0) {
      set('play.currentSongname', res.Body[0]);
      OnSongNameChange();
    }
    return res.Body[0];
  } else if (res.Command == 'terminate') {
    localStorage.clear();
    mbInfo = {};
    window.location = '?kicked';
  }
}
