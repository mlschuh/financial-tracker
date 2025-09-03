<!-- src/components/FullCalendarPanel.vue -->
<template>
  <div class="fullcalendar-panel">
    <FullCalendar ref="fullCalendar" :options="calendarOptions" />
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
  height: "100%", // Use 100% of container height
  contentHeight: "auto", // Let content determine height within container
  aspectRatio: undefined, // Remove aspect ratio constraints
  headerToolbar: {
    left: "prev,next today",
    center: "title",
    right: "",
  },
  // Ensure calendar fits within container
  handleWindowResize: true,
  dayMaxEvents: 3, // Limit events per day to prevent overflow
  moreLinkClick: "popover", // Show popover for "more" events
}));
</script>

<style scoped>
.fullcalendar-panel {
  height: 100%;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  padding: 10px;
}

/* Ensure FullCalendar takes full height of container */
:deep(.fc) {
  height: 100% !important;
  display: flex;
  flex-direction: column;
}

:deep(.fc-view-harness) {
  flex: 1;
  overflow: hidden;
}

:deep(.fc-scroller) {
  overflow-y: auto !important;
  flex: 1;
}

:deep(.fc-daygrid-body) {
  overflow: hidden;
}

/* Ensure events don't overflow cells */
:deep(.fc-event) {
  font-size: 11px;
  margin: 1px 0;
  overflow: hidden;
  text-overflow: ellipsis;
}

:deep(.fc-daygrid-event) {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

/* Limit day cell height */
:deep(.fc-daygrid-day-frame) {
  min-height: 80px;
  max-height: 120px;
  overflow: hidden;
}

/* Style the "more" link */
:deep(.fc-more-link) {
  font-size: 10px;
  color: #666;
}
</style>
