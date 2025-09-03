<!-- Alternative version with more explicit scrolling control -->
<template>
  <div class="fullcalendar-panel">
    <div class="calendar-header">
      <h3>Financial Calendar</h3>
    </div>
    <div class="calendar-wrapper">
      <FullCalendar ref="fullCalendar" :options="calendarOptions" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, nextTick } from "vue";
import FullCalendar from "@fullcalendar/vue3";
import dayGridPlugin from "@fullcalendar/daygrid";
import rrulePlugin from "@fullcalendar/rrule";
import { useAppStore } from "@/stores/appStore";
import type { EventApi } from "@fullcalendar/core";

const store = useAppStore();

const calendarEvents = computed(() => {
  return store.appState.eventOccurances.map((occurrence) => {
    const account = store.accounts.find((a) => a.id === occurrence.accountId);
    const displayAmount = (occurrence.amount / 100).toFixed(2);

    return {
      id: occurrence.id,
      title: `${occurrence.eventName} ($${displayAmount})`,
      start: occurrence.date,
      end: occurrence.date,
      allDay: true,
      color: account?.color || "#3498db",
      extendedProps: {
        eventId: occurrence.eventId,
        accountId: occurrence.accountId,
        eventType: occurrence.eventType,
        amount: occurrence.amount,
      },
    };
  });
});

const calendarOptions = computed(() => ({
  plugins: [dayGridPlugin, rrulePlugin],
  initialView: "dayGridMonth",
  events: calendarEvents.value,
  eventClick: (info: { event: EventApi }) => {
    store.selectEventOccurrence(info.event.id);
  },
  height: "auto", // Let content determine height
  contentHeight: "auto", // Auto-size content
  aspectRatio: 1.35, // Maintain reasonable aspect ratio
  headerToolbar: {
    left: "prev,next today",
    center: "title",
    right: "",
  },
  handleWindowResize: true,
  dayMaxEvents: false, // Show all events
  dayGridEventMinHeight: 15,
  // Allow natural sizing
  fixedWeekCount: false, // Don't always show 6 weeks
}));
</script>

<style scoped>
.fullcalendar-panel {
  height: 100%;
  display: flex;
  flex-direction: column;
  padding: 10px;
  overflow: hidden;
}

.calendar-header {
  flex-shrink: 0;
  margin-bottom: 10px;
  padding-bottom: 10px;
  border-bottom: 1px solid #eee;
}

.calendar-header h3 {
  margin: 0;
  color: #333;
  font-size: 16px;
}

.calendar-wrapper {
  flex: 1;
  overflow-y: auto; /* Enable scrolling at wrapper level */
  overflow-x: hidden;
  border: 1px solid #ddd;
  border-radius: 4px;
  background: white;
}

/* FullCalendar styling for natural sizing */
:deep(.fc) {
  height: auto !important; /* Let it size naturally */
}

:deep(.fc-view-harness) {
  overflow: visible;
}

:deep(.fc-scroller) {
  overflow: visible !important; /* Don't clip content */
}

/* Allow day cells to expand as needed */
:deep(.fc-daygrid-day-frame) {
  min-height: 80px;
  padding: 2px;
}

/* Rest of the event styling remains the same */
:deep(.fc-event) {
  font-size: 11px;
  margin: 1px;
  padding: 2px 4px;
  border-radius: 3px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  cursor: pointer;
  transition: opacity 0.2s;
}

:deep(.fc-event:hover) {
  opacity: 0.8;
}

/* Style the scrollbar for the wrapper */
.calendar-wrapper::-webkit-scrollbar {
  width: 8px;
}

.calendar-wrapper::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 4px;
}

.calendar-wrapper::-webkit-scrollbar-thumb {
  background: #c1c1c1;
  border-radius: 4px;
}

.calendar-wrapper::-webkit-scrollbar-thumb:hover {
  background: #a8a8a8;
}
</style>
