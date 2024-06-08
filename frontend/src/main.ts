// 引入全局样式
import './assets/css/style.css';

import {createApp} from 'vue'
import App from './App.vue'
import {setupStore} from '@/src/stores'
import {setupNaiveDiscreteApi} from '@/src/plugins'
import router from '@/src/router'


async function bootstrap() {
    const app = createApp(App)
    // 挂载路由
    app.use(router)
    // 挂载状态管理
    setupStore(app)
    // 挂载 naive-ui 脱离上下文的 Api
    setTimeout(() => {
        setupNaiveDiscreteApi()
    })
    // 挂载应用
    app.mount('#app', true)
}

void bootstrap()

