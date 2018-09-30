function playAudio () {
  get('render.play.audio').play();
  get('render.play.button').src = 'img/pause.png';
}

function pauseAudio () {
  get('render.play.audio').pause();
  get('render.play.button').src = 'img/play.png';
}
