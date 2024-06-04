<script setup lang="ts">
import {RouterLink} from "vue-router";
import {API_URL} from "@/main";
import type {Image} from "@/models";

let props = defineProps<{image: Image}>();
let id = props.image.ID;
</script>

<template>
  <div class="cell">
    <RouterLink :to="`/image/${id}`">
      <div class="card">
        <div class="card-image">
          <figure class="image is-1by1">
            <img
                :src="`${API_URL}/api/${id}/image.png`"
                alt="Placeholder image"
            />
          </figure>
        </div>
        <div class="card-footer">
          <a class="card-footer-item">Gå til</a>
        </div>
      </div>
    </RouterLink>
  </div>
</template>

<style scoped>
.card-image {
  max-width: 250px;
  max-height: 250px;
}

.card {
  max-width: 250px;
}
</style>