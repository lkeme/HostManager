<script setup lang="ts">
import HeaderPanel from "@/src/layouts/HeaderPanel/index.vue";
import FooterPanel from "@/src/layouts/FooterPanel/index.vue";
</script>

<template>
  <main class="flex flex-col min-h-screen w-screen relative overflow-hidden duration-150 transition-all">
    <!-- 头部 h-10 -->
    <header class="flex absolute h-10 top-0 left-0 right-0">
      <HeaderPanel/>
    </header>
    <!-- 主体 -->
    <section class="flex absolute top-10 bottom-7 left-0 right-0">
      <slot></slot>
    </section>
    <!-- 底部 -->
    <footer class="flex absolute h-7 bottom-0 left-0 right-0">
      <FooterPanel/>
    </footer>
  </main>
</template>

<style lang="scss" scoped>
</style>