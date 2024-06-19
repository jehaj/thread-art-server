<script setup lang="ts">

import {useRoute} from "vue-router";

const route = useRoute();
import {API_URL} from "@/main";
import AnimatedImage from "@/components/AnimatedImage.vue";
import {ref} from "vue";

let pointIndices: number[];
let numberOfPoints: number;

let pointsFetched = ref(false);

async function getPoints() {
  let pointsResponse = await fetch(`${API_URL}/api/${route.params.id}/points`);
  let pointsJson = await pointsResponse.json();
  pointIndices = pointsJson.PointIndex;
  numberOfPoints = pointsJson.NumberOfPoints;
  pointsFetched.value = true;
}

getPoints();
</script>

<template>
  <section class="section">
    <div class="container">
      <h1 class="is-size-1">Sådan!</h1>
      <p class="is-size-5">Dette er siden tilhørende billedet med id <span
          class="has-text-warning has-text-weight-bold">{{ $route.params.id }}</span>.</p>
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
