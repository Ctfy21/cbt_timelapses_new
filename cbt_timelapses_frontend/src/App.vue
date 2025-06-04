<script setup>

import { useOrderStore } from "@/orderStore";
import {ref} from "vue";

const roomValue = ref('')
const cameraValue = ref('')
const startDateValue = ref('')
const endDateValue = ref('')
const showSettings = ref(false)
const settingsId = ref('')
const settingsName = ref('')

const protocol = window.location.protocol
const port = window.location.port
const hostname = window.location.hostname

const store = useOrderStore();

store.bindEvents()

const openSettings = () => {
  showSettings.value = true
}

const closeSettings = () => {
  showSettings.value = false
  settingsId.value = ''
  settingsName.value = ''
}

const saveSettings = () => {
  // Here you can add logic to save the settings
  closeSettings()
}

</script>

<template>
  <div class="container-fluid">
    <!-- Header -->
    <header class="bg-body text-white p-3 mb-4">
      <div class="d-flex justify-content-between align-items-center">
        <h1 class="h3 mb-0">Generate Timelapse</h1>
        <button class="btn btn-light" @click="openSettings">
          <i class="bi bi-gear"></i> Настройки
        </button>
      </div>
    </header>

    <!-- Settings Modal -->
    <div class="modal fade" :class="{ 'show d-block': showSettings }" tabindex="-1" v-if="showSettings">
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title">Настройки</h5>
            <button type="button" class="btn-close" @click="closeSettings"></button>
          </div>
          <div class="modal-body">
            <div class="mb-3">
              <label class="form-label">Id</label>
              <input type="text" class="form-control" v-model="settingsId">
            </div>
            <div class="mb-3">
              <label class="form-label">Введите название</label>
              <input type="text" class="form-control" v-model="settingsName">
            </div>
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-primary" @click="saveSettings">Сохранить</button>
          </div>
        </div>
      </div>
    </div>
    <div class="modal-backdrop fade" :class="{ 'show': showSettings }" v-if="showSettings"></div>

    <!-- Main Content -->
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
              <a :href="protocol + '//'+ hostname + ':' + port + '/download/' + order.room + '/' + order.camera + '/timelapses/' + 'output_' + order.startDate.toISOString().slice(0,10) + '_00-00-00_to_' + order.endDate.toISOString().slice(0,10) + '_00-00-00.mp4'" v-if="order.status === 200" download>
                <BIconDownload> Скачать </BIconDownload>
              </a>
            </div>
          </li>
        </ul>

      </div>
    </div>
  </div>

</template>


<style scoped>

</style>
