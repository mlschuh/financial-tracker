<template>
  <div class="chart-container">
    <canvas ref="chartCanvas"></canvas>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, computed } from "vue";
import { Chart, registerables } from "chart.js";
import "chartjs-adapter-date-fns";
import { useFullState } from "../stores/fullstate.js";

const fullStateStore = useFullState();

// Register Chart.js components globally
Chart.register(...registerables);

const chartCanvas = ref(null); // Reference to the canvas element
let chartInstance = null; // Variable to store the Chart.js instance

let chartLabels = computed(() => {
  const start = new Date(fullStateStore.selectedDateRange.start); // Ensure input is converted to Date objects
  const end = new Date(fullStateStore.selectedDateRange.end);
  const dates = [];

  // Generate dates
  for (let d = new Date(start); d <= end; d.setDate(d.getDate() + 1)) {
    dates.push(new Date(d)); // Push a copy of the current date
  }

  return dates;
});

// Example chart configuration (replace or extend this as needed)
const chartConfig = {
  type: "line", // Chart type (e.g., 'line', 'bar', 'doughnut', etc.)
  data: {
    labels: computed(() => chartLabels),
    datasets: [
      {
        label: "Value Over Time",
        data: [
          { x: "2025-01-01T00:00:00Z", y: 10 },
          { x: "2025-01-02T00:00:00Z", y: 20 },
          { x: "2025-01-03T00:00:00Z", y: 15 },
        ],
        borderColor: "blue",
        backgroundColor: "rgba(0, 0, 255, 0.1)",
        fill: true,
      },
    ],
  },
  options: {
    responsive: true,
    maintainAspectRatio: false,
    //   scales: {
    //     y: {
    //       beginAtZero: true,
    //     },
    //   },
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
  },
};

// Initialize the Chart.js instance
onMounted(() => {
  if (chartCanvas.value) {
    chartInstance = new Chart(chartCanvas.value, chartConfig);
    // window.addEventListener("resize", () => {
    //     console.log("Rezising window")
    //   chartInstance.resize();
    // });
  }
});

// Destroy the Chart.js instance on unmount to avoid memory leaks
onUnmounted(() => {
  if (chartInstance) {
    chartInstance.destroy();
  }
});
</script>

<style scoped>
.chart-container {
  position: relative;
  height: 100%;
  width: 80vw;
}
</style>
