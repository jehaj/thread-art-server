import "bulma/css/bulma.min.css";

import {createApp} from "vue";
import App from "./App.vue";
import router from "./router";

const BASE_URL = import.meta.env.BASE_URL;
const BASE_URL_EMPTY = BASE_URL.length < 3;
export const API_URL = BASE_URL_EMPTY ? "http://localhost:8080" : BASE_URL;

const app = createApp(App);

app.use(router);

app.mount("#app");
