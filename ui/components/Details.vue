<template>
  <div v-if="fullStateStore.addEventDate != ''"> 
    Add Modal
  <AddModal 
      :is-open="isModalOpen"
      :date="fullStateStore.addEventDate.dateStr"
      @close="closeModal"
      @submit="handleAddEvent"
    />
  </div>
  <div class="details-container">
    {{ fullStateStore.selectedEventOccurance }}
  </div>
</template>

<script setup>
import { useFullState } from "../stores/fullstate.js";
import AddModal from "./AddModal.vue"

const fullStateStore = useFullState();
</script>

<style scoped>
.details-container {
  position: relative;
  height: 100%;
  width: 100%;
}
</style>
