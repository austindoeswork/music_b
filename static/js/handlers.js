function BodyReadyHandler () {
  window.history.pushState(null, 'music_b', '#/host/' + get('room.name'));
  changePage(window.location.hash);

  const audio = document.getElementById('audio');
  audio.addEventListener('ended', AudioEndedHandler);
  window.setTimeout(function(){get('render.host.audio').play();}, 700);

  const canvas = get('render.host.canvas');
  canvas.width = window.innerWidth;
  canvas.height = window.innerHeight;

  const audioCtx = new (window.AudioContext || window.webkitAudioContext)();
  const source = audioCtx.createMediaElementSource(audio);
  const analyser = audioCtx.createAnalyser();
  source.connect(analyser);
  analyser.connect(audioCtx.destination);
  analyser.fftSize = 2048;

  draw(analyser);
}

function AudioLoadedHandler () {
  // this.AudioEndedHandler();
}

function AudioEndedHandler () {
  const audio = get('render.host.audio');
  audio.pause();

  const button = get('render.host.button');
  if (get('play.queue').length == 0 ) {
    button.src = 'img/sad.png';
    set('play.currentSongname', 'No songs in queue :[');
    set('play.loading', true);
  } else {
    const source = get('render.host.source');

    let currentSong = get('play.queue').shift();
    const srcUrl = PROTOCOL + '://' + URL + '/song/' + currentSong.id;

    if (typeof srcUrl != 'undefined') {
      source.src = srcUrl;
      audio.load();
      playAudio();

      set('play.currentSongname', currentSong.query);
    }
  }

  get('render.host.songName').innerHTML = get('play.currentSongname');
}
