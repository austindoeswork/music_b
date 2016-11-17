var currentSongname = "welcome to music_b";
var audioLoading = true;
var playQueue = [];

function playAudio() {
  var button = document.getElementById("playButton");
  var audio = document.getElementById("audio");
  audio.play();
  button.src = "img/pause.png";
}

function pauseAudio() {
  var button = document.getElementById("playButton");
  var audio = document.getElementById("audio");

  audio.pause();
  button.src = "img/play.png";
}
