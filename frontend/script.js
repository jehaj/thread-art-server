class Point {
  constructor(x, y) {
    this.x = x;
    this.y = y;
  }
}

function constrain(x, min, max) {
  if (x < min) x = min;
  if (max < x) x = max;
  return x;
}

const canvasSize = 400;
const takeIn = 20;
const lineWidth = 1;
const transparency = 80;
let lineColor;


const ctx = document.getElementById("t-canvas").getContext("2d");

ctx.globalCompositeOperation = "destination-over";
ctx.strokeStyle = "rgba(0, 0, 0, 0.2)";

let runButton = document.querySelector("#t-button");

let pointer = 0;
let circleSize = 0;
let pointList = [];
let circle = new Map();

fetch(window.location.origin + "/example.txt").then((data) => {
  data.blob().then((dataBlob => {
    dataBlob.text().then(res => {
      circleSize = parseInt(res.split("\n")[0]);
      for (let point of res.split("\n")[1].slice(0, -1).split(",")) {
        pointList.push(parseInt(point));
      }
      for (let i = 0; i < circleSize; i++) {
        let angle = (2 * Math.PI / circleSize) * i;
        let x = constrain(
          (Math.cos(angle) * canvasSize) / 2 + canvasSize / 2,
          0,
          canvasSize - 1
        );
        let y = constrain(
          (Math.sin(angle) * canvasSize) / 2 + canvasSize / 2,
          0,
          canvasSize - 1
        );
        let tPoint = new Point(x, y);
        circle.set(i, tPoint);
      }
      ctx.beginPath();
      for (let i = 0; i < pointList.length; i++) {
        let t_point = circle.get(pointList[pointer]);
        ctx.lineTo(t_point.x, t_point.y);
        pointer++;
      }
      ctx.stroke();
    })
  }));
})

let doLoop = false;

runButton.addEventListener("click", (event) => {
  ctx.clearRect(0, 0, 400, 400); // clear canvas
  if (runButton.value == "Run") {
    pointer = 0;
    runButton.value = "Reset";
    doLoop = true;
    window.requestAnimationFrame(draw);
  } else {
    runButton.value = "Run";
    doLoop = false;
  }
});

function draw() {
  ctx.beginPath();
  for (let i = 0; i < takeIn; i++) {
    if (pointer >= pointList.length) {
      break;
    }
    let t_point = circle.get(pointList[pointer]);
    ctx.lineTo(t_point.x, t_point.y);
    pointer++;
  }
  ctx.stroke();

  if (pointer >= pointList.length) {
    doLoop = false;
  }

  if(runButton.value == "Run") {
    ctx.clearRect(0, 0, 400, 400); // clear canvas
  }

  if (doLoop) {
    window.requestAnimationFrame(draw);
  }
}

const imageFile = document.getElementById("image-file");
const uploadButton = document.getElementById("upload");
const apiURL = window.location.origin + "/api/";

uploadButton.addEventListener("click", upload);

function upload(event) {
  event.preventDefault();
  uploadButton.toggleAttribute("disabled");
  const infoText = document.getElementById("info-text");
  infoText.removeAttribute("hidden");
  infoText.innerText = "Uploading the image. Please wait!";
  let uploadData = new FormData();
  uploadData.append("image", new Blob([imageFile.files[0]]));
  fetch(apiURL, {
    method: "POST",
    mode: "same-origin",
    body: uploadData
  }).then((res) => {
    if (res.status == 200) {
      res.text().then((result) => {
        const id = result.split(" ")[4];
        window.location.assign(`art.html?id=${id}`);
      })
    } else {
      uploadButton.toggleAttribute("disabled");
      res.text().then(s => {
        infoText.innerText = s;
      });
    }
  });
}
