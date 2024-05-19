<script setup>

import { useOrderStore } from "@/orderStore";
import {ref, toRaw} from "vue";
import axios from "axios";


const roomValue = ref('')
const cameraValue = ref('')
const startDateValue = ref('')
const endDateValue = ref('')

const store = useOrderStore();

store.bindEvents()

function downloadImage(order){
  const existingOrderIndex = store.orders.findIndex((temp) => {
    return temp.id === order.id;
  });
  const url = `http://192.168.42.119:5000/download/${order.room}/${order.camera}/timelapses/output_${order.startDate.toISOString().slice(0,10)}_00-00-00_to_${order.endDate.toISOString().slice(0,10)}_00-00-00.mp4`
  axios({
    url: url,
    method: 'GET',
    responseType: 'blob',
    onDownloadProgress: function (progressEvent){
      store.orders[existingOrderIndex].downloaderValue = ((progressEvent.loaded / progressEvent.total) * 100).toFixed() + "%"
    }
  }).then((response) => {
    store.orders[existingOrderIndex].downloaderValue = ''
    // create file link in browser's memory
    const href = URL.createObjectURL(response.data);
    // create "a" HTML element with href to file & click
    const link = document.createElement('a');
    link.href = href;
    link.setAttribute('download', `timelapse_${order.room}_${order.camera}_${order.startDate.toISOString().slice(0,10)}_${order.endDate.toISOString().slice(0,10)}`); //or any other extension
    document.body.appendChild(link);
    link.click();
    // clean up "a" element & remove ObjectURL
    document.body.removeChild(link);
    URL.revokeObjectURL(href);
  });
}

</script>

<template>
  <div class="row app justify-content-center p-2">
    <form @submit.prevent="store.addNewOrder(roomValue, cameraValue, startDateValue, endDateValue)" class="row w-50 gy-4">
      <div class="input-group">
        <span class="input-group-text">Комната</span>
        <select class="form-select" v-model="roomValue" aria-label="Пример выбора по умолчанию">
          <option v-for="(cameras, nameRoom) in store.folders" :key="nameRoom" :value="nameRoom">{{ nameRoom }}</option>
        </select>
      </div>
      <div class="input-group">
        <span class="input-group-text">Камера</span>
        <select class="form-select" v-model="cameraValue" aria-label="Пример выбора по умолчанию">
          <option v-for="(cameras, ) in store.folders[roomValue]" :key="cameras" :value="cameras">{{ cameras }}</option>
        </select>
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
          <div v-if="order.downloaderValue === ''" class="p-1 d-flex col border align-items-center justify-content-center">
            <button v-if="order.status === 200">
              <BIconDownload @click="downloadImage(toRaw(order))"> Скачать </BIconDownload>
            </button>
          </div>
          <div v-else class="p-1 d-flex col border align-items-center justify-content-center">
            <label>{{ order.downloaderValue }}</label>
          </div>
        </li>
      </ul>

    </div>
  </div>

</template>.


<style scoped>

</style>