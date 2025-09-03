<!-- src/components/ToastNotification.vue -->
<template>
  <Transition name="toast">
    <div v-if="store.showToast" :class="['toast', `toast-${store.toastType}`]">
      <div class="toast-icon">
        <span v-if="store.toastType === 'success'">✓</span>
        <span v-else-if="store.toastType === 'error'">✕</span>
        <span v-else-if="store.toastType === 'warning'">⚠</span>
        <span v-else-if="store.toastType === 'info'">ℹ</span>
      </div>
      <div class="toast-message">
        {{ store.toastMessage }}
      </div>
    </div>
  </Transition>
</template>

<script setup lang="ts">
import { useAppStore } from "@/stores/appStore";

const store = useAppStore();
</script>

<style scoped>
.toast {
  position: fixed;
  top: 20px;
  right: 20px;
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 15px 20px;
  border-radius: 6px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  z-index: 1000;
  max-width: 400px;
  font-weight: 500;
}

.toast-success {
  background-color: #27ae60;
  color: white;
}

.toast-error {
  background-color: #e74c3c;
  color: white;
}

.toast-warning {
  background-color: #f39c12;
  color: white;
}

.toast-info {
  background-color: #3498db;
  color: white;
}

.toast-icon {
  font-size: 16px;
  font-weight: bold;
  min-width: 20px;
  text-align: center;
}

.toast-message {
  flex: 1;
}

.toast-enter-active,
.toast-leave-active {
  transition: all 0.3s ease;
}

.toast-enter-from {
  opacity: 0;
  transform: translateX(100%);
}

.toast-leave-to {
  opacity: 0;
  transform: translateX(100%);
}
</style>
