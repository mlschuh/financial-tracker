<!-- src/components/ToastNotification.vue -->
<template>
  <Transition name="toast">
    <div v-if="store.showToast" class="toast">
      {{ store.toastMessage }}
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
  background-color: #e74c3c;
  color: white;
  padding: 15px 20px;
  border-radius: 4px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  z-index: 1000;
  max-width: 400px;
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
