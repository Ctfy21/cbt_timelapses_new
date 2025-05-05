<script setup>

import { useOrderStore } from "@/orderStore";
import {ref} from "vue";

const roomValue = ref('')
const cameraValue = ref('')
const startDateValue = ref('')
const endDateValue = ref('')
  
const ipAddress = process.env.VUE_APP_PATH_START
  
const store = useOrderStore();

store.bindEvents()

</script>

<template>
  <div class="row app justify-content-center p-2">
    <form @submit.prevent="store.addNewOrder(roomValue, cameraValue, startDateValue, endDateValue)" class="row w-50 gy-4">
      <div class="input-group">
        <span class="input-group-text">Комната</span>
        <select class="form-select" v-model="roomValue" aria-label="Пример выбора по умолчанию" required>
          <option v-for="(cameras, nameRoom) in store.folders" :key="nameRoom" :value="nameRoom">{{ nameRoom }}</option>
        </select>
      </div>
      <div class="input-group">
        <span class="input-group-text">Камера</span>
        <select class="form-select" v-model="cameraValue" aria-label="Пример выбора по умолчанию" required>
          <option v-for="(cameras, ) in store.folders[roomValue]" :key="cameras" :value="cameras">{{ cameras }}</option>
        </select>
      </div>
      <div class="input-group">
        <span class="input-group-text">Дата начала</span>
        <input type="date" class="form-control" v-model="startDateValue" required>
      </div>
      <div class="input-group">
        <span class="input-group-text">Дата конца</span>
        <input type="date" class="form-control" v-model="endDateValue" required>
      </div>
      <div class="d-flex justify-content-center">
        <button type="submit" class="btn btn-primary" > Создать </button>
      </div>
    </form>
    <div class="d-flex justify-content-center mt-2">
        <p class="text-danger" v-if="store.error != ''">{{ store.error }}</p>
    </div>


    <hr class="row mt-5">

    <div class="row w-75">
        <div class="row gy-4">
          <div class="p-1 d-flex col align-items-center justify-content-center">
            <h5 class="h5 text-center">Название комнаты</h5>
          </div>
          <div class="p-1 d-flex col align-items-center justify-content-center">
            <h5 class="h5 text-center">Название камеры</h5>
          </div>
          <div class="p-1 d-flex col align-items-center justify-content-center">
            <h5 class="h5 text-center">Дата начала</h5>
          </div>
          <div class="p-1 d-flex col align-items-center justify-content-center">
            <h5 class="h5 text-center">Дата конца</h5>
          </div>
          <div class="p-1 d-flex col align-items-center justify-content-center">
            <h5 class="h5 text-center">Статус</h5>
          </div>
          <div class="p-1 d-flex col align-items-center justify-content-center">
            <h5 class="h5 text-center">Скачать</h5>
          </div>
        </div>
      <ul>
        <li class="rounded row gy-4 m-1" v-for="order in store.orders" :key="order.id">
          <div class="p-1 d-flex col border align-items-center justify-content-center">
            <label class="">{{ order.room }}</label>
          </div>
          <div class="p-1 d-flex col border align-items-center justify-content-center">
            <label class="">{{ order.camera }}</label>
          </div>
          <div class="p-1 d-flex col border align-items-center justify-content-center">
            <label>{{ order.startDate.getDate() }}.{{ order.startDate.getMonth()+1 }}.{{ order.startDate.getFullYear() }}</label>
          </div>
          <div class="p-1 d-flex col border align-items-center justify-content-center">
            <label>{{ order.endDate.getDate() }}.{{ order.endDate.getMonth()+1 }}.{{ order.endDate.getFullYear() }}</label>
          </div>
          <div class="p-1 d-flex col border align-items-center justify-content-center">
            <div class="spinner-border text-warning" v-if="order.status === 300"></div>
            <BIconCheckCircleFill class="text-success" v-else-if="order.status === 200"></BIconCheckCircleFill>
            <BIconXCircleFill class="text-danger" v-else></BIconXCircleFill>
          </div>
          <div class="p-1 d-flex col border align-items-center justify-content-center">
            <a :href="'http://'+ ipAddress +':5000/download/' + order.room + '/' + order.camera + '/timelapses/' + 'output_' + order.startDate.toISOString().slice(0,10) + '_00-00-00_to_' + order.endDate.toISOString().slice(0,10) + '_00-00-00.mp4'" v-if="order.status === 200" download>
              <BIconDownload> Скачать </BIconDownload>
            </a>
          </div>
        </li>
      </ul>

    </div>
  </div>

</template>.


<style scoped>

</style>
