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

// Reactive calendar options
const calendarOptions = computed(() => ({
  plugins: [dayGridPlugin, timeGridPlugin, interactionPlugin],
  initialView: "dayGridMonth",
  editable: true,
  selectable: true,
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
    console.log(`Comparing ${id} == ${event.eventId}-${event.date}`);
    if (id == `${event.eventId}-${event.date}`) {
      return true;
    }
    return false;
  });
};

const formatEvents = (eventList) => {
  // Format the events if necessary (example: ensure they match FullCalendar's format)
  return eventList.map((event) => {
    // console.log("Working through events");

    return {
      id: `${event.eventId}-${event.date}`, // Event ID
      title: `${event.eventName} (${event.amount})`, // Event title
      start: event.date.split("T")[0], // Start date/time (ISO string)
      // end: event.end, // End date/time (optional, ISO string)
      allDay: true, // Mark as all-day event if specified
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

<style lang="css"></style>
