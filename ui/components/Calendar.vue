<script setup>
import { ref, onMounted, computed } from "vue";
import FullCalendar from "@fullcalendar/vue3";
import dayGridPlugin from "@fullcalendar/daygrid";
import timeGridPlugin from "@fullcalendar/timegrid";
import interactionPlugin from "@fullcalendar/interaction";
import { useFullState } from "../stores/fullstate.js";

const fullStateStore = useFullState();
const isModalOpen = ref(false);
const selectedDate = ref("");

// Helper function to find an account's color by its ID
const getAccountColor = (accountId) => {
  const account = fullStateStore.accounts.find((acc) => acc.id === accountId);
  return account ? account.color : "#9E9E9E"; // Default to a neutral gray if no color is found
};

// Helper function to get border color based on eventType
const getTypeBorderColor = (eventType) => {
  if (eventType === "income") {
    return "#4CAF50"; // Green for income
  } else if (eventType === "expense") {
    return "#F44336"; // Red for expense
  }
  return "#2196F3"; // Default blue or another neutral color
};

// Reactive calendar options
const calendarOptions = computed(() => ({
  plugins: [dayGridPlugin, timeGridPlugin, interactionPlugin],
  initialView: "dayGridMonth",
  editable: true,
  selectable: true,
  // Use a computed property for events so it reacts to changes in fullStateStore.eventOccurances
  events: formatEvents(fullStateStore.eventOccurances),
  height: "100%",
  contentHeight: "100%",
  fixedWeekCount: false,
  eventClick: handleEventClick,
  dateClick: handleDateClick,
  datesSet: handleDateSet,
}));

const handleDateSet = (dateRange) => {
  // console.log(dateRange)
  // fullStateStore.selectedDateRange = dateRange
};

const handleDateClick = (date) => {
  console.log(date);
  fullStateStore.addEventDate = date;
};

const closeModal = () => {
  isModalOpen.value = false;
};

const findEventOccuranceById = (id) => {
  return fullStateStore.eventOccurances.find((event) => {
    // Corrected comparison logic to match how you construct the id
    if (id === `${event.eventId}-${event.date}`) {
      return true;
    }
    return false;
  });
};

const formatEvents = (eventList) => {
  return eventList.map((event) => {
    const accountColor = getAccountColor(event.accountId); // Get the background color from the account
    const typeBorderColor = getTypeBorderColor(event.eventType); // Get the border color from the event type

    return {
      id: `${event.eventId}-${event.date}`, // Event ID
      title: `${event.eventName} (${event.amount})`, // Event title
      start: event.date.split("T")[0], // Start date/time (ISO string)
      allDay: true, // Mark as all-day event if specified
      backgroundColor: accountColor, // Background color based on account
      borderColor: typeBorderColor, // Border color based on event type
      extendedProps: {
        accountId: event.accountId,
        eventType: event.eventType, // Include eventType in extendedProps if needed later
      },
    };
  });
};

// Function to handle event click
const handleEventClick = (info) => {
  console.log("Event clicked:", info.event.id);
  // You can now access `info.event` to handle the clicked event
  fullStateStore.selectedEventOccurance = findEventOccuranceById(info.event.id);
};

const handleAddEvent = (newEvent) => {
  console.log(newEvent);
};

onMounted(() => {
  fullStateStore.updateFromServer();
  // Refresh the calendar every 5 s
  setInterval(() => {
    fullStateStore.updateFromServer();
  }, 5000);
});
</script>

<template>
  <FullCalendar :options="calendarOptions">
    <template v-slot:eventContent="arg">
      <b>{{ arg.timeText }}</b>
      <i>{{ arg.event.title }}</i>
    </template>
  </FullCalendar>
</template>

<style lang="css">
/* Target FullCalendar event elements */
.fc-event {
  border-width: 2px !important; /* Adjust this value for desired thickness */
}

/* Optional: If you want to ensure the border shows well around the text */
/* This might push the text slightly inward, but makes the border more visible */
.fc-event-main {
  padding-left: 2px; /* Adjust if needed to keep text from touching border */
  padding-right: 2px;
}
</style>
