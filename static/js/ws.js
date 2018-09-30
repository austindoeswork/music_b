// Websocket Helpers

// init
function initWs (url, port, name, cb) {
  if (PROTOCOL == 'http') {
    url = 'ws://' + url;
  } else if (PROTOCOL == 'https') {
    url = 'wss://' + url;
  }

  if (!!port) {
    url += ':' + port;
  }

  url += '/ws?name=' + name;

  ws = new WebSocket(url);
  set('ws', ws);

  ws.onopen = function (e) {
    get('render.host.partyName').innerHTML = `you're hosting "${get('room.name')}"`;

    BodyReadyHandler();

    requestSong('all star');
  };

  ws.onmessage = function(e) {
    parseResponse(e.data);
  }
}

// commands we can use
function requestSong (query) {
  get('ws').send(JSON.stringify({
    type: 'push',
    body: query,
  }));
}

function nextSong () {
  get('ws').send(JSON.stringify({
    type: 'pop',
    body: null,
  }));
}

function parseResponse (r) {
  res = JSON.parse(r);

  if (res.type == 'queue') {
    set('play.queue', res.Body);
    if (get('play.loading')) {
      AudioEndedHandler();
    }
  } else if (res.type == 'error') {
    errs = get('errors');
    errs.push(res.Body);
  }
}
