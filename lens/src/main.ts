import "bulma/css/bulma.min.css";

import {createApp, ref} from "vue";
import App from "./App.vue";
import router from "./router";

const IS_PRODUCTION = import.meta.env.MODE === "production";
console.log(`Running in '${IS_PRODUCTION ? "production" : "development"}' mode.`);
export const API_URL = IS_PRODUCTION ? "" : "http://localhost:8080";
export const DEMO_USER_ID: string = import.meta.env.VITE_DEMO_USER_ID;

console.log(DEMO_USER_ID);

export const userID = ref("");

const app = createApp(App);

app.use(router);

app.mount("#app");
