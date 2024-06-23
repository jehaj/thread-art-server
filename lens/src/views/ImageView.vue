<script setup lang="ts">

import {useRoute} from "vue-router";

const route = useRoute();
import {API_URL} from "@/main";
import AnimatedImage from "@/components/AnimatedImage.vue";
import {ref} from "vue";

let pointIndices: number[];
let numberOfPoints: number;

let pointsFetched = ref(false);

let downloadImage = ref<HTMLAnchorElement>();
let downloadText = ref<HTMLAnchorElement>();

async function getPoints() {
  let pointsResponse = await fetch(`${API_URL}/api/${route.params.id}/points`);
  let pointsJson = await pointsResponse.json();
  pointIndices = pointsJson.PointIndex;
  numberOfPoints = pointsJson.NumberOfPoints;
  pointsFetched.value = true;

  let textBlob = new Blob([numberOfPoints + "\n", pointIndices.join(", ")], {type: "text/plain"});
  if (downloadText.value === undefined) return;
  downloadText.value.setAttribute("href", URL.createObjectURL(textBlob));
}

getPoints();
</script>

<template>
  <section class="section">
    <div class="container">
      <h1 class="is-size-1">Sådan!</h1>
      <p class="is-size-5">Dette er siden tilhørende billedet med id <span
          class="has-text-warning has-text-weight-bold">{{ $route.params.id }}</span>.</p>
      <p>Tryk nedenfor for at hente hhv. billedet og punkterne.</p>
      <div class="buttons has-addons mt-3">
        <!-- download attributen virker kun når domænet er det samme eller der er tale om blob. -->
        <a :href="`${API_URL}/api/${route.params.id}/out.png`" download class="button is-link">Billede (png)</a>
        <a ref="downloadText" download class="button is-link">Punkter (txt)</a>
      </div>
      <div class="columns mt-3 is-justify-content-space-between">
        <div class="column is-one-third">
          <img :src="`${API_URL}/api/${route.params.id}/in.png`" alt="">
        </div>
        <div class="column is-one-third">
          <img :src="`${API_URL}/api/${route.params.id}/out.png`" alt="">
        </div>
        <div class="column is-one-third">
          <AnimatedImage v-if="pointsFetched" :numberOfPoints="numberOfPoints" :pointIndices="pointIndices"/>
        </div>
      </div>
    </div>
  </section>
</template>

<style>
.column.is-one-third {
  max-width: 350px;
}
</style>
