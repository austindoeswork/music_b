const Pages = {
  home: '',
  host: '#/host/:id',
  client: '#/client/:id'
};

var __STATE;

function initState () {
  __STATE = {
    page: 'home',
    errors: [],
    room: {
      id: null,
      name: '',
      createSuccess: false,
      joinSuccess: false,
    },
    roomId: null,
    roomName: null,
    play: {
      queue: [],
      loading: true,
      currentSongname: null,
    },
    ws: null,
    render: {
      home: {
        el: document.getElementById('splash'),
        title: document.getElementById('splash-title'),
        input: document.getElementById('partyNameInput'),
        submit: document.getElementById('partyNameSubmit'),
        fail: document.getElementById('splash-fail'),
      },
      client: {
        el: document.getElementById('client'),
      },
      host: {
        mode: 'songInfo',
        el: document.getElementById('host'),
        title: document.getElementById('host-title'),
        audio: document.getElementById('audio'),
        source: document.getElementById('source'),
        canvas: document.getElementById('canvas'),
        button: document.getElementById('playButton'),

        swappyIcon: document.getElementById('swappy-icon'),

        songInfo: document.getElementById('host-song-info'),
        partyName: document.getElementById('host-partyName'),
        songName: document.getElementById('host-songName'),

        control: document.getElementById('host-control'),
        controller: document.getElementById('host-control-input'),

        skip: document.getElementById('skippy-dippy-doo'),
      },
    },
  };
}

function get (item) {
  if (typeof __STATE === 'undefined') {
    return undefined;
  }

  // Get an item from the state
  const path = item.split('.');

  out = __STATE;
  for (let i = 0; i < path.length; i++) {
    out = out[path[i]]
  }

  return out;
}

function set (item, to) {
  if (typeof __STATE === 'undefined') {
    return undefined;
  }

  // Change an item in the state
  const path = item.split('.');

  out = __STATE;
  for (let i = 0; i < path.length-1; i++) {
    out = out[path[i]];
  }

  out[path.pop()] = to;
}

function changePage (to) {
  const parts = to.split('/')
  for (let k in Pages) {
    const against = Pages[k].split('/');
    let found = true
    for (let i = 0; i < against.length; i++) {
      if (against[i][0] == ':') {
        continue;
      } else {
        if (against[i] != parts[i]) {
          found = false;
          break;
        }
      }
    }

    if (found) {
      // hide old page
      get(`render.${get('page')}.el`).classList.toggle('hidden', true);
      // show the new page
      get(`render.${k}.el`).classList.toggle('hidden', false);
      // set the new page in the state
      set('page', Pages[k]);

      return Pages[k];
    }
  }

  return null;
}

function host_swappado (to) {
  const cur = get('render.host.mode');

  if (typeof to === 'undefined') {
    if (cur == 'songInfo') {
      to = 'controller';
    } else if (cur == 'controller') {
      to = 'songInfo';
    }
  } else if (cur == to) {
    return;
  }

  let old, yung; // new is a keyword :/

  if (to == 'controller') {
    old = get('render.host.songInfo');
    yung = get('render.host.control');

    get('render.host.swappyIcon').classList.toggle('fa-calculator', false);
    get('render.host.swappyIcon').classList.toggle('fa-music', true);

    get('render.host.controller').focus();
  } else if (to == 'songInfo') {
    old = get('render.host.control');
    yung = get('render.host.songInfo');

    get('render.host.swappyIcon').classList.toggle('fa-calculator', true);
    get('render.host.swappyIcon').classList.toggle('fa-music', false);
  }

  old.classList.toggle('rotatedOut', true);
  yung.classList.toggle('rotatedOut', false);

  set('render.host.mode', to);
}
