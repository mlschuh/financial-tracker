<!-- EventModal.vue -->
<template>
  <div class="modal-overlay">
    <div class="modal-content">
      <div class="bg-white p-6 rounded-lg shadow-xl w-full max-w-md">
        <h2 class="text-xl font-bold mb-4">Add New Event for {{ date }}</h2>

        <form @submit.prevent="handleSubmit" class="space-y-4">
          <div>
            <label class="block text-sm font-medium mb-1">Name</label>
            <input
              v-model="formData.name"
              type="text"
              class="w-full p-2 border rounded"
              placeholder="Event name"
            />
          </div>

          <div>
            <label class="block text-sm font-medium mb-1">Amount</label>
            <input
              v-model="formData.amount"
              type="number"
              class="w-full p-2 border rounded"
              placeholder="Enter amount"
            />
          </div>

          <div>
            <label class="block text-sm font-medium mb-1">Account</label>
            <input
              v-model="formData.account"
              type="text"
              class="w-full p-2 border rounded"
              placeholder="Account name"
            />
          </div>

          <div class="flex justify-end space-x-2 mt-6">
            <button
              type="button"
              @click="handleClose"
              class="px-4 py-2 border rounded hover:bg-gray-100"
            >
              Cancel
            </button>
            <button
              type="submit"
              :disabled="isSubmitting"
              class="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600 disabled:opacity-50"
            >
              {{ isSubmitting ? "Saving..." : "Save" }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useFullState } from "../stores/fullstate.js";
import { reactive, ref } from "vue";

const fullStateStore = useFullState();

const props = defineProps({
  isOpen: {
    type: Boolean,
    required: true,
  },
  date: {
    type: String,
    required: true,
  },
});

const emit = defineEmits(["close", "submit"]);

const isSubmitting = ref(false);
const formData = reactive({
  name: "",
  amount: "",
  account: "",
});

const resetForm = () => {
  formData.name = "";
  formData.amount = "";
  formData.account = "";
};

const handleClose = () => {
  resetForm();
  emit("close");
};

const handleSubmit = async () => {
  console.log(`Submitting ${formData}`);
  emit("submit");
  resetForm();
  emit("close");

  // try {
  //   isSubmitting.value = true;
  //   const response = await fetch("/events", {
  //     method: "POST",
  //     headers: {
  //       "Content-Type": "application/json",
  //     },
  //     body: JSON.stringify({
  //       date: props.date,
  //       ...formData,
  //     }),
  //   });

  //   if (!response.ok) {
  //     throw new Error("Failed to create event");
  //   }

  //   emit("submit");
  //   resetForm();
  //   emit("close");
  // } catch (error) {
  //   console.error("Error creating event:", error);
  //   // You might want to show an error message to the user here
  // } finally {
  //   isSubmitting.value = false;
  // }
};
</script>

<style scoped>
/* Modal container (background overlay) */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.5); /* Semi-transparent background */
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000; /* Ensure it's above other elements */
  opacity: 0; /* Hidden by default */
  visibility: hidden; /* Hidden by default */
  transition: opacity 0.3s ease, visibility 0.3s ease; /* Smooth transition */
}

/* Show modal overlay */
.modal-overlay.active {
  opacity: 1;
  visibility: visible;
}

/* Modal content */
.modal-content {
  background-color: #fff;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
  max-width: 500px; /* Limit the modal width */
  width: 100%; /* Make it responsive */
  position: relative;
  overflow: hidden; /* Prevent overflow issues */
}

/* Close button */
.modal-close {
  position: absolute;
  top: 10px;
  right: 10px;
  background: none;
  border: none;
  font-size: 1.5rem;
  cursor: pointer;
}
</style>
