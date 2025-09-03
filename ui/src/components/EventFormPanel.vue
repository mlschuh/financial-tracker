<!-- src/components/EventFormPanel.vue -->
<template>
  <div class="event-form-panel">
    <div class="form-header">
      <h3 v-if="store.isCreatingNewEvent">Create New Event</h3>
      <h3 v-else>Edit Event</h3>
      <button
        v-if="!store.isCreatingNewEvent"
        @click="store.startCreatingNewEvent()"
        class="btn-secondary"
      >
        Add New Event
      </button>
    </div>

    <form @submit.prevent="handleSubmit" class="event-form">
      <div class="form-group">
        <label for="name">Name:</label>
        <input
          id="name"
          v-model="formData.name"
          type="text"
          required
          class="form-input"
        />
      </div>

      <div class="form-group">
        <label for="category">Category:</label>
        <input
          id="category"
          v-model="formData.category"
          type="text"
          class="form-input"
        />
      </div>

      <div class="form-group">
        <label for="account">Account:</label>
        <select
          id="account"
          v-model="formData.account"
          required
          class="form-input"
        >
          <option value="">Select Account</option>
          <option
            v-for="account in store.accounts"
            :key="account.id"
            :value="account.id"
          >
            {{ account.name }}
          </option>
        </select>
      </div>

      <div class="form-group">
        <label for="amount">Amount ($):</label>
        <input
          id="amount"
          v-model.number="amountInDollars"
          type="number"
          step="0.01"
          min="0"
          required
          class="form-input"
        />
      </div>

      <div class="form-group">
        <label for="type">Type:</label>
        <div class="radio-group">
          <label class="radio-label">
            <input v-model="formData.type" type="radio" value="income" />
            Income
          </label>
          <label class="radio-label">
            <input v-model="formData.type" type="radio" value="expense" />
            Expense
          </label>
        </div>
      </div>

      <div class="form-group">
        <label for="start">Start Date:</label>
        <input
          id="start"
          v-model="startDate"
          type="datetime-local"
          required
          class="form-input"
        />
      </div>

      <div class="form-group">
        <label for="rrule">Recurrence Rule:</label>
        <input
          id="rrule"
          v-model="formData.rrule"
          type="text"
          placeholder="e.g., FREQ=MONTHLY;BYMONTHDAY=1"
          class="form-input"
        />
        <small class="form-help">Leave empty for single occurrence</small>
      </div>

      <div
        v-if="!store.isCreatingNewEvent && store.editingEvent?.exceptions"
        class="exceptions-section"
      >
        <h4>Exceptions:</h4>
        <div
          v-if="Object.keys(store.editingEvent.exceptions).length === 0"
          class="no-exceptions"
        >
          No exceptions defined
        </div>
        <div v-else class="exceptions-list">
          <div
            v-for="[date, exception] in Object.entries(
              store.editingEvent.exceptions
            )"
            :key="date"
            class="exception-item"
          >
            <strong>{{ date }}:</strong> {{ exception.type }}
            <span v-if="exception.amount">
              - ${{ (exception.amount / 100).toFixed(2) }}</span
            >
          </div>
        </div>
      </div>

      <div class="form-actions">
        <button
          v-if="!store.isCreatingNewEvent && store.selectedEventOccurrenceId"
          @click="handleIgnoreSelectedDate"
          type="button"
          class="btn-warning"
        >
          Ignore Selected Date
        </button>
        <button
          v-if="!store.isCreatingNewEvent && store.editingEvent"
          @click="handleDeleteEvent"
          type="button"
          class="btn-danger"
        >
          Delete Event
        </button>

        <button
          v-if="store.isCreatingNewEvent"
          type="submit"
          class="btn-primary"
        >
          Add Event
        </button>

        <button v-else type="submit" class="btn-primary">Update Event</button>
      </div>
    </form>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from "vue";
import { useAppStore } from "@/stores/appStore";
import type { EventCreate } from "@/types/api";

const store = useAppStore();

const formData = ref<EventCreate>({
  name: "",
  category: "",
  account: "",
  amount: 0,
  start: "",
  type: "expense",
  rrule: "",
  exceptions: {},
});

