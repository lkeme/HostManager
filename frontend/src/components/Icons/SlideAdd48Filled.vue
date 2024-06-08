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
  <svg :height="props.size" :width="props.size" fill="none" viewBox="0 0 48 48">
    <g fill="none">
      <path
          d="M10.25 8A6.25 6.25 0 0 0 4 14.25v19.5A6.25 6.25 0 0 0 10.25 40h12.746A12.962 12.962 0 0 1 22 35c0-7.18 5.82-13 13-13c3.493 0 6.664 1.378 9 3.62V14.25A6.25 6.25 0 0 0 37.75 8h-27.5zM35 46c6.075 0 11-4.925 11-11s-4.925-11-11-11s-11 4.925-11 11s4.925 11 11 11zm0-19a1 1 0 0 1 1 1v6h6a1 1 0 1 1 0 2h-6v6a1 1 0 1 1-2 0v-6h-6a1 1 0 1 1 0-2h6v-6a1 1 0 0 1 1-1z"
          fill="currentColor"
          stroke="currentColor"/>
    </g>
  </svg>
</template>

<style lang="scss" scoped>

</style>