<script setup lang="ts">
import AnimatedImage from "@/components/AnimatedImage.vue";
import {ref} from "vue";
import {useRouter} from "vue-router";
import {API_URL, userID} from "@/main";

const uploadButton = ref<HTMLInputElement | null>(null);
const router = useRouter();

/**
 * Show the error message that the Response, res, has.
 * @param res the body will be consumed on this Response.
 */
async function showErrorMessage(res: Response): Promise<void> {
  let alertString = "Der er sket en fejl! Prøv igen senere.";
  let errorMsg = await res.text();
  if (errorMsg.length > 0) alertString += ` Se fejlbeskeden: ${errorMsg}`;
  alert(alertString);
}

/**
 * goToRouteGiven expects a Response, res, that has the id to navigate to as body text.
 * @param res will consume the body of this response.
 */
async function goToRouteGivenIn(res: Response): Promise<void> {
  let id = await res.text();
  localStorage.setItem("userID", id);
  userID.value = id;
  await router.push(`/user/${id}`);
}

/**
 * Uploads the image in the uploadButton.value. Checks if there is any file before trying. Shows an error message,
 * if the server sends an error back. It will automatically navigate to the ID returned.
 */
async function uploadImage() {
  if (!uploadButton.value) return;
  if (!uploadButton.value.files) return;
  if (uploadButton.value?.files.length == 0) {
    alert("Ingen filer er valgt.");
    return;
  }
  let file = uploadButton.value.files[0];
  let formData = new FormData();
  let headers = userID.value ? {"Authorization": `Basic ${userID.value}`} : new Headers();
  formData.set("image", file);
  let res = await fetch(API_URL + "/api/upload", {
    method: "POST",
    body: formData,
    headers: headers,
  });
  if (!res.ok) {
    await showErrorMessage(res);
    return;
  }
  await goToRouteGivenIn(res);
}
</script>

<template>
  <main>
    <section class="hero is-primary">
      <div class="container">
        <div class="hero-body">
          <div class="columns">
            <div class="column">
              <p class="title">Prøv tråd kunst</p>
              <p class="subtitle mb-2">Simpelt & gratis</p>
              <p>Du kan nemt afprøve hvordan et billede vil se ud lavet i tråde. Hvis det er noget, som du kan lide, så
                kan
                du prøve at lave det selv i virkeligheden</p>
            </div>
            <div class="column is-justify-content-center is-flex">
              <AnimatedImage/>
            </div>
          </div>
        </div>
      </div>
    </section>
    <section class="hero is-warning">
      <div class="container">
        <div class="hero-body">
          <p class="title">Kom i gang</p>
          <p class="subtitle mb-2">Det simpelt!</p>
          <div class="columns">
            <div class="column">
              <p>Billedet du uploader må være <kbd>.jpg</kbd> eller <kbd>.png</kbd>. Det må have en maksimal størrelse
                på 10 MB og det skal være indenfor de to størrelser 400 x 400 og 1024 x 1024.</p>
            </div>

            <div class="column is-flex is-align-items-center is-justify-content-center is-gap-2">
              <div class="file is-link mb-0">
                <label class="file-label">
                  <input class="file-input" type="file" accept="image/png, image/jpeg" ref="uploadButton"/>
                  <span class="file-cta">
                <span class="file-label">Vælg en fil...</span>
              </span>
                </label>
              </div>
              <button class="button is-warning is-light" @click="uploadImage">Upload!</button>
            </div>
          </div>
        </div>
      </div>
    </section>
  </main>
</template>
