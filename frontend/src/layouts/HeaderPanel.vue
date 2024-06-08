<script setup lang="ts">
import {Component, computed, h, onMounted, ref} from "vue";
import {GetAppSettingAttr} from "@/wailsjs/go/controller/Config";
import {Quit, WindowHide, WindowMinimise, WindowToggleMaximise} from "@/wailsjs/runtime";
import WindowHidden from "@/src/components/Icons/WindowHidden.vue";
import WindowMin from "@/src/components/Icons/WindowMin.vue";
import WindowMax from "@/src/components/Icons/WindowMax.vue";
import WindowClose from "@/src/components/Icons/WindowClose.vue";
import WindowRestore from "@/src/components/Icons/WindowRestore.vue";
import MoonOutline from "@/src/components/Icons/MoonOutline.vue";
import SunnyOutline from "@/src/components/Icons/SunnyOutline.vue";
import ContrastOutline from "@/src/components/Icons/ContrastOutline.vue";
import {Home, LockClosed, LockOpen, LogOutOutline as LogoutIcon, Pencil as EditIcon} from '@vicons/ionicons5'
import {NIcon} from 'naive-ui'
import {useThemeStore} from "@/src/stores";
import {storeToRefs} from "pinia";
import {IsEncrypted, IsLogin, Logout} from "@/wailsjs/go/controller/Auth";
import {IsEncryptedResponse, IsLoginResponse, LogoutResponse} from "@/src/types/response/auth";
import {GetAppSettingAttrResponse} from "@/src/types/response/config";
import {useRouter} from "vue-router";
import {type FormSchema, ModalForm, type Recordable} from 'naive-ui-form'

// -----------change password modal start-----------
const showChangePasswordModal = ref(false)
const changePasswordSchemas: FormSchema[] = [
  {
    field: 'old_password',
    label: '老密码',
    type: 'input',
    required: true
  },
  {
    field: 'password',
    label: '新密码',
    type: 'input',
    required: true
  },
  {
    field: 'confirm_password',
    label: '确认密码',
    type: 'input',
    required: true,
  }
]

const changePasswordLoading = ref(false)

function handleSubmit(values: Recordable) {
  console.log(values)
  changePasswordLoading.value = true
  setTimeout(() => {
    changePasswordLoading.value = false
    showChangePasswordModal.value = false
  }, 2000)
}

// -----------change password modal end-----------


const router = useRouter()
const windowRestore = ref(false)

const themeStore = useThemeStore()
const {modeState} = storeToRefs(themeStore)


const renderIcon = (icon: Component, color: string = '') => {
  return () => {
    return h(NIcon, {
      color: color ?? 'var(--icon-color)',
    }, {
      default: () => h(icon)
    })
  }
}

const appName = ref<string>('')
const appIcon = ref<string>('')
const isEncrypted = ref<boolean>(false)
const isLogin = ref<boolean>(false)

const moreOptions = ref([
  {
    label: '回到首页',
    key: 'home',
    icon: renderIcon(Home)
  },
  {
    show: isLogin,
    label: '分割线0',
    key: 'divider0',
    type: 'divider'
  },
  {
    show: isLogin,
    disabled: true,
    label: isEncrypted ? '加密数据' : '明文数据',
    key: 'is_encrypted',
    icon: isEncrypted ? renderIcon(LockClosed, '#0e7a0d') : renderIcon(LockOpen, '#ffd700')
  },
  {
    show: isLogin,
    label: '修改密码',
    key: 'change_password',
    icon: renderIcon(EditIcon),
  },
  {
    show: isLogin,
    label: '退出登录',
    key: 'logout',
    icon: renderIcon(LogoutIcon)
  },
  {
    label: '分割线1',
    key: 'divider1',
    show: true,
    type: 'divider'
  },
  {
    label: '隐藏',
    key: 'hide',
    icon: renderIcon(WindowHidden)
  },
  {
    label: '最小化',
    key: 'minimise',
    icon: renderIcon(WindowMin)
  },
  {
    label: '最大化',
    key: 'maximise',
    icon: renderIcon(WindowMax)
  },
  {
    label: '分割线2',
    key: 'divider2',
    show: true,
    type: 'divider'
  },
  {
    label: '关闭',
    key: 'close',
    icon: renderIcon(WindowClose)
  }
])


