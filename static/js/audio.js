var audioLoading = true;
var playQueue = [];

function playAudio() {
  var button = document.getElementById("playButton");
  var audio = document.getElementById("audio");
  audio.play();
  button.src = button.src.replace(/play/i, "pause");
}

function pauseAudio() {
  var button = document.getElementById("playButton");
  var audio = document.getElementById("audio");

  audio.pause();
  button.src = button.src.replace(/pause/i, "play");
}
