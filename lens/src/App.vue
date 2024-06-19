<script setup lang="ts">
import {RouterLink, RouterView} from "vue-router";
import {userID} from "@/main";

userID.value = localStorage.getItem("userID") || "";

</script>

<template>
  <header>
    <nav class="navbar" role="navigation" aria-label="main navigation">
      <div class="container">
        <div class="navbar-brand">
          <RouterLink class="navbar-item" to="/">Hjem</RouterLink>
          <RouterLink class="navbar-item" :to="`/user/${userID}`" v-if="userID !== ''">Bruger</RouterLink>
        </div>
      </div>
    </nav>
  </header>

  <RouterView/>
</template>
<style scoped>
.router-link-active {
  background-color: #70ffb1;
}
</style>
