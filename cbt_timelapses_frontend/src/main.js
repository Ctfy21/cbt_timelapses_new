import 'bootstrap/dist/css/bootstrap.css'
import 'bootstrap-vue/dist/bootstrap-vue.css'

import { BootstrapIconsPlugin } from "bootstrap-icons-vue";

import VCalendar from 'v-calendar';
import 'v-calendar/style.css';

import { createApp } from 'vue'
import App from './App.vue'

import { createPinia } from "pinia";

const pinia = createPinia()
const app = createApp(App)

app.use(BootstrapIconsPlugin)
app.use(VCalendar, {})
app.use(pinia)
app.mount('#app')
