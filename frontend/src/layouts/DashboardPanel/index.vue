<template>
  <div class="flex flex-row w-full h-full">
    <!-- Host Tab-->
    <div class="min-w-[300px]">
      <HostTab/>
    </div>
    <!-- Divider-->
    <div class="min-w-[1px]" @mousedown="leftResize($event)">
      <n-divider vertical/>
    </div>
    <!-- Host Info-->
    <div class="flex-1">
      <HostData/>
    </div>
  </div>
</template>

<script setup lang="ts">
import {ref} from "vue";
import HostTab from "@/src/layouts/DashboardPanel/HostTab.vue";
import HostData from "@/src/layouts/DashboardPanel/HostData.vue";


const leftBoxWidth = ref(300); // 默认宽度

const style = ref({
  width: leftBoxWidth + 'px',
  minWidth: '300px',
  maxWidth: '600px'
});

const leftResize = (e: MouseEvent) => {
  document.onselectstart = function () {
    return false;
  }; //解决拖动会选中文字的问题
  const startX = e.clientX;
  const startWidth = leftBoxWidth.value;
  const mouseMove = (documentE: MouseEvent) => {
    // 80 是左侧菜单宽度
    leftBoxWidth.value = startWidth + documentE.clientX - startX - 80;
  };
  const mouseUp = () => {
    document.removeEventListener('mousemove', mouseMove);
    document.removeEventListener('mouseup', mouseUp);
    // 拖拽完记得重新设置可以选中
    document.onselectstart = function () {
      return true;
    };
  };
  document.addEventListener('mousemove', mouseMove);
  document.addEventListener('mouseup', mouseUp);
};
</script>

<style lang="scss" scoped>
//.drag-bar {
//  width: 10px;
//  cursor: col-resize;
//  z-index: 999;
//
//  &:hover {
//    background-color: $color;
//  }
//}

.drag-bar {
  position: absolute;
  top: 0;
  left: -5px;
  height: 100%;
}

.n-divider {
  margin: 0;
  height: 100%;
}

.n-card {
  ::v-deep(.n-card-header) {
    flex-basis: 75%;
    padding: 0;
    background-color: unset;
    align-items: unset;
  }

  ::v-deep(.n-card__action) {
    flex-basis: 25%;
    padding: 0;
    //background-color: unset;
  }
}

</style>