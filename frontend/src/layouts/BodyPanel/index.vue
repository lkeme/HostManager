<template>
  <div class="flex flex-row w-full h-full relative drawer drawer-end duration-75" >
    <!-- Host Tab-->
    <div class="flex flex-col w-[300px]">
      <HostTab/>
    </div>
    <n-divider vertical />
<!--        <div class="w-[300px] absolute bottom-0 top-0 left-0  bg-amber-100">-->
<!--        <div class="flex-1 flex-col bg-amber-700">-->
<!--          &lt;!&ndash; Host List&ndash;&gt;-->
<!--&lt;!&ndash;          <HostList/>&ndash;&gt;-->
<!--          &lt;!&ndash; Host Info&ndash;&gt;-->
<!--&lt;!&ndash;          <HostInfo/>&ndash;&gt;-->
<!--        </div>-->

  </div>


  <!--  <n-card class="flex h-full hidden">-->
  <!--    <template #header>-->
  <!--      <n-data-table-->
  <!--          class="h-full w-full"-->
  <!--          :columns="columns"-->
  <!--          :data="data"-->
  <!--          :bordered="false"-->
  <!--          virtual-scroll-->
  <!--          flex-height-->
  <!--          flex-width-->
  <!--          size="small"-->
  <!--      />-->
  <!--    </template>-->

  <!--    <template #action>-->
  <!--      这里显示列的信息-->
  <!--    </template>-->
  <!--  </n-card>-->


</template>

<script setup lang="ts">
import type {DataTableColumns} from 'naive-ui'
import HostTab from "@/src/layouts/BodyPanel/HostTab.vue";
import HostList from "@/src/layouts/BodyPanel/HostList.vue";
import HostInfo from "@/src/layouts/BodyPanel/HostInfo.vue";

type RowData = {
  key: number
  type: string
  name: string
  ip: string
  uname: string
  ping: string
}

const columns: DataTableColumns<RowData> = [
  {
    type: 'selection',
    fixed: 'left'
  },
  {
    title: '类型',
    key: 'type',
    fixed: 'left',
    width: 100,
  },
  {
    title: '名称',
    key: 'name',
    fixed: 'left',
    width: 120,
  },
  {
    title: '地址',
    key: 'ip',
    fixed: 'left'
  },
  {
    title: '用户名',
    key: 'uname',
    width: 120,
    fixed: 'left'
  },
  {
    title: '延迟',
    key: 'ping',
    width: 80,
    fixed: 'right'
  }
]

const data: RowData[] = Array.from({length: 200}).map((_, index) => ({
  key: index,
  type: Math.random() * 10 > 5 ? 'RDP' : 'SSH',
  name: `Host ${index + 1}`,
  ip: `192.168.1.${index + 1}:${Math.floor(index + Math.random() * 100)}`,
  uname: Math.random() * 10 > 5 ? 'admin' : 'root',
  ping: `${Math.floor(Math.random() * 100)}ms`,
}))


</script>


<style lang="scss" scoped>
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