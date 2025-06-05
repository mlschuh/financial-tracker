<script setup>
import { ref, computed } from "vue";
import calendar from "../components/Calendar.vue";
import chart from "../components/Chart.vue";
import detailsPanel from "../components/Details.vue";
import { Splitpanes, Pane } from "splitpanes";
import "splitpanes/dist/splitpanes.css";

const collapse = false;
</script>

<template>
  <splitpanes class="default-theme" horizontal>
    <pane>
      <splitpanes>
        <pane :size="80"><calendar /></pane>
        <pane>List</pane>
      </splitpanes>
    </pane>
    <pane>
      <splitpanes>
        <pane :size="80"><chart /></pane>
        <pane><details /></pane>
      </splitpanes>
    </pane>
  </splitpanes>
</template>

<style scoped></style>
