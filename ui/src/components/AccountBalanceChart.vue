<!-- src/components/AccountBalanceChart.vue -->
<template>
  <div class="chart-panel">
    <div class="chart-controls">
      <h3>Account Balances</h3>
      <div class="date-range-buttons">
        <button
          v-for="range in dateRanges"
          :key="range.months"
          @click="store.setChartDateRange(range.months)"
          :class="[
            'range-btn',
            { active: store.chartDateRangeMonths === range.months },
          ]"
        >
          {{ range.label }}
        </button>
      </div>
    </div>
    <div class="chart-container">
      <Line
        v-if="chartData.datasets.length > 0"
        :data="chartData"
        :options="chartOptions"
      />
      <div v-else class="no-data">
        No balance data available for the selected time range
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from "vue";
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
  TimeScale,
} from "chart.js";
import annotationPlugin from "chartjs-plugin-annotation";
import "chartjs-adapter-date-fns";
import { Line } from "vue-chartjs";
import { useAppStore } from "@/stores/appStore";

ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
  TimeScale,
  annotationPlugin
);

const store = useAppStore();

const dateRanges = [
  { label: "3 months", months: 3 },
  { label: "6 months", months: 6 },
  { label: "12 months", months: 12 },
  { label: "24 months", months: 24 },
];

const chartData = computed(() => {
  const balances = store.filteredAccountBalances;
  const accounts = store.accounts;

  // Group balances by account
  const accountBalances = new Map<
    string,
    Array<{ date: string; balance: number }>
  >();

  balances.forEach((balance) => {
    if (!accountBalances.has(balance.accountId)) {
      accountBalances.set(balance.accountId, []);
    }
    accountBalances.get(balance.accountId)!.push({
      date: balance.date,
      balance: balance.balance,
    });
  });

  // Sort balances by date for each account
  accountBalances.forEach((balanceArray) => {
    balanceArray.sort(
      (a, b) => new Date(a.date).getTime() - new Date(b.date).getTime()
    );
  });

  // Get all unique dates and sort them
  const allDates = Array.from(new Set(balances.map((b) => b.date))).sort(
    (a, b) => new Date(a).getTime() - new Date(b).getTime()
  );

  const datasets = Array.from(accountBalances.entries()).map(
    ([accountId, accountBalanceData]) => {
      const account = accounts.find((a) => a.id === accountId);

      // Create data points with actual Date objects for better time scale handling
      const dataPoints = allDates
        .map((date) => {
          const balanceEntry = accountBalanceData.find((b) => b.date === date);
          return {
            x: new Date(date),
            y: balanceEntry ? balanceEntry.balance / 100 : null, // Convert to dollars
          };
        })
        .filter((point) => point.y !== null); // Remove null values

      return {
        label: account?.name || `Account ${accountId}`,
        data: dataPoints,
        borderColor: account?.color || "#3498db",
        backgroundColor: account?.color || "#3498db",
        tension: 0.1,
        spanGaps: true,
      };
    }
  );

  return {
    datasets,
  };
});

const chartOptions = computed(() => {
  const today = new Date();
  today.setHours(0, 0, 0, 0); // Set to start of day for precise comparison

  return {
    responsive: true,
    maintainAspectRatio: false,
    interaction: {
      mode: "index" as const,
      intersect: false,
    },
    plugins: {
      legend: {
        position: "top" as const,
      },
      title: {
        display: false,
      },
      tooltip: {
        callbacks: {
          label: function (context: any) {
            return `${context.dataset.label}: $${context.parsed.y.toFixed(2)}`;
          },
          title: function (context: any) {
            const date = new Date(context[0].parsed.x);
            return date.toLocaleDateString();
          },
        },
      },
      // Today's date annotation
      annotation: {
        annotations: {
          todayLine: {
            type: "line" as const,
            xMin: today,
            xMax: today,
            borderColor: "#ff6b6b",
            borderWidth: 2,
            borderDash: [5, 5],
            label: {
              display: true,
              content: "Today",
              position: "start" as const,
              yAdjust: -10,
              backgroundColor: "#ff6b6b",
              color: "white",
              borderRadius: 4,
              padding: {
                x: 6,
                y: 4,
              },
              font: {
                size: 11,
                weight: "bold" as const,
              },
            },
          },
        },
      },
    },
    scales: {
      x: {
        type: "time" as const,
        time: {
          unit: "day" as const,
          displayFormats: {
            day: "MMM dd",
          },
          tooltipFormat: "MMM dd, yyyy",
        },
        title: {
          display: true,
          text: "Date",
        },
      },
      y: {
        beginAtZero: false,
        title: {
          display: true,
          text: "Balance ($)",
        },
        ticks: {
          callback: function (value: any) {
            return "$" + value.toFixed(2);
          },
        },
      },
    },
    elements: {
      point: {
        radius: 3,
        hoverRadius: 6,
      },
    },
  };
});
</script>

<style scoped>
.chart-panel {
  height: 100%;
  padding: 15px;
  display: flex;
  flex-direction: column;
}

.chart-controls {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
  flex-shrink: 0;
}

.chart-controls h3 {
  margin: 0;
  color: #333;
}

.date-range-buttons {
  display: flex;
  gap: 5px;
}

.range-btn {
  padding: 6px 12px;
  border: 1px solid #ddd;
  background-color: #fff;
  color: #333;
  border-radius: 4px;
  cursor: pointer;
  font-size: 12px;
  transition: all 0.2s;
}

.range-btn:hover {
  background-color: #f8f9fa;
}

.range-btn.active {
  background-color: #3498db;
  color: white;
  border-color: #3498db;
}

.chart-container {
  flex: 1;
  position: relative;
  min-height: 0;
}

.no-data {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: #666;
  font-style: italic;
}
</style>