const handleSelect = async (key: string | number) => {
  // $message.info(String(key))
  switch (key) {
    case 'home':
      await handleHome()
      break
    case 'change_password':
      await (async () => {
        showChangePasswordModal.value = true
      })();
      break
    case 'logout':
      await handleLogout()
      break
    case 'hide':
      await handleHide()
      break
    case 'minimise':
      await handleMinimise()
      break
    case 'maximise':
      await handleMaximise()
      break
    case 'close':
      await handleClose()
      break
  }
}

// 回到首页
const handleHome = async () => {
  await router.push('/')
}

const handleLogout = async () => {
  const response: LogoutResponse = await Logout()
  const {data} = response
  if (data.status) {
    $message.success('退出成功')
    await router.push('/login')
  } else {
    $message.error(data.message)
  }
}

const handleHide = async () => {
  WindowHide()
}

const handleMinimise = async () => {
  WindowMinimise()
}

const handleMaximise = async () => {
  windowRestore.value = !windowRestore.value
  WindowToggleMaximise()
}

const handleClose = async () => {
  Quit()
}

const modeStateText = computed(() => {
  switch (modeState.value) {
    case 'light':
      return '日间模式'
    case 'dark':
      return '夜间模式'
    case 'auto':
      return '跟随系统'
  }
})


onMounted(async () => {
  await GetAppSettingAttr('appName').then((response: GetAppSettingAttrResponse) => {
    appName.value = response.data.value;
  });

  await GetAppSettingAttr('appIcon').then((response: GetAppSettingAttrResponse) => {
    appIcon.value = response.data.value;
  });

  await IsEncrypted().then((response: IsEncryptedResponse) => {
    isEncrypted.value = response.data.encrypted;
  });

  await IsLogin().then((response: IsLoginResponse) => {
    isLogin.value = response.data.is_login;
  });
});


</script>

<template>
  <div class="flex flex-col h-full w-full">
    <!--    top-0 z-50-->
    <div class="flex flex-grow items-center">
      <div class="flex flex-grow" style="--wails-draggable:drag;user-select: none" @dblclick="handleMaximise()">
        <n-dropdown
            :options="moreOptions"
            :show-arrow="true"
            size="small"
            placement="top-end"
            trigger="click"
            @select="handleSelect">
          <n-image preview-disabled v-show="appIcon" :src="appIcon" class="w-7 h-7 mx-2 rounded-md" @dblclick.stop
                   alt="icon"/>

        </n-dropdown>
        <p class="text-lg font-bold">{{ appName }}</p>
      </div>
      <div class="flex items-center ml-auto h-full cursor-pointer ">
        <n-tooltip :delay="1000" :show-arrow="false">
          {{ modeStateText }}
          <template #trigger>
            <button
                class="focus:outline-none w-10 h-full flex items-center justify-center hover:shadow-md"
                @click="themeStore.toggleDarkMode">
              <moon-outline v-if="modeState === 'dark'" size="20"/>
              <sunny-outline v-if="modeState === 'light'" size="20"/>
              <contrast-outline v-if="modeState === 'auto'" size="20"/>
            </button>
          </template>
        </n-tooltip>
        <n-divider vertical class="m-2"/>
        <n-tooltip :delay="1000" :show-arrow="false">
          最小化
          <template #trigger>
            <button
                class="focus:outline-none w-10 h-full flex items-center justify-center hover:shadow-md"
                @click="handleMinimise()">
              <window-min size="20"/>
            </button>
          </template>
        </n-tooltip>
        <n-tooltip :delay="1000" :show-arrow="false">
          最大化
          <template #trigger>
            <button
                class="focus:outline-none w-10 h-full flex items-center justify-center hover:shadow-md"
                @click="handleMaximise()">
              <window-max v-if="!windowRestore" size="20"/>
              <window-restore v-else size="20"/>
            </button>
          </template>
        </n-tooltip>
        <n-tooltip :delay="1000" :show-arrow="false">
          关闭
          <template #trigger>
            <button
                class="focus:outline-none w-10 h-full flex items-center justify-center hover:shadow-md hover:bg-red-600"
                @click="Quit()">
              <window-close size="20"/>
            </button>
          </template>
        </n-tooltip>
      </div>
    </div>
    <n-divider/>
  </div>
  <ModalForm
      v-model:show="showChangePasswordModal"
      :schemas="changePasswordSchemas"
      title="修改密码"
      :loading="changePasswordLoading"
      @submit="handleSubmit"
  />
</template>

<style lang="scss" scoped>
.n-divider {
  margin: 0;
}
</style>