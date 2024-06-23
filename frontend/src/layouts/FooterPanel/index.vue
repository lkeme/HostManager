<script setup lang="ts">
import {BrowserOpenURL, WindowSetAlwaysOnTop} from "@/wailsjs/runtime";
import {type Component, computed, h, onMounted, ref} from "vue";
import Setting from "@/src/components/Icons/Setting.vue";
import Pin32Filled from "@/src/components/Icons/Pin32Filled.vue";
import {NIcon} from "naive-ui";
import {type FormSchema, ModalForm, type ModalFormInstance, type Recordable} from "naive-ui-form";
import {useThemeStore} from '@/src/stores'
import {storeToRefs} from 'pinia'
import {GetAppSettingAttr} from "@/wailsjs/go/controller/Config";
import {GetAppSettingAttrResponse} from "@/src/types/response/config";
import {IsEncrypted, IsLogin} from "@/wailsjs/go/controller/Auth";
import {IsEncryptedResponse, IsLoginResponse} from "@/src/types/response/auth";

const themeStore = useThemeStore()
const { themeColors} = storeToRefs(themeStore)

const props = defineProps({
  showPin2Top: {
    type: Boolean,
    default: false,
    required: false
  }
})

const isPinned = ref<boolean>(false);

const pinWindow = () => {
  isPinned.value = !isPinned.value
  WindowSetAlwaysOnTop(isPinned.value)
  isPinned.value ? $message.success("已置顶") : $message.success("已取消置顶")
}

const renderIcon = (icon: Component) => {
  return () => {
    return h(NIcon, null, {
      default: () => h(icon)
    })
  }
}

const preferencesOptions = computed(() => {
  return [
    {
      label: '偏好设置',
      key: 'preferences',
      icon: renderIcon(Setting)
    },
    {
      label: '报告错误',
      key: 'report',
    },
    {
      label: '使用指南',
      key: 'help',
    },
    {
      label: '检查更新',
      key: 'update',
    },
    {
      type: 'divider',
      key: 'd1',
    },
    {
      label: '关于',
      key: 'about',
    },
  ]
})

const onSelectPreferenceMenu = (key: string) => {
  switch (key) {
    case 'preferences':
      break
    case 'update':
      break
    case 'report':
      break
    case 'help':
      break
    case 'about':
      showAboutModal.value = true
      break
  }
}

// 格式：XXXX年XX月XX日XX时XX分XX秒 星期X
const complement = function (value: any) {
  return value < 10 ? `0${value}` : value;
};

const formatDate = (date: any) => {
  const time = new Date(date);
  const year = time.getFullYear();
  const month = complement(time.getMonth() + 1);
  const day = complement(time.getDate());
  const hour = complement(time.getHours());
  const minute = complement(time.getMinutes());
  const second = complement(time.getSeconds());
  const week = '星期' + '日一二三四五六'.charAt(time.getDay());
  return `${year}年${month}月${day}日 ${week} ${hour}:${minute}:${second}`;
};

const nowTime = ref("")
onMounted(() => {
  setInterval(() => {
    nowTime.value = formatDate(new Date())
  })
});


// 关于弹窗
const showAboutModal = ref(false);
const showAboutModalRef = ref<ModalFormInstance | null>(null);
const showAboutModalSchemas: FormSchema[] = [
  {
    field: 'template',
    type: 'slot',
    slot: 'template',
    componentProps: {
      showFeedback: false
    }
  }
];
const showAboutModalLoading = ref(false);
const handleShowAboutModal = async (values: Recordable) => {
  console.log(values)
};

const onOpenSource = () => {
  BrowserOpenURL('https://github.com/lkeme/HostManager')
}

const onOpenWebsite = () => {
  BrowserOpenURL('https://github.com/')
}

const appIcon = ref<string>('')

onMounted(async () => {
  await GetAppSettingAttr('appIcon').then((response: GetAppSettingAttrResponse) => {
    appIcon.value = response.data.value;
  });
});


</script>

<template>
  <div class="flex flex-col h-full w-full">
    <n-divider/>
    <div class="flex flex-1 items-center justify-between mx-2 bottom-0 z-50">
      <p>拥有<span class="font-bold">23</span>个主机</p>
      <!--隐藏/显示-->
      <div class="group">
        <div class="group-hover:hidden">
          Happy Ending
        </div>
        <div class="hidden group-hover:flex">
          <n-gradient-text gradient="linear-gradient(90deg, red 0%, green 50%, blue 100%)">
            {{ nowTime }}
          </n-gradient-text>
        </div>
      </div>
      <div class="flex gap-x-2.5 ">
        <!--      absolute right-2 bottom-2-->
        <n-tooltip :delay="1000" :show-arrow="false">
          窗口置顶
          <template #trigger>
            <n-icon
                v-show="props.showPin2Top"
                :size="20"
                :component="Pin32Filled"
                class="hover:cursor-pointer"
                :depth="isPinned ? 2 : 4"
                @click="pinWindow()"
            />
          </template>
        </n-tooltip>
        <n-dropdown
            :options="preferencesOptions"
            :show-arrow="false"
            trigger="click"
            size="medium"
            placement="top-end"
            @select="onSelectPreferenceMenu">
          <n-tooltip :delay="1000" :show-arrow="false">
            偏好设置
            <template #trigger>
              <n-icon
                  :size="20"
                  :component="Setting"
                  class="hover:cursor-pointer"
                  :depth="2"
              />
            </template>
          </n-tooltip>
        </n-dropdown>
      </div>
    </div>
  </div>
<!--  不显示按钮-->
  <ModalForm
      title=""
      v-model:show="showAboutModal"
      ref="showAboutModalRef"
      :show-icon="false"
      :loading="showAboutModalLoading"
      @submit="handleShowAboutModal"
      :schemas="showAboutModalSchemas"
      negativeText=""
      positiveText=""
  >
    <template #template>
      <n-space :size="10" :wrap="false" :wrap-item="false" align="center" vertical>
        <n-avatar :size="120" :src="appIcon ?? 'https://q1.qlogo.cn/g?b=qq&nk=10007&s=100'" color="#0000"></n-avatar>
        <div class="font-bold text-xl m-1">HostManager</div>
        <n-text>(v1.0.24)</n-text>
        <n-space :size="5" :wrap="false" :wrap-item="false" align="center">
          <n-text class="cursor-pointer hover:underline decoration-1" @click="onOpenSource">源码地址</n-text>
          <n-divider vertical />
          <n-text class="cursor-pointer hover:underline decoration-1" @click="onOpenWebsite">官方网站</n-text>
        </n-space>
        <div class="text-xs text-gray-500">
          Copyright © 2024 HostManager.glc All rights reserved
        </div>
      </n-space>
    </template>
  </ModalForm>

</template>

<style lang="scss" scoped>
.n-divider {
  margin: 0;
}
</style>