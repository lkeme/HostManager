import type {Router, RouteRecordRaw} from 'vue-router'
import {createRouter, createWebHashHistory} from 'vue-router'
import {IsLogin} from "@/wailsjs/go/controller/Auth";
import {IsLoginResponse} from "@/src/types/response/auth";

const routes: RouteRecordRaw[] = [
    {
        path: '/',
        name: 'Dashboard',
        component: () => import('@/src/views/dashboard/index.vue'), // 注意这里要带上 文件后缀.vue
        meta: {
            title: 'Index/Dashboard',
            requiresAuth: true
        }
    },
    // {
    //     path: '/index',
    //     name: 'Dashboard',
    //     component: () => import('@/src/views/dashboard/index.vue'), // 注意这里要带上 文件后缀.vue
    //     meta: {
    //         title: 'Index/Dashboard',
    //         requiresAuth: true
    //     }
    // },
    {
        path: '/login',
        name: 'Login',
        component: () => import('@/src/views/login/index.vue'), // 注意这里要带上 文件后缀.vue
        meta: {
            title: 'Login',
            requiresAuth: false
        }
    },
    {
        name: 'nofound',
        path: "/:nofound(.*)*",
        redirect: '/index',
        meta: {
            title: '404',
            requiresAuth: false
        }
    }
]

// 创建路由
const router: Router = createRouter({
    // history: createWebHistory(),
    history: createWebHashHistory(),
    routes,
})

// 设置导航守卫 有token 直接进入首页或指定的页面（无法进入登录页面）
// router.beforeEach(async (to, _, __) => {
router.beforeEach(async (to, _) => {
    // instead of having to check every route record with
    // to.matched.some(record => record.meta.requiresAuth)
    const res: IsLoginResponse = await IsLogin()
    console.log(res)

    if (to.meta.requiresAuth === true && !res.data.is_login) {
        return {
            path: '/login',
            query: {redirect: to.fullPath},
        }
    } else {
        if (to.meta.requiresAuth === 'guest' && res.data.is_login) {
            return {
                path: '/',
                query: {redirect: to.fullPath},
            }
        }
    }

    // if (from.path === '/' && from.matched && from.matched.length === 0) {
    //
    // }
    // next();
})


// 设置页面标题
router.afterEach((to) => {
    if (to.meta.title) {
        document.title = to.meta.title as string;
    }
})

// 导出路由
export default router