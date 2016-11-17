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

function checkQueueReady() {
  if (playQueue.length > 0) {
    AudioEndedHandler();
  } else {
    window.setTimeout(checkQueueReady, 200);
  }
}

function AudioEndedHandler() {
  if (playQueue.length == 0) {
    requestSong();
    checkQueueReady();
  }

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

function CheckRoomJoin(depth) {
  if (true) {
    document.getElementById("createpage").style.display = 'none';
    document.getElementById("loading").style.display = 'none';
    document.getElementById("playerpage").style.display = 'block';

    document.getElementById("playButton").style.webkitAnimationPlayState = 'running';
    document.getElementById("title").style.webkitAnimationPlayState = 'running';
    document.getElementById("canvas").style.webkitAnimationPlayState = 'running';

    BodyReadyHandler();

    requestSong();
  } else if (depth > 0) {
    window.setTimeout(function(){CheckRoomJoin(depth-1);}, 200);
  } else {
    window.location = "?fail";
  }
}

function TryCreatingRoom() {
  var roomname = document.getElementById("submit").value;
  mbInfo.roomName = roomname;
  mbInfo.id = roomname
  createWS("austindoes.work/ws", "");

  ws.onopen = function(e) {
    ws.onmessage = function(e) {
      parseResponse(e.data);
    }

    createRoom();
    CheckRoomJoin();
  }
}
