<script lang="ts" setup>
import {Base64} from 'js-base64'
import {onMounted, reactive, ref} from 'vue';
import {useRoute, useRouter} from 'vue-router';
import {LockClosedOutline} from '@vicons/ionicons5';
import {IsRegister, Login, Register} from "@/wailsjs/go/controller/Auth";
import {useRememberStore} from '@/src/stores';
import {IsRegisterResponse, LoginResponse, RegisterResponse} from "@/src/types/response/auth";


// --------------registerForm start----------------
const registerFormRef = ref();
const registerFormInline = reactive({
  password: '',
  confirm_password: '',
});

const handleRegisterSubmit = (e: MouseEvent) => {
  const next = async () => {
    $message.loading('尝试注册中...');
    loading.value = true;

    try {
      const response: RegisterResponse = await Register(registerFormInline.password, registerFormInline.confirm_password);
      const {data} = response;
      $message.destroyAll();
      //
      if (data.status) {
        const toPath = decodeURIComponent((route.query?.redirect || '/') as string);
        $message.success('注册成功，即将进入系统');
        //
        if (route.name === '/login') {
          await router.replace('/');
        } else {
          await router.replace(toPath)
        }
      } else {
        $message.error(data.message || '注册失败');
      }
    } finally {
      loading.value = false;
    }
  };

  e.preventDefault();
  registerFormRef.value?.validate(async (errors) => {
    if (errors) {
      $message.error('注册信息填写有误，请检查后重新填写');
      return;
    }
    //
    $dialog.warning({
      title: '注册提示',
      content: '注册密码非常重要，如果密码丢失，数据不可找回，请牢记您的密码',
      negativeText: '再想想',
      positiveText: '继续注册？',
      onPositiveClick: next,
      onNegativeClick: () => {
        return;
      }
    });
  });
};

// --------------registerForm end----------------

// --------------loginForm start----------------
const loginFormRef = ref();
const loginFormInline = reactive({
  password: '',
  remember: false,
});

const handleLoginSubmit = (e: MouseEvent) => {
  e.preventDefault();
  loginFormRef.value?.validate(async (errors) => {
    if (errors) {
      $message.error('登录信息填写有误，请检查后重新填写');
      return;
    }
    //
    $message.loading('尝试登录中...');
    loading.value = true;

    try {
      const response: LoginResponse = await Login(loginFormInline.password);
      const {data} = response;
      $message.destroyAll();
      //
      if (data.status) {
        saveLoginInfo();
        const toPath = decodeURIComponent((route.query?.redirect || '/') as string);
        $message.success('登录成功');
        //
        if (route.name === '/login') {
          await router.replace('/');
        } else {
          await router.replace(toPath)
        }
      } else {
        $message.error(data.message || '登录失败');
      }
    } finally {
      loading.value = false;
    }
  });
};

const remember = useRememberStore();
const saveLoginInfo = () => {
  if (loginFormInline.remember) {
    remember.setRememberInfo(Base64.encode(JSON.stringify(loginFormInline)));
    return;
  }
  remember.removeRememberInfo();
}

const readerLoginInfo = () => {
  const rememberInfo = remember.getRememberInfo()
  if (rememberInfo) {
    const info = JSON.parse(Base64.decode(rememberInfo));
    loginFormInline.password = info.password;
    loginFormInline.remember = info.remember;
  }
}
// --------------loginForm end----------------


const loading = ref(false);

const router = useRouter();
const route = useRoute();


// --------------common start----------------
const rules = {
  password: [
    {required: true, message: '请输入密码', trigger: 'blur'},
    {min: 8, max: 16, message: '密码长度应为8-16', trigger: 'blur'},
    {pattern: /^(?![0-9]+$)(?![a-zA-Z]+$)[0-9A-Za-z]{8,16}$/, message: '密码必须包含数字和字母', trigger: 'blur'},
  ],
  confirm_password: [
    {required: true, message: '请确认密码', trigger: 'blur'},
    {
      validator: (_, value, callback) => {
        if (value === registerFormInline.password) {
          callback();
        } else {
          callback(new Error('两次输入密码不匹配'));
        }
      }, trigger: 'blur'
    },
  ],
};

// 判断是否注册
const isRegister = ref(false);
onMounted(() => {
  IsRegister().then((res: IsRegisterResponse) => {
    isRegister.value = res.data.is_register;
  });
  readerLoginInfo();
});
// --------------common end----------------


</script>

<template>
  <div class="flex flex-col h-full w-full view-auth-bg">
    <div class="flex-1 px-8 py-4 max-w-[384px] min-w-[320px] mx-auto">
      <div class="py-8 text-center">
        <div class="">
          <img src="@/src/assets/images/logo.svg" alt="logo" class="mx-auto"/>
        </div>
        <!--          <div class="text-sm text-[#808695]">这里是描述~~~~~</div>-->
        <div class="text-sm text-red-500" v-if="!isRegister">当前还未注册</div>
      </div>
      <div class="view-auth-register-form" v-if="!isRegister">
        <n-form ref="registerFormRef" label-placement="left" size="large" :model="registerFormInline" :rules="rules">
          <n-form-item path="password">
            <n-input v-model:value="registerFormInline.password" type="password" showPasswordOn="click"
                     placeholder="请输入密码">
              <template #prefix>
                <n-icon size="18" color="#808695">
                  <LockClosedOutline/>
                </n-icon>
              </template>
            </n-input>
          </n-form-item>
          <n-form-item path="confirm_password">
            <n-input v-model:value="registerFormInline.confirm_password" type="password" showPasswordOn="click"
                     placeholder="请再次输入密码">
              <template #prefix>
                <n-icon size="18" color="#808695">
                  <LockClosedOutline/>
                </n-icon>
              </template>
            </n-input>
          </n-form-item>
          <n-form-item>
            <n-button type="primary" @click="handleRegisterSubmit" size="large" :loading="loading" block>
              注册
            </n-button>
          </n-form-item>
        </n-form>
      </div>
      <div class="view-auth-login-form" v-else>
        <n-form
            ref="loginFormRef"
            label-placement="left"
            size="large"
            :model="loginFormInline"
            :rules="rules"
        >
          <n-form-item path="password">
            <n-input
                v-model:value="loginFormInline.password"
                type="password"
                showPasswordOn="click"
                placeholder="请输入密码"
            >
              <template #prefix>
                <n-icon size="18" color="#808695">
                  <LockClosedOutline/>
                </n-icon>
              </template>
            </n-input>
          </n-form-item>
          <n-form-item>
            <n-checkbox v-model:checked="loginFormInline.remember">记住密码</n-checkbox>
          </n-form-item>
          <n-form-item>
            <n-button type="primary" @click="handleLoginSubmit" size="large" :loading="loading" block>
              登录
            </n-button>
          </n-form-item>
        </n-form>
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
@media (min-width: 768px) {
  .view-auth-bg {
    background-image: url('@/src/assets/images/login.svg');
    background-repeat: no-repeat;
    background-position: 50%;
    background-size: 100%;
  }
}
</style>