<script lang="ts" setup>
import type {DataTableColumns} from 'naive-ui'
import {NButton} from "naive-ui"
import {h, ref} from "vue"

type HostData = {
  key: number
  type: string
  name: string
  ip: string
  uname: string
  ping: string
}

const createColumns = ({play}: { play: (row: HostData) => void }): DataTableColumns<HostData> => {
  return [
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
    },
    {
      title: 'Action',
      key: 'actions',
      render(row) {
        return h(
            NButton,
            {
              strong: true,
              tertiary: true,
              size: 'small',
              onClick: () => play(row)
            },
            {default: () => 'Play'}
        )
      }
    }
  ]
}

const columns: DataTableColumns<HostData> = createColumns({
  play(row: HostData) {
    active.value = true
    console.log(row)
  }
})

const data: HostData[] = Array.from({length: 200}).map((_, index) => ({
  key: index,
  type: Math.random() * 10 > 5 ? 'RDP' : 'SSH',
  name: `Host ${index + 1}`,
  ip: `192.168.1.${index + 1}:${Math.floor(index + Math.random() * 100)}`,
  uname: Math.random() * 10 > 5 ? 'admin' : 'root',
  ping: `${Math.floor(Math.random() * 100)}ms`,
}))

const pagination = false as const

const active = ref(false)


</script>

<template>
  <n-data-table
      class="h-full w-full"
      :pagination="pagination"
      :columns="columns"
      :data="data"
      :bordered="false"
      virtual-scroll
      flex-height
      flex-width
      size="small"
  />

  <n-drawer v-model:show="active" :width="502">
    <n-drawer-content>
      <template #header>
        Header
      </template>
      <template #footer>
        <n-button>Footer</n-button>
      </template>
    </n-drawer-content>
  </n-drawer>
</template>

<style scoped lang="scss">

</style>