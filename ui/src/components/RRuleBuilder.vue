<!-- src/components/RRuleBuilder.vue -->
<template>
  <div class="rrule-builder">
    <div class="rrule-row">
      <label class="rrule-label">Repeat:</label>
      <select v-model="frequency" class="rrule-select" @change="updateRRule">
        <option value="">No repeat (single event)</option>
        <option value="DAILY">Daily</option>
        <option value="WEEKLY">Weekly</option>
        <option value="MONTHLY">Monthly</option>
        <option value="YEARLY">Yearly</option>
      </select>
    </div>

    <template v-if="frequency">
      <!-- Interval -->
      <div class="rrule-row">
        <label class="rrule-label">Every:</label>
        <div class="rrule-interval">
          <input
            v-model.number="interval"
            type="number"
            min="1"
            max="99"
            class="rrule-input-small"
            @input="updateRRule"
          />
          <span class="rrule-interval-label">{{ intervalLabel }}</span>
        </div>
      </div>

      <!-- Weekly: Day selection -->
      <div v-if="frequency === 'WEEKLY'" class="rrule-row">
        <label class="rrule-label">On days:</label>
        <div class="day-selector">
          <label v-for="day in weekDays" :key="day.value" class="day-checkbox">
            <input
              v-model="selectedDays"
              type="checkbox"
              :value="day.value"
              @change="updateRRule"
            />
            <span class="day-label">{{ day.short }}</span>
          </label>
        </div>
      </div>

      <!-- Monthly: Day of month or day of week -->
      <div v-if="frequency === 'MONTHLY'" class="rrule-row">
        <label class="rrule-label">By:</label>
        <div class="monthly-options">
          <label class="radio-option">
            <input
              v-model="monthlyType"
              type="radio"
              value="monthday"
              @change="updateRRule"
            />
            Day of month
            <input
              v-if="monthlyType === 'monthday'"
              v-model.number="monthDay"
              type="number"
              min="1"
              max="31"
              class="rrule-input-small"
              @input="updateRRule"
            />
          </label>
          <label class="radio-option">
            <input
              v-model="monthlyType"
              type="radio"
              value="weekday"
              @change="updateRRule"
            />
            Day of week
            <select
              v-if="monthlyType === 'weekday'"
              v-model="monthlyWeekday"
              class="rrule-select-small"
              @change="updateRRule"
            >
              <option value="1MO">First Monday</option>
              <option value="1TU">First Tuesday</option>
              <option value="1WE">First Wednesday</option>
              <option value="1TH">First Thursday</option>
              <option value="1FR">First Friday</option>
              <option value="1SA">First Saturday</option>
              <option value="1SU">First Sunday</option>
              <option value="2MO">Second Monday</option>
              <option value="2TU">Second Tuesday</option>
              <option value="2WE">Second Wednesday</option>
              <option value="2TH">Second Thursday</option>
              <option value="2FR">Second Friday</option>
              <option value="2SA">Second Saturday</option>
              <option value="2SU">Second Sunday</option>
              <option value="3MO">Third Monday</option>
              <option value="3TU">Third Tuesday</option>
              <option value="3WE">Third Wednesday</option>
              <option value="3TH">Third Thursday</option>
              <option value="3FR">Third Friday</option>
              <option value="3SA">Third Saturday</option>
              <option value="3SU">Third Sunday</option>
              <option value="4MO">Fourth Monday</option>
              <option value="4TU">Fourth Tuesday</option>
              <option value="4WE">Fourth Wednesday</option>
              <option value="4TH">Fourth Thursday</option>
              <option value="4FR">Fourth Friday</option>
              <option value="4SA">Fourth Saturday</option>
              <option value="4SU">Fourth Sunday</option>
              <option value="-1MO">Last Monday</option>
              <option value="-1TU">Last Tuesday</option>
              <option value="-1WE">Last Wednesday</option>
              <option value="-1TH">Last Thursday</option>
              <option value="-1FR">Last Friday</option>
              <option value="-1SA">Last Saturday</option>
              <option value="-1SU">Last Sunday</option>
            </select>
          </label>
        </div>
      </div>

      <!-- End options -->
      <div class="rrule-row">
        <label class="rrule-label">Ends:</label>
        <div class="end-options">
          <label class="radio-option">
            <input
              v-model="endType"
              type="radio"
              value="never"
              @change="updateRRule"
            />
            Never
          </label>
          <label class="radio-option">
            <input
              v-model="endType"
              type="radio"
              value="until"
              @change="updateRRule"
            />
            On date
            <input
              v-if="endType === 'until'"
              v-model="untilDate"
              type="date"
              class="rrule-input-date"
              @input="updateRRule"
            />
          </label>
          <label class="radio-option">
            <input
              v-model="endType"
              type="radio"
              value="count"
              @change="updateRRule"
            />
            After
            <input
              v-if="endType === 'count'"
              v-model.number="count"
              type="number"
              min="1"
              max="999"
              class="rrule-input-small"
              @input="updateRRule"
            />
            occurrences
          </label>
        </div>
      </div>

      <!-- Generated RRULE preview -->
      <div class="rrule-preview">
        <label class="rrule-label">Generated rule:</label>
        <div class="rrule-preview-text">{{ generatedRRule || "None" }}</div>
        <div class="rrule-description">{{ humanReadableDescription }}</div>
      </div>
    </template>

    <!-- Manual input fallback -->
    <div class="rrule-row">
      <label class="rrule-label">
        <input
          v-model="useManualInput"
          type="checkbox"
          @change="toggleManualInput"
        />
        Advanced: Edit RRULE manually
      </label>
    </div>

    <div v-if="useManualInput" class="rrule-row">
      <input
        v-model="manualRRule"
        type="text"
        placeholder="e.g., FREQ=MONTHLY;BYMONTHDAY=1"
        class="rrule-input"
        @input="onManualInput"
      />
      <small class="form-help">
        Enter a complete RRULE string.
        <a
          href="https://icalendar.org/iCalendar-RFC-5545/3-8-5-3-recurrence-rule.html"
          target="_blank"
        >
          Learn more about RRULE format
        </a>
      </small>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from "vue";

