// src/stores/appStore.ts
import { defineStore } from "pinia";
import { ref, computed } from "vue";
import { apiClient } from "@/utils/api";
import type {
  Account,
  Event,
  EventCreate,
  EventOccurrence,
  AccountBalance,
  AppStateResponse,
  Exception,
  ToastType,
} from "@/types/api";

export const useAppStore = defineStore("app", () => {
  // State
  const accounts = ref<Account[]>([]);
  const events = ref<Event[]>([]);
  const appState = ref<AppStateResponse>({
    eventOccurances: [],
    accountBalances: [],
    events: [],
    accounts: [],
  });
  const selectedEventOccurrenceId = ref<string>("");
  const editingEvent = ref<Event | null>(null);
  const isCreatingNewEvent = ref<boolean>(true);
  const chartDateRangeMonths = ref<number>(3);
  const toastMessage = ref<string>("");
  const toastType = ref<ToastType>("info"); // Add toast type
  const showToast = ref<boolean>(false);

  // Computed
  const filteredAccountBalances = computed(() => {
    const now = new Date();
    const startDate = new Date(
      now.getFullYear(),
      now.getMonth() - chartDateRangeMonths.value,
      now.getDate()
    );
    const endDate = new Date(
      now.getFullYear(),
      now.getMonth() + chartDateRangeMonths.value,
      now.getDate()
    );

    return appState.value.accountBalances.filter((balance) => {
      const balanceDate = new Date(balance.date);
      return balanceDate >= startDate && balanceDate <= endDate;
    });
  });

  // Actions
  const showToastMessage = (message: string, type: ToastType = "info") => {
    toastMessage.value = message;
    toastType.value = type;
    showToast.value = true;
    setTimeout(() => {
      showToast.value = false;
    }, 5000);
  };

  const handleApiError = (error: any) => {
    const errorMessage = error.response?.data?.error || "An error occurred";
    showToastMessage(errorMessage, "error");
  };

  const fetchAppState = async () => {
    try {
      const data = await apiClient.getAppState();
      appState.value = data;
      accounts.value = data.accounts;
      events.value = data.events;
    } catch (error) {
      handleApiError(error);
    }
  };

  const createAccount = async (accountData: {
    name: string;
    color: string;
  }) => {
    try {
      await apiClient.createAccount(accountData);
      await fetchAppState();
      showToastMessage("Account created successfully", "success");
    } catch (error) {
      handleApiError(error);
    }
  };

  const createEvent = async (eventData: EventCreate) => {
    try {
      await apiClient.createEvent(eventData);
      await fetchAppState();
      isCreatingNewEvent.value = false;
      showToastMessage("Event created successfully", "success");
    } catch (error) {
      handleApiError(error);
    }
  };

  const updateEvent = async (eventId: string, eventData: EventCreate) => {
    try {
      // First delete the existing event
      await apiClient.deleteEvent(eventId);

      // Then create a new event with the updated data
      await apiClient.createEvent(eventData);

      // Refresh the app state
      await fetchAppState();

      // Clear the selection since the old event ID no longer exists
      selectedEventOccurrenceId.value = "";
      editingEvent.value = null;

      showToastMessage("Event updated successfully", "success");
    } catch (error) {
      handleApiError(error);
      // If there was an error, refresh the state to ensure consistency
      await fetchAppState();
    }
  };

  const selectEventOccurrence = (occurrenceId: string) => {
    selectedEventOccurrenceId.value = occurrenceId;

    // Extract event ID from occurrence ID (assuming format: eventId-date)
    const eventId = occurrenceId.split("-")[0];
    const foundEvent = events.value.find((e) => e.id === eventId);

    if (foundEvent) {
      editingEvent.value = { ...foundEvent };
      isCreatingNewEvent.value = false;
    }
  };

  const setChartDateRange = (months: number) => {
    chartDateRangeMonths.value = months;
  };

  const addEventException = async (eventId: string, date: string) => {
    const event = events.value.find((e) => e.id === eventId);
    if (!event) return;

    const updatedEvent: EventCreate = {
      name: event.name,
      category: event.category,
      account: event.account,
      amount: event.amount,
      start: event.start,
      rrule: event.rrule,
      type: event.type,
      exceptions: {
        ...event.exceptions,
        [date]: { type: "skip" },
      },
    };

    await updateEvent(eventId, updatedEvent);
  };

  const startCreatingNewEvent = () => {
    isCreatingNewEvent.value = true;
    editingEvent.value = null;
    selectedEventOccurrenceId.value = "";
  };

  const deleteEvent = async (eventId: string) => {
    try {
      await apiClient.deleteEvent(eventId);
      await fetchAppState();

      // Clear selection if the deleted event was selected
      if (editingEvent.value?.id === eventId) {
        editingEvent.value = null;
        selectedEventOccurrenceId.value = "";
        isCreatingNewEvent.value = true;
      }

      showToastMessage("Event deleted successfully", "success");
    } catch (error) {
      handleApiError(error);
    }
  };

  return {
    // State
    accounts,
    events,
    appState,
    selectedEventOccurrenceId,
    editingEvent,
    isCreatingNewEvent,
    chartDateRangeMonths,
    toastMessage,
    toastType, // Add this to the return
    showToast,
    // Computed
    filteredAccountBalances,
    // Actions
    fetchAppState,
    createAccount,
    createEvent,
    updateEvent,
    deleteEvent,
    selectEventOccurrence,
    setChartDateRange,
    addEventException,
    startCreatingNewEvent,
    showToastMessage,
  };
});
