<script lang="ts" setup>
const props = defineProps({
  size: {
    type: [Number, String],
    default: 14,
  },
  strokeWidth: {
    type: [Number, String],
    default: 4,
  },
})
</script>

<template>
  <svg id="icon" :width="props.size" :height="props.size" viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg">
    <circle cx="50" cy="50" r="40" stroke="white" :stroke-width=props.strokeWidth fill="none" />
    <line x1="50" y1="10" x2="50" y2="90" stroke="white" stroke-width="3" />
    <line x1="10" y1="50" x2="90" y2="50" stroke="white" stroke-width="3" />
  </svg>
</template>

<style lang="scss" scoped>
@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

#icon {
  animation: spin 2s infinite linear;
}
</style>