const props = defineProps<{
  modelValue: string;
}>();

const emit = defineEmits<{
  "update:modelValue": [value: string];
}>();

// Core recurrence settings
const frequency = ref<string>("");
const interval = ref<number>(1);
const endType = ref<string>("never");
const untilDate = ref<string>("");
const count = ref<number>(10);

// Weekly settings
const selectedDays = ref<string[]>([]);

// Monthly settings
const monthlyType = ref<string>("monthday");
const monthDay = ref<number>(1);
const monthlyWeekday = ref<string>("1MO");

// Manual input
const useManualInput = ref<boolean>(false);
const manualRRule = ref<string>("");

const weekDays = [
  { value: "MO", short: "Mon", full: "Monday" },
  { value: "TU", short: "Tue", full: "Tuesday" },
  { value: "WE", short: "Wed", full: "Wednesday" },
  { value: "TH", short: "Thu", full: "Thursday" },
  { value: "FR", short: "Fri", full: "Friday" },
  { value: "SA", short: "Sat", full: "Saturday" },
  { value: "SU", short: "Sun", full: "Sunday" },
];

const intervalLabel = computed(() => {
  if (interval.value === 1) {
    switch (frequency.value) {
      case "DAILY":
        return "day";
      case "WEEKLY":
        return "week";
      case "MONTHLY":
        return "month";
      case "YEARLY":
        return "year";
      default:
        return "";
    }
  } else {
    switch (frequency.value) {
      case "DAILY":
        return "days";
      case "WEEKLY":
        return "weeks";
      case "MONTHLY":
        return "months";
      case "YEARLY":
        return "years";
      default:
        return "";
    }
  }
});

