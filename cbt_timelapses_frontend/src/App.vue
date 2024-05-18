<script setup>

import { useOrderStore } from "@/orderStore";
import { ref } from "vue";
import {downloadImage} from "@/downloader";

const store = useOrderStore();

const roomValue = ref('')
const cameraValue = ref('')
const startDateValue = ref('')
const endDateValue = ref('')

store.bindEvents()

</script>

<template>
  <div class="row app justify-content-center p-2">
    <form @submit.prevent="store.addNewOrder(roomValue, cameraValue, startDateValue, endDateValue)" class="row w-50 gy-4">
      <div class="input-group">
        <span class="input-group-text">Комната</span>
        <input type="text" class="form-control input" v-model="roomValue" placeholder="Название комнаты">
      </div>
      <div class="input-group">
        <span class="input-group-text">Камера</span>
        <input type="text" class="form-control input" v-model="cameraValue" placeholder="Название камеры">
      </div>
      <div class="input-group">
        <span class="input-group-text">Дата начала</span>
        <input type="date" class="form-control" v-model="startDateValue" >
      </div>
      <div class="input-group">
        <span class="input-group-text">Дата конца</span>
        <input type="date" class="form-control" v-model="endDateValue">
      </div>
      <div class="d-flex justify-content-center">
        <button type="submit" class="btn btn-primary" > Создать </button>
      </div>
    </form>

    <hr class="row mt-5">

    <div class="row w-75">
        <div class="rounded row gy-4">
          <div class="p-1 d-flex col align-items-center justify-content-center">
            <h5 class="h5">Название комнаты</h5>
          </div>
          <div class="p-1 d-flex col align-items-center justify-content-center">
            <h5 class="h5">Название камеры</h5>
          </div>
          <div class="p-1 d-flex col align-items-center justify-content-center">
            <h5 class="h5">Дата начала</h5>
          </div>
          <div class="p-1 d-flex col align-items-center justify-content-center">
            <h5 class="h5">Дата конца</h5>
          </div>
          <div class="p-1 d-flex col align-items-center justify-content-center">
            <h5 class="h5">Статус</h5>
          </div>
          <div class="p-1 d-flex col align-items-center justify-content-center">
            <h5 class="h5">Скачать</h5>
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
            <button v-if="order.status === 200">
              <BIconDownload @click="downloadImage(order.room, order.camera, order.startDate, order.endDate)"> Скачать </BIconDownload>
            </button>
          </div>
        </li>
      </ul>

    </div>
  </div>

</template>.


<style scoped>

</style>