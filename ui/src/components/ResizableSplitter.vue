<!-- src/components/ResizableSplitter.vue -->
<template>
  <div class="resizable-container" ref="containerRef">
    <!-- Upper pane -->
    <div
      class="resizable-pane upper-pane"
      :style="{ height: upperPaneHeight + 'px' }"
    >
      <slot name="upper" />
    </div>

    <!-- Divider -->
    <div
      class="resizable-divider"
      @mousedown="startResize"
      :style="{ top: upperPaneHeight + 'px' }"
    >
      <div class="divider-handle">
        <div class="divider-dots">
          <span></span>
          <span></span>
          <span></span>
        </div>
      </div>
    </div>

    <!-- Lower pane -->
    <div
      class="resizable-pane lower-pane"
      :style="{
        height: lowerPaneHeight + 'px',
        top: upperPaneHeight + dividerHeight + 'px',
      }"
    >
      <slot name="lower" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from "vue";

interface Props {
  minUpperHeight?: number;
  minLowerHeight?: number;
  initialUpperRatio?: number;
}

const props = withDefaults(defineProps<Props>(), {
  minUpperHeight: 200,
  minLowerHeight: 150,
  initialUpperRatio: 0.6,
});

const containerRef = ref<HTMLElement>();
const containerHeight = ref(600);
const upperPaneHeight = ref(360);
const dividerHeight = 6;
const isResizing = ref(false);

const lowerPaneHeight = computed(() => {
  return containerHeight.value - upperPaneHeight.value - dividerHeight;
});

const updateContainerHeight = () => {
  if (containerRef.value) {
    const rect = containerRef.value.getBoundingClientRect();
    containerHeight.value = rect.height;

    // Update upper pane height based on current ratio or initial ratio
    const currentRatio =
      upperPaneHeight.value / (containerHeight.value - dividerHeight);
    if (currentRatio === 0 || !isResizing.value) {
      upperPaneHeight.value = Math.floor(
        (containerHeight.value - dividerHeight) * props.initialUpperRatio
      );
    }
  }
};

const startResize = (event: MouseEvent) => {
  isResizing.value = true;
  const startY = event.clientY;
  const startHeight = upperPaneHeight.value;

  const onMouseMove = (e: MouseEvent) => {
    const deltaY = e.clientY - startY;
    let newUpperHeight = startHeight + deltaY;

    // Enforce minimum heights
    newUpperHeight = Math.max(props.minUpperHeight, newUpperHeight);
    newUpperHeight = Math.min(
      containerHeight.value - props.minLowerHeight - dividerHeight,
      newUpperHeight
    );

    upperPaneHeight.value = newUpperHeight;

    // Save to localStorage for persistence
    const ratio = newUpperHeight / (containerHeight.value - dividerHeight);
    localStorage.setItem("splitter-ratio", ratio.toString());
  };

  const onMouseUp = () => {
    isResizing.value = false;
    document.removeEventListener("mousemove", onMouseMove);
    document.removeEventListener("mouseup", onMouseUp);
    document.body.style.cursor = "";
    document.body.style.userSelect = "";
  };

  document.addEventListener("mousemove", onMouseMove);
  document.addEventListener("mouseup", onMouseUp);
  document.body.style.cursor = "ns-resize";
  document.body.style.userSelect = "none";

  event.preventDefault();
};

const loadSavedRatio = () => {
  const savedRatio = localStorage.getItem("splitter-ratio");
  if (savedRatio) {
    const ratio = parseFloat(savedRatio);
    if (ratio > 0 && ratio < 1) {
      upperPaneHeight.value = Math.floor(
        (containerHeight.value - dividerHeight) * ratio
      );
    }
  }
};

let resizeObserver: ResizeObserver;

onMounted(() => {
  updateContainerHeight();
  loadSavedRatio();

  // Watch for container size changes
  if (window.ResizeObserver) {
    resizeObserver = new ResizeObserver(() => {
      updateContainerHeight();
    });

    if (containerRef.value) {
      resizeObserver.observe(containerRef.value);
    }
  }

  // Fallback for browsers without ResizeObserver
  window.addEventListener("resize", updateContainerHeight);
});

onUnmounted(() => {
  if (resizeObserver && containerRef.value) {
    resizeObserver.unobserve(containerRef.value);
  }
  window.removeEventListener("resize", updateContainerHeight);
});
</script>

<style scoped>
.resizable-container {
  height: 100%;
  position: relative;
  overflow: hidden;
}

.resizable-pane {
  position: absolute;
  left: 0;
  right: 0;
  overflow: hidden;
}

.upper-pane {
  top: 0;
  border-bottom: none;
}

.lower-pane {
  bottom: 0;
}

.resizable-divider {
  position: absolute;
  left: 0;
  right: 0;
  height: 6px;
  background-color: #e0e0e0;
  cursor: ns-resize;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: background-color 0.2s;
  z-index: 10;
}

.resizable-divider:hover {
  background-color: #d0d0d0;
}

.resizable-divider:active {
  background-color: #c0c0c0;
}

.divider-handle {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 40px;
  height: 4px;
  background-color: #999;
  border-radius: 2px;
}

.divider-dots {
  display: flex;
  gap: 2px;
}

.divider-dots span {
  width: 3px;
  height: 3px;
  background-color: #666;
  border-radius: 50%;
}

/* Visual feedback during resize */
.resizable-divider:active .divider-handle {
  background-color: #3498db;
}

.resizable-divider:active .divider-dots span {
  background-color: #2980b9;
}
</style>