const amountInDollars = computed({
  get: () => formData.value.amount / 100,
  set: (value: number) => {
    formData.value.amount = Math.round(value * 100);
  },
});

const startDate = computed({
  get: () => {
    if (!formData.value.start) return "";
    return new Date(formData.value.start).toISOString().slice(0, 16);
  },
  set: (value: string) => {
    formData.value.start = new Date(value).toISOString();
  },
});

// Watch for editing event changes
watch(
  () => store.editingEvent,
  (newEvent) => {
    if (newEvent && !store.isCreatingNewEvent) {
      formData.value = { ...newEvent };
    }
  },
  { immediate: true }
);

// Watch for creating new event
watch(
  () => store.isCreatingNewEvent,
  (isCreating) => {
    if (isCreating) {
      formData.value = {
        name: "",
        category: "",
        account: "",
        amount: 0,
        start: new Date().toISOString(),
        type: "expense",
        rrule: "",
        exceptions: {},
      };
    }
  }
);

const handleSubmit = async () => {
  if (store.isCreatingNewEvent) {
    await store.createEvent(formData.value);
  } else if (store.editingEvent) {
    await store.updateEvent(store.editingEvent.id, formData.value);
  }
};

const handleIgnoreSelectedDate = async () => {
  if (!store.selectedEventOccurrenceId || !store.editingEvent) return;

  // Extract date from occurrence ID
  const occurrenceIdParts = store.selectedEventOccurrenceId.split("-");
  if (occurrenceIdParts.length < 2) return;

  const dateStr = occurrenceIdParts.slice(1).join("-");
  const date = new Date(dateStr).toISOString().split("T")[0]; // Get YYYY-MM-DD format

  await store.addEventException(store.editingEvent.id, date);
};

const handleDeleteEvent = async () => {
  if (!store.editingEvent) return;

  if (
    confirm(`Are you sure you want to delete "${store.editingEvent.name}"?`)
  ) {
    await store.deleteEvent(store.editingEvent.id);
  }
};
</script>

<style scoped>
.event-form-panel {
  height: 100%;
  overflow-y: auto;
}

.form-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding-bottom: 10px;
  border-bottom: 1px solid #ddd;
}

.form-header h3 {
  margin: 0;
}

.event-form {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.form-group {
  display: flex;
  flex-direction: column;
}

.form-group label {
  margin-bottom: 5px;
  font-weight: bold;
  color: #333;
}

.form-input {
  padding: 8px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
}

.form-input:focus {
  outline: none;
  border-color: #3498db;
  box-shadow: 0 0 0 2px rgba(52, 152, 219, 0.2);
}

.radio-group {
  display: flex;
  gap: 15px;
}

.radio-label {
  display: flex;
  align-items: center;
  gap: 5px;
  font-weight: normal;
}

.form-help {
  margin-top: 5px;
  color: #666;
  font-size: 12px;
}

.exceptions-section {
  margin: 20px 0;
  padding: 15px;
  background-color: #f8f9fa;
  border-radius: 4px;
}

.exceptions-section h4 {
  margin: 0 0 10px 0;
  color: #333;
}

.no-exceptions {
  color: #666;
  font-style: italic;
}

.exceptions-list {
  display: flex;
  flex-direction: column;
  gap: 5px;
}

.exception-item {
  padding: 5px 0;
  color: #555;
}

.form-actions {
  display: flex;
  flex-direction: column;
  gap: 10px;
  margin-top: 20px;
}

.btn-primary,
.btn-secondary,
.btn-warning {
  padding: 10px 20px;
  border: none;
  border-radius: 4px;
  font-size: 14px;
  cursor: pointer;
  transition: background-color 0.2s;
}

.btn-primary {
  background-color: #3498db;
  color: white;
}

.btn-primary:hover {
  background-color: #2980b9;
}

.btn-secondary {
  background-color: #95a5a6;
  color: white;
}

.btn-secondary:hover {
  background-color: #7f8c8d;
}

.btn-warning {
  background-color: #f39c12;
  color: white;
}

.btn-warning:hover {
  background-color: #d68910;
}

.btn-danger {
  background-color: #e74c3c;
  color: white;
}

.btn-danger:hover {
  background-color: #c0392b;
}
</style>
