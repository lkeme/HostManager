import type {RouteRecordRaw} from "vue-router";

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

export default routes