const generatedRRule = computed(() => {
  if (useManualInput.value) {
    return manualRRule.value;
  }

  if (!frequency.value) return "";

  const parts = [`FREQ=${frequency.value}`];

  if (interval.value > 1) {
    parts.push(`INTERVAL=${interval.value}`);
  }

  // Add frequency-specific rules
  if (frequency.value === "WEEKLY" && selectedDays.value.length > 0) {
    parts.push(`BYDAY=${selectedDays.value.join(",")}`);
  }

  if (frequency.value === "MONTHLY") {
    if (monthlyType.value === "monthday") {
      parts.push(`BYMONTHDAY=${monthDay.value}`);
    } else {
      parts.push(`BYDAY=${monthlyWeekday.value}`);
    }
  }

  // Add end conditions
  if (endType.value === "until" && untilDate.value) {
    const until =
      new Date(untilDate.value)
        .toISOString()
        .replace(/[:-]/g, "")
        .split("T")[0] + "T000000Z";
    parts.push(`UNTIL=${until}`);
  } else if (endType.value === "count") {
    parts.push(`COUNT=${count.value}`);
  }

  return parts.join(";");
});

const humanReadableDescription = computed(() => {
  if (!frequency.value) return "No recurrence";

  let description = "";

  // Base frequency
  if (interval.value === 1) {
    switch (frequency.value) {
      case "DAILY":
        description = "Daily";
        break;
      case "WEEKLY":
        description = "Weekly";
        break;
      case "MONTHLY":
        description = "Monthly";
        break;
      case "YEARLY":
        description = "Yearly";
        break;
    }
  } else {
    switch (frequency.value) {
      case "DAILY":
        description = `Every ${interval.value} days`;
        break;
      case "WEEKLY":
        description = `Every ${interval.value} weeks`;
        break;
      case "MONTHLY":
        description = `Every ${interval.value} months`;
        break;
      case "YEARLY":
        description = `Every ${interval.value} years`;
        break;
    }
  }

  // Add day specifications
  if (frequency.value === "WEEKLY" && selectedDays.value.length > 0) {
    const dayNames = selectedDays.value
      .map((day) => weekDays.find((wd) => wd.value === day)?.short || day)
      .join(", ");
    description += ` on ${dayNames}`;
  }

  if (frequency.value === "MONTHLY") {
    if (monthlyType.value === "monthday") {
      description += ` on day ${monthDay.value}`;
    } else {
      const weekdayMap: { [key: string]: string } = {
        "1MO": "first Monday",
        "1TU": "first Tuesday",
        "1WE": "first Wednesday",
        "1TH": "first Thursday",
        "1FR": "first Friday",
        "1SA": "first Saturday",
        "1SU": "first Sunday",
        "2MO": "second Monday",
        "2TU": "second Tuesday",
        "2WE": "second Wednesday",
        "2TH": "second Thursday",
        "2FR": "second Friday",
        "2SA": "second Saturday",
        "2SU": "second Sunday",
        "3MO": "third Monday",
        "3TU": "third Tuesday",
        "3WE": "third Wednesday",
        "3TH": "third Thursday",
        "3FR": "third Friday",
        "3SA": "third Saturday",
        "3SU": "third Sunday",
        "4MO": "fourth Monday",
        "4TU": "fourth Tuesday",
        "4WE": "fourth Wednesday",
        "4TH": "fourth Thursday",
        "4FR": "fourth Friday",
        "4SA": "fourth Saturday",
        "4SU": "fourth Sunday",
        "-1MO": "last Monday",
        "-1TU": "last Tuesday",
        "-1WE": "last Wednesday",
        "-1TH": "last Thursday",
        "-1FR": "last Friday",
        "-1SA": "last Saturday",
        "-1SU": "last Sunday",
      };
      description += ` on the ${
        weekdayMap[monthlyWeekday.value] || monthlyWeekday.value
      }`;
    }
  }

  // Add end conditions
  if (endType.value === "until" && untilDate.value) {
    const date = new Date(untilDate.value).toLocaleDateString();
    description += `, until ${date}`;
  } else if (endType.value === "count") {
    description += `, ${count.value} times`;
  }

  return description;
});

// Initialize from prop
watch(
  () => props.modelValue,
  (newValue) => {
    if (newValue && !useManualInput.value) {
      parseRRule(newValue);
    }
    manualRRule.value = newValue || "";
  },
  { immediate: true }
);

