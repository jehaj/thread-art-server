<script setup lang="ts">
import {RouterLink} from "vue-router";
import {API_URL} from "@/main";
import type {Image} from "@/models";

let props = defineProps<{ image: Image }>();
let id = props.image.ID;
</script>

<template>
  <div class="cell">
    <RouterLink :to="`/image/${id}`">
      <div class="card">
        <div class="card-image">
          <figure class="image is-1by1">
            <img
                :src="`${API_URL}/api/${id}/in.png`"
                alt="Placeholder image"
            />
            <img
                v-if="props.image.Finished"
                class="behind"
                :src="`${API_URL}/api/${id}/out.png`"
                alt="Placeholder image"
            />
          </figure>
        </div>
        <div class="card-footer"
             :class="{'has-background-link-light': props.image.Finished, 'has-background-warning-light': !props.image.Finished}">
          <a class="card-footer-item is-primary" v-if="props.image.Finished">Gå til</a>
          <a class="card-footer-item " v-else>Kom igen senere</a>
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

.behind {
  position: absolute;
  z-index: 1;
  transition: opacity 0.3s;
  opacity: 0;
}

.card:hover {
  .behind {
    opacity: 1;
    transition: opacity 0.3s;
  }
}
</style>