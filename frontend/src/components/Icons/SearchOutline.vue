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
  <svg  fill="none" viewBox="0 0 512 512">
    <path
        d="M456.69 421.39L362.6 327.3a173.81 173.81 0 0 0 34.84-104.58C397.44 126.38 319.06 48 222.72 48S48 126.38 48 222.72s78.38 174.72 174.72 174.72A173.81 173.81 0 0 0 327.3 362.6l94.09 94.09a25 25 0 0 0 35.3-35.3zM97.92 222.72a124.8 124.8 0 1 1 124.8 124.8a124.95 124.95 0 0 1-124.8-124.8z"
        fill="currentColor"/>
  </svg>

  <svg :height="props.size" :width="props.size" viewBox="0 0 512 512">
    <path d="M221.09 64a157.09 157.09 0 1 0 157.09 157.09A157.1 157.1 0 0 0 221.09 64z" fill="none"
          stroke="currentColor" stroke-miterlimit="10" :stroke-width=props.strokeWidth></path>
    <path fill="none" stroke="currentColor" stroke-linecap="round" stroke-miterlimit="10" :stroke-width=props.strokeWidth
          d="M338.29 338.29L448 448"></path>
  </svg>
</template>

<style lang="scss" scoped>

</style>