const parseRRule = (rruleString: string) => {
  if (!rruleString) return;

  const parts = rruleString.split(";");
  const rules: { [key: string]: string } = {};

  parts.forEach((part) => {
    const [key, value] = part.split("=");
    if (key && value) {
      rules[key] = value;
    }
  });

  // Parse frequency
  frequency.value = rules.FREQ || "";

  // Parse interval
  interval.value = parseInt(rules.INTERVAL) || 1;

  // Parse weekly days
  if (rules.BYDAY && frequency.value === "WEEKLY") {
    selectedDays.value = rules.BYDAY.split(",");
  }

  // Parse monthly settings
  if (frequency.value === "MONTHLY") {
    if (rules.BYMONTHDAY) {
      monthlyType.value = "monthday";
      monthDay.value = parseInt(rules.BYMONTHDAY) || 1;
    } else if (rules.BYDAY) {
      monthlyType.value = "weekday";
      monthlyWeekday.value = rules.BYDAY;
    }
  }

  // Parse end conditions
  if (rules.UNTIL) {
    endType.value = "until";
    // Convert RRULE date format back to HTML date format
    const dateStr = rules.UNTIL.replace("T000000Z", "");
    const year = dateStr.substring(0, 4);
    const month = dateStr.substring(4, 6);
    const day = dateStr.substring(6, 8);
    untilDate.value = `${year}-${month}-${day}`;
  } else if (rules.COUNT) {
    endType.value = "count";
    count.value = parseInt(rules.COUNT) || 10;
  } else {
    endType.value = "never";
  }
};

const updateRRule = () => {
  if (!useManualInput.value) {
    emit("update:modelValue", generatedRRule.value);
  }
};

const toggleManualInput = () => {
  if (useManualInput.value) {
    emit("update:modelValue", manualRRule.value);
  } else {
    manualRRule.value = generatedRRule.value;
    emit("update:modelValue", generatedRRule.value);
  }
};

const onManualInput = () => {
  if (useManualInput.value) {
    emit("update:modelValue", manualRRule.value);
  }
};
</script>

<style scoped>
.rrule-builder {
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding: 15px;
  border: 1px solid #e9ecef;
  border-radius: 6px;
  background-color: #f8f9fa;
}

.rrule-row {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.rrule-label {
  font-weight: 600;
  color: #495057;
  font-size: 13px;
}

.rrule-select,
.rrule-input {
  padding: 6px 10px;
  border: 1px solid #ced4da;
  border-radius: 4px;
  font-size: 14px;
  background-color: white;
}

.rrule-select-small,
.rrule-input-small,
.rrule-input-date {
  padding: 4px 8px;
  border: 1px solid #ced4da;
  border-radius: 4px;
  font-size: 13px;
  background-color: white;
  margin-left: 8px;
}

.rrule-input-small {
  width: 60px;
}

.rrule-interval {
  display: flex;
  align-items: center;
  gap: 8px;
}

.rrule-interval-label {
  font-size: 14px;
  color: #6c757d;
}

.day-selector {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.day-checkbox {
  display: flex;
  flex-direction: column;
  align-items: center;
  cursor: pointer;
  padding: 6px;
  border-radius: 4px;
  transition: background-color 0.2s;
}

.day-checkbox:hover {
  background-color: #e9ecef;
}

.day-checkbox input[type="checkbox"] {
  margin-bottom: 4px;
}

.day-label {
  font-size: 12px;
  color: #495057;
}

.monthly-options,
.end-options {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.radio-option {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  cursor: pointer;
}

.rrule-preview {
  margin-top: 8px;
  padding: 12px;
  background-color: #e7f3ff;
  border-left: 4px solid #3498db;
  border-radius: 4px;
}

.rrule-preview-text {
  font-family: "Monaco", "Menlo", monospace;
  font-size: 12px;
  color: #2c3e50;
  margin-bottom: 6px;
  word-break: break-all;
}

.rrule-description {
  font-size: 13px;
  color: #3498db;
  font-weight: 500;
}

.form-help {
  margin-top: 4px;
  color: #6c757d;
  font-size: 12px;
}

.form-help a {
  color: #3498db;
  text-decoration: none;
}

.form-help a:hover {
  text-decoration: underline;
}
</style>
