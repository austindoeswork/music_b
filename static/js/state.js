const Pages = {
  home: '',
  control: '#/control/:id',
  play: '#/play/:id'
};

function initState () {
  localStorage.state = {
    page: 'home',
    room: {
      id: null,
      name: '',
      createSuccess: false,
      joinSuccess: false,
    }
    roomId: null,
    roomName: null,
    play: {
      queue: [],
      qLength: '0',
      audioLoading: false,
      currentSongname: null,
    },
    api:  {
      ws: null,
      gotFirst: false,
    },
    render: {
      home: {
        el: document.getElementById('splash'),
        title: document.getElementById('splash-title'),
        input: document.getElementById('partyNameInput'),
        submit: document.getElementById('partyNameSubmit'),
        fail: document.getElementById('splash-fail'),
      },
      control: {
        el: document.getElementById('client'),
      },
      play: {
        el: document.getElementById('host'),
        title: document.getElementById('host-title'),
        loading: document.getElementById('play-loading'),
        audio: document.getElementById('audio'),
        canvas: document.getElementById('canvas'),
        button: document.getElementById('playButton'),
        partyName: document.getElementById('host-partyName'),
        songName: document.getElementById('host-songName'),
      },
    }
  }
}

get (item) {
  if (typeof localStorage.state === 'undefined') {
    return undefined;
  }

  // Get an item from the state
  const path = item.split('.');

  out = localStorage.state;
  for (let i = 0; i < path.length; i++) {
    out = out[path[i]]
  }

  return out;
}

set (item, to) {
  if (typeof localStorage.state === 'undefined') {
    return undefined;
  }

  // Change an item in the state
  const path = item.split('.');

  out = localStorage.state;
  for (let i = 0; i < path.length-1; i++) {
    out = out[path[i]];
  }

  out[path.pop()] = to;
}

changePage (to) {
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
      get(`render.${Pages[k]}.el`).classList.toggle('hidden', true);
      // set the new page in the state
      set('page', Pages[k]);

      return Pages[k];
    }
  }

  return null;
}
