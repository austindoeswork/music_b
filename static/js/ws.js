// Websocket Helpers

// init
function initWs (url, port) {
  url = 'wss://' + url;
  if (port != '') {
    url += ':' + port;
  }

  set('ws', new WebSocket(url));
}

// commands we can use
function requestSong (query) {
  get('ws').send(JSON.stringify({
    type: 'push',
    body: [query]
  }));
}

function nextSong () {
  get('ws').send(JSON.stringify({
    type: 'pop',
    body: []
  }));
}

function parseResponse (r) {
  res = JSON.parse(r);

  if (res.type == 'queue') {
    set('play.queue', res.Body);
  } else if (res.type == 'error') {
    errs = get('errors');
    errs.push(res.Body);
  }
}
