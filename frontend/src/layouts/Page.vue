<script setup lang="ts">
import HeaderPanel from "@/src/layouts/HeaderPanel.vue";
import FooterPanel from "@/src/layouts/FooterPanel.vue";
import BodyPanel from "@/src/layouts/BodyPanel/index.vue";
</script>

<template>
  <main class="flex flex-col h-screen w-screen overflow-hidden">
    <!-- 头部 -->
    <header class="flex h-10">
      <HeaderPanel/>
    </header>
    <!-- 主体 -->
    <section class="flex-1 overflow-y-auto">
      <slot></slot>
    </section>
    <!-- 底部 -->
    <footer class="flex h-7">
      <FooterPanel/>
    </footer>
  </main>
</template>


<style lang="scss" scoped>


</style>