function playAudio () {
  get('render.host.audio').play();
  get('render.host.button').src = 'img/pause.png';
}

function pauseAudio () {
  get('render.host.audio').pause();
  get('render.host.button').src = 'img/play.png';
}
