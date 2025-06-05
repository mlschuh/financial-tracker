<template>
  <div class="chart-container">
    <div class="date-range-selectors">
      <label for="before-select">Before:</label>
      <select
        id="before-select"
        v-model="selectedBefore"
        @change="updateDateRange"
      >
        <option
          v-for="option in durationOptions"
          :key="option.value"
          :value="option.value"
        >
          {{ option.label }}
        </option>
      </select>

      <label for="after-select">After:</label>
      <select
        id="after-select"
        v-model="selectedAfter"
        @change="updateDateRange"
      >
        <option
          v-for="option in durationOptions"
          :key="option.value"
          :value="option.value"
        >
          {{ option.label }}
        </option>
      </select>
    </div>
    <Line :data="reactiveChartData" :options="chartOptions" ref="chartCanvas" />
    <!-- <canvas ref="chartCanvas"></canvas> -->
  </div>
</template>

<script setup>
import { Line } from "vue-chartjs";
import { ref, onMounted, onUnmounted, computed, watch } from "vue";
import { Chart, registerables } from "chart.js";
import "chartjs-adapter-date-fns";
import { useFullState } from "../stores/fullstate.js";
import {
  addMonths,
  subMonths,
  startOfDay,
  endOfDay,
  isWithinInterval,
} from "date-fns";

const fullStateStore = useFullState();

// Register Chart.js components globally
Chart.register(...registerables);

const chartCanvas = ref(null); // Reference to the canvas element

// Define localStorage keys
const LOCALSTORAGE_BEFORE_KEY = "chartDateRangeBefore";
const LOCALSTORAGE_AFTER_KEY = "chartDateRangeAfter";

// Combo box state
const selectedBefore = ref(
  parseInt(localStorage.getItem(LOCALSTORAGE_BEFORE_KEY)) || 1
);
const selectedAfter = ref(
  parseInt(localStorage.getItem(LOCALSTORAGE_AFTER_KEY)) || 1
);
const durationOptions = [
  { label: "1 month", value: 1 },
  { label: "3 months", value: 3 },
  { label: "6 months", value: 6 },
  { label: "12 months", value: 12 },
  { label: "24 months", value: 24 },
  { label: "48 months", value: 48 },
];

/**
 * Updates the selectedDateRange in the store based on the 'Before' and 'After' selections.
 */
const updateDateRange = () => {
  const today = new Date(); // Or a fixed reference date if needed

  // Calculate the 'start' date based on 'Before' selection
  const newStartDate = subMonths(today, selectedBefore.value);
  // Ensure the start date is the beginning of that day
  fullStateStore.selectedDateRange.start =
    startOfDay(newStartDate).toISOString();

  // Calculate the 'end' date based on 'After' selection
  const newEndDate = addMonths(today, selectedAfter.value);
  // Ensure the end date is the end of that day
  fullStateStore.selectedDateRange.end = endOfDay(newEndDate).toISOString();
};

watch(selectedBefore, (newValue) => {
  localStorage.setItem(LOCALSTORAGE_BEFORE_KEY, newValue.toString());
  updateDateRange(); // Recalculate date range when value changes
});

watch(selectedAfter, (newValue) => {
  localStorage.setItem(LOCALSTORAGE_AFTER_KEY, newValue.toString());
  updateDateRange(); // Recalculate date range when value changes
});

// Create a computed map for quick account name and color lookup
const accountDetailsMap = computed(() => {
  const map = {};
  if (fullStateStore.accounts) {
    fullStateStore.accounts.forEach((account) => {
      map[account.id] = {
        name: account.name,
        color: account.color || getRandomColor(), // Use provided color, or fallback to random
      };
    });
  }
  return map;
});

// This computed property processes accountBalances to create datasets for the chart.
const chartDataSets = computed(() => {
  const accountBalances = fullStateStore.accountBalances;
  const start = new Date(fullStateStore.selectedDateRange.start);
  const end = new Date(fullStateStore.selectedDateRange.end);
  const detailsMap = accountDetailsMap.value; // Use the computed map

  // Group balances by accountId
  const balancesByAccount = accountBalances.reduce((acc, item) => {
    // Filter out balances outside the selected date range
    const itemDate = new Date(item.date);
    if (isWithinInterval(itemDate, { start, end })) {
      if (!acc[item.accountId]) {
        acc[item.accountId] = [];
      }
      acc[item.accountId].push({ x: itemDate, y: item.balance });
    }
    return acc;
  }, {});

  // Convert grouped balances into Chart.js dataset format
  const datasets = Object.keys(balancesByAccount).map((accountId) => {
    // Sort data points by date for proper line drawing
    const sortedData = balancesByAccount[accountId].sort(
      (a, b) => a.x.getTime() - b.x.getTime()
    );

    // Look up the account details, default to accountId if not found
    const accountDetail = detailsMap[accountId] || {
      name: `Account: ${accountId}`,
      color: getRandomColor(),
    };

    return {
      label: `Account: ${accountDetail.name}`, // You might want to map accountId to a more friendly name
      data: sortedData,
      borderColor: accountDetail.color, // Assign a random color for different accounts
      backgroundColor: accountDetail.color, // Lighter background color for fill
      fill: false, // Set to true if you want the area under the line filled
      tension: 0.1, // Smoothness of the line
      stepped: true,
    };
  });

  return datasets;
});

// This computed property will provide the full chart data object
// that Vue-Chartjs expects. When `chartDataPoints` changes, this
// computed property will re-evaluate, and the `Line` component will
// detect the change in its `data` prop and re-render.
const reactiveChartData = computed(() => {
  return {
    // Labels are technically not strictly necessary for time scales when data has {x, y}
    // but it's good practice to provide them or at least know what they represent.
    // For time scale, Chart.js typically uses the 'x' values from your dataset's data.
    labels: [],
    datasets: chartDataSets.value,
  };
});
const chartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  animation: {
    duration: 300, // Animation duration in milliseconds (e.g., 500ms = 0.5 seconds)
    // You can also change the easing function for different animation effects
    // easing: 'linear', // 'linear', 'easeInQuad', 'easeOutQuad', 'easeInOutQuad', etc.
  },
  scales: {
    x: {
      type: "time", // Use the timescale
      time: {
        unit: "day", // Display ticks by day
      },
      title: {
        display: true,
        text: "Date",
      },
    },
    y: {
      beginAtZero: true,
      title: {
        display: true,
        text: "Value",
      },
    },
  },
};

// Helper function to generate a random color for chart lines
const getRandomColor = (alpha = 1) => {
  const r = Math.floor(Math.random() * 255);
  const g = Math.floor(Math.random() * 255);
  const b = Math.floor(Math.random() * 255);
  return `rgba(${r}, ${g}, ${b}, ${alpha})`;
};

// Initialize the Chart.js instance
onMounted(() => {
  updateDateRange(); // Set initial dates based on default combo box values
});
</script>

<style scoped>
.chart-container {
  position: relative;
  height: 80%;
  width: 100%;
  display: flex;
  flex-direction: column; /* Arrange items vertically */
  align-items: center; /* Center horizontally */
}

.date-range-selectors {
  margin-bottom: 20px; /* Space between selectors and chart */
  display: flex;
  gap: 15px; /* Space between labels and selects */
  align-items: center;
}

.date-range-selectors label {
  font-weight: bold;
}

.date-range-selectors select {
  padding: 8px;
  border: 1px solid #ccc;
  border-radius: 4px;
  cursor: pointer;
  min-width: 120px; /* Ensure consistent width */
}
</style>
