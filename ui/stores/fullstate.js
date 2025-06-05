import { defineStore } from "pinia";
import axios from "axios";
import { ref } from "vue";

export const useFullState = defineStore("fullState", () => {
  const eventOccurances = ref([]);
  const accountBalances = ref([]);
  const events = ref([]);
  const formattedEvents = ref([]);
  const accounts = ref([]);
  const selectedEventOccurance = ref({});
  const addEventDate = ref("");
  const selectedDateRange = ref({});

  async function updateFromServer() {
    try {
      const response = await axios.get("/api/state");
      this.eventOccurances = response.data.eventOccurances;
      this.accountBalances = response.data.accountBalances;
      this.events = response.data.events;
      this.accounts = response.data.accounts;
    } catch (err) {
      this.error = err.message;
    }
  }

  return {
    eventOccurances,
    accountBalances,
    events,
    formattedEvents,
    accounts,
    selectedEventOccurance,
    addEventDate,
    selectedDateRange,
    updateFromServer,
  };
});
