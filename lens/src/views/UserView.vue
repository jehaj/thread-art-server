<script setup lang="ts">
import {ref} from "vue";
import ImageBox from "@/components/ImageBox.vue";
import {API_URL} from "@/main";
import type {Image} from "@/models";
import {useRoute} from "vue-router";

let images = ref<Image[]>();

async function loadImages(id: string) {
  console.log(`Loading images from user ${id}.`);
  let response = await fetch(API_URL + `/api/user/${id}`);
  let json = await response.json();
  images.value = json.Images;
}

const route = useRoute();

loadImages(route.params.id as string);
</script>

<template>
  <section class="section">
    <div class="container">
      <h1>Hej! Dit bruger ID er <span class="has-text-danger has-text-weight-bold">{{ $route.params.id }}</span>.</h1>
      <div class="grid">
        <ImageBox v-for="image in images" :key="image.ID" :image="image"/>
      </div>
    </div>
  </section>

</template>

<style>

</style>
