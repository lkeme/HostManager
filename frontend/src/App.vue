<script lang="ts" setup>
import {storeToRefs} from 'pinia'
import {useThemeStore} from '@/src/stores'
import {dateZhCN, zhCN} from 'naive-ui'

const themeStore = useThemeStore()
const {theme, themeOverrides} = storeToRefs(themeStore)
</script>

<template>
<!--  class="w-full h-full"-->
  <n-config-provider
      :locale="zhCN"
      :date-locale="dateZhCN"
      :theme="theme"
      :theme-overrides="themeOverrides"
      class="fill-height">
    <n-global-style/>
    <RouterView/>
  </n-config-provider>
</template>


<style lang="scss" scoped>
.n-layout-header,
.n-layout-footer {
  background: rgba(128, 128, 128, 0.2);
  padding: 24px;
}

.n-layout-sider {
  background: rgba(128, 128, 128, 0.3);
}

.n-layout-content {
  background: rgba(128, 128, 128, 0.4);
}
</style>
