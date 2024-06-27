<script setup lang="ts">
import {onMounted, ref} from "vue";

const props = defineProps<{
  pointIndices: number[],
  numberOfPoints: number
}>();

class Point {
  x: number;
  y: number;

  constructor(x: number, y: number) {
    this.x = x;
    this.y = y;
  }
}

function constrain(x: number, min: number, max: number): number {
  if (x < min) x = min;
  if (max < x) x = max;
  return x;
}

const canvasSize = 350;

const canvas = ref<HTMLCanvasElement>();

onMounted(() => {
  console.log("Done loading, now animating :-).");
  if (!canvas.value) return;
  const ctx = canvas.value.getContext("2d");

  if (!ctx) return;
  ctx.globalCompositeOperation = "destination-over";
  ctx.strokeStyle = "rgba(0, 0, 0, 0.1)";

  let circle = new Map();
  const circleSize = props.numberOfPoints;

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
  let pointer = 0;
  let lastTimestamp: number;

  ctx.clearRect(0, 0, 350, 350);

  function draw(timestamp: number): void {
    if (lastTimestamp === 0) {
      lastTimestamp = timestamp;
    }
    let timePassed = timestamp - lastTimestamp;
    lastTimestamp = timestamp;
    if (!ctx) return;
    ctx.beginPath();
    for (let i = 0; i < Math.ceil(timePassed / 2); i++) {
      if (pointer >= props.pointIndices.length) {
        break;
      }
      let t_point = circle.get(props.pointIndices[pointer]);
      ctx.lineTo(t_point.x, t_point.y);
      pointer++;
    }
    ctx.stroke();

    if (pointer < props.pointIndices.length) {
      window.requestAnimationFrame(draw);
    } else {
      console.log("Done animating :-D");
    }
  }

  window.requestAnimationFrame(draw);
});
</script>

<template>
  <canvas ref="canvas" width="350" height="350">
    An error has occurred. A <kbd>canvas</kbd> element with an example is
    supposed to be here.
  </canvas>
</template>

<style scoped>
canvas {
  width: 100%;
}
</style>