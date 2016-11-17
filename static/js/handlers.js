function ButtonClickHandler() {
  if (audioLoading) {
    return;
  }

  var button = document.getElementById("playButton");
  var playing = button.src.includes("pause");

  if (playing) {
    pauseAudio();
  } else {
    playAudio();
  }
}

function AudioLoadedHandler() {
  audioLoading = false;
}

function AudioEndedHandler() {
  audioLoading = true;

  var audio = document.getElementById("audio");
  var source = document.getElementById("source");

  playQueue.pop();
  var srcUrl = playQueue[0];
  requestSong();

  source.src = srcUrl;
  audio.pause();
  audio.load();
  playAudio();
}

function BodyReadyHandler() {
  var audio = document.getElementById("audio");
  audio.addEventListener("ended", AudioEndedHandler);
  window.setTimeout(function(){document.getElementById("audio").play();}, 700);

  var canvas = document.getElementById("canvas");
  canvas.width = window.innerWidth;
  canvas.height = window.innerHeight;

  var audioCtx = new (window.AudioContext || window.webkitAudioContext)();
  var source = audioCtx.createMediaElementSource(audio);
  var analyser = audioCtx.createAnalyser();
  source.connect(analyser);
  analyser.connect(audioCtx.destination);
  analyser.fftSize = 2048;

  draw(analyser);
}
