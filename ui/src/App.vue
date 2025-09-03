<!-- src/App.vue -->
<template>
  <div class="app">
    <div class="main-container">
      <div class="left-section">
        <div class="upper-left-pane">
          <FullCalendarPanel />
        </div>
        <div class="lower-left-pane">
          <AccountBalanceChart />
        </div>
      </div>
      <div class="right-section">
        <EventFormPanel />
      </div>
    </div>
    <ToastNotification />
  </div>
</template>

<script setup lang="ts">
import { onMounted } from "vue";
import { useAppStore } from "@/stores/appStore";
import FullCalendarPanel from "@/components/FullCalendarPanel.vue";
import EventFormPanel from "@/components/EventFormPanel.vue";
import AccountBalanceChart from "@/components/AccountBalanceChart.vue";
import ToastNotification from "@/components/ToastNotification.vue";

const store = useAppStore();

onMounted(() => {
  store.fetchAppState();
});
</script>

<style scoped>
.app {
  height: 100vh;
  overflow: hidden;
}

.main-container {
  display: flex;
  height: 100%;
}

.left-section {
  width: 70%;
  display: flex;
  flex-direction: column;
}

.upper-left-pane {
  height: 60%;
  border-right: 1px solid #ddd;
  border-bottom: 1px solid #ddd;
}

.lower-left-pane {
  height: 40%;
  border-right: 1px solid #ddd;
}

.right-section {
  width: 30%;
  padding: 20px;
  background-color: #f9f9f9;
}
</style>
