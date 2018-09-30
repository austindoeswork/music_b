function playAudio () {
  set('play.loading', false);
  get('render.host.audio').play();
  get('render.host.button').src = 'img/pause.png';
}

function pauseAudio () {
  set('play.loading', true);
  get('render.host.audio').pause();
  get('render.host.button').src = 'img/play.png';
}
