<script setup>
import { defineComponent } from "vue";
import calendar from "../components/Calendar.vue";
import chart from "../components/Chart.vue";
import detailsPanel from "../components/Details.vue";

</script>

<template>
  <div class="container">
    <div class="first-column">
      <div class="calendar">
        <calendar />
      </div>
      <div class="chart">
        Chart Pane
        <chart />
      </div>
    </div>
    <div>
      <h2>Details</h2>
      <detailsPanel />
    </div>
  </div>
</template>

<style scoped>
/* Container sets up the grid layout */
.container {
  display: grid;
  grid-template-columns: 1fr 350px; /* Two columns: First column is 2/3 width, second is 1/3 width */
  height: 100vh; /* Full viewport height */
}

/* First column with two rows */
.first-column {
  display: grid;
  grid-template-rows: 2fr 200px; /* Top row is 2/3 height, bottom row is 1/3 height */
}

/* Styling for rows and columns */
.top-row,
.bottom-row,
.second-column {
  display: flex;
  align-items: center;
  justify-content: center;
  border: 1px solid #ccc;
  background-color: #f9f9f9;
}

.chart {
  width: 100%;
  height: 100%;
  background-color: pink;
}

</style>
