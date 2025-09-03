<!-- src/components/FullCalendarPanel.vue -->
<template>
  <div class="fullcalendar-panel">
    <FullCalendar ref="fullCalendar" :options="calendarOptions" />
  </div>
</template>

<script setup lang="ts">
import { computed } from "vue";
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
  height: "auto",
  headerToolbar: {
    left: "prev,next today",
    center: "title",
    right: "",
  },
}));
</script>

<style scoped>
.fullcalendar-panel {
  height: 100%;
  padding: 10px;
}

:deep(.fc-event) {
  font-size: 12px;
}

:deep(.fc-daygrid-event) {
  white-space: normal;
}
</style>
