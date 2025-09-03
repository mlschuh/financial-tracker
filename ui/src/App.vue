<!-- src/App.vue -->
<template>
  <div class="app">
    <div class="main-container">
      <div class="left-section">
        <ResizableSplitter
          :min-upper-height="200"
          :min-lower-height="150"
          :initial-upper-ratio="0.6"
        >
          <template #upper>
            <div class="upper-left-pane">
              <FullCalendarPanel />
            </div>
          </template>
          <template #lower>
            <div class="lower-left-pane">
              <AccountBalanceChart />
            </div>
          </template>
        </ResizableSplitter>
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
import ResizableSplitter from "@/components/ResizableSplitter.vue";

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
  border-right: 1px solid #ddd;
  display: flex;
  flex-direction: column;
  height: 100%;
}

.upper-left-pane {
  height: 100%;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.lower-left-pane {
  height: 100%;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.right-section {
  width: 30%;
  padding: 20px;
  background-color: #f9f9f9;
  overflow-y: auto;
}
</style>
