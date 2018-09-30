function draw (analyser) {
  const drawVisual = requestAnimationFrame(function () {
    draw(analyser);
  });

  // Set up constants for the waveform calculation
  const SIZE = window.innerWidth;
  const HG = window.innerHeight;
  const DEV = 10.0;
  const SLICE = 100.0;
  const YMOD = 10.0;
  const DXMOD = 2.5;
  const DYMOD = 8.0;
  const RAD = document.getElementById('playButton').getBoundingClientRect().width / 2;
  const CX = SIZE/2.0;
  const CY = HG/2.0;

  // dataArray contains the raw audio data for some small time duration
  const bufferLength = analyser.frequencyBinCount;
  let dataArray = new Uint8Array(bufferLength);
  analyser.getByteTimeDomainData(dataArray);

  // generate the points of the waveform from the audio data by translating them into polar coordinates
  // angle is the position in time, length is the amplitude of the waveform at that position in time
  let points = [];
  for (let i = 0; i < bufferLength; i++) {
    const v = dataArray[i] / 4.0;
    const theta = 2*Math.PI*(i/bufferLength);

    const r = v + RAD;
    const dx = CX + r*Math.cos(theta+Math.PI/2);
    const dy = CY + r*Math.sin(theta+Math.PI/2);

    points.push({x: dx, y: dy});
  }

  // clear the canvas, set stroke
  const canvas = get('render.play.canvas');
  const canvasCtx = canvas.getContext('2d');

  canvasCtx.fillStyle = 'rgba(48, 48, 48, 1)';
  canvasCtx.fillRect(0, 0, canvas.width, canvas.height);
  canvasCtx.lineWidth = 2;
  canvasCtx.strokeStyle = 'rgb(255, 160, 64)';

  canvasCtx.beginPath();

  // now create a curve from the previously calculated points
  canvasCtx.moveTo(points[0].x, points[0].y);
  for (i = 1; i < points.length - 2; i++) {
    const dx = (points[i].x + points[i + 1].x) / 2;
    const dy = (points[i].y + points[i + 1].y) / 2;
    canvasCtx.quadraticCurveTo(points[i].x, points[i].y, dx, dy);
  }

  // final dx, dy. this is to close the loop
  const fdx = (points[0].x + points[i+1].x)/2;
  const fdy = (points[0].y + points[i+1].y)/2;

  canvasCtx.quadraticCurveTo(points[i+1].x, points[i+1].y, fxc, fyc);
  canvasCtx.quadraticCurveTo(points[0].x, points[0].y, fxc, fyc);

  canvasCtx.stroke();
};
