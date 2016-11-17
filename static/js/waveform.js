
// draw an oscilloscope of the current audio source
function draw(analyser) {

  var drawVisual = requestAnimationFrame(function(){
    draw(analyser);
  });

  var SIZE = window.innerWidth;
  var HG = window.innerHeight;
  var DEV = 10.0;
  var SLICE = 100.0;
  var YMOD = 10.0;
  var DXMOD = 2.5;
  var DYMOD = 8.0;
  var CENTERX, CENTERY;
  var RAD;

  RAD = document.getElementById("playButton").getBoundingClientRect().width;
  RAD /= 2;

  var rect = document.getElementById("playButton").getBoundingClientRect();
  CENTERX = SIZE/2.0;
  CENTERY = HG/2.0;

  var bufferLength = analyser.frequencyBinCount;
  var dataArray = new Uint8Array(bufferLength);

  var canvas = document.getElementById("canvas");
  var canvasCtx = canvas.getContext("2d");

  analyser.getByteTimeDomainData(dataArray);

  canvasCtx.fillStyle = 'rgba(48, 48, 48, 1)';
  canvasCtx.fillRect(0, 0, canvas.width, canvas.height);
  canvasCtx.lineWidth = 2;
  canvasCtx.strokeStyle = 'rgb(255, 160, 64)';

  canvasCtx.beginPath();

  var points = [];
  for (var i = 0; i < bufferLength; i++) {
    var v = dataArray[i] / 4.0;

    var theta = 2*Math.PI*(i/bufferLength);
    var r = v + RAD;
    var dx = CENTERX + r*Math.cos(theta);
    var dy = CENTERY + r*Math.sin(theta);

    points.push({"x": dx, "y": dy});
  }

  canvasCtx.moveTo(points[0].x, points[0].y);
  for (i = 1; i < points.length - 2; i++) {
    var xc = (points[i].x + points[i + 1].x) / 2;
    var yc = (points[i].y + points[i + 1].y) / 2;
    canvasCtx.quadraticCurveTo(points[i].x, points[i].y, xc, yc);
  }
  var lxc = (points[0].x + points[i+1].x)/2;
  var lyc = (points[0].y + points[i+1].y)/2;

  canvasCtx.quadraticCurveTo(points[i+1].x, points[i+1].y, lxc, lyc);
  canvasCtx.quadraticCurveTo(points[0].x, points[0].y, lxc, lyc);

  canvasCtx.stroke();
};
