import "bulma/css/bulma.min.css";

import {createApp, ref} from "vue";
import App from "./App.vue";
import router from "./router";

const BASE_URL = import.meta.env.BASE_URL;
const BASE_URL_EMPTY = BASE_URL.length < 3;
export const API_URL = BASE_URL_EMPTY ? "http://localhost:8080" : BASE_URL;
export const DEMO_USER_ID: string = import.meta.env.VITE_DEMO_USER_ID;

console.log(DEMO_USER_ID);

export const userID = ref("");

const app = createApp(App);

app.use(router);

app.mount("#app");
