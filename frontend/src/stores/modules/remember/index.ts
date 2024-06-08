import {defineStore} from "pinia";

const id = "remember_info";

export const useRememberStore = defineStore(id, {
    state: () => ({
        id: id,
        remember: '',
    }),
    getters: {
        getRemember: (state) => state.remember,
    },
    actions: {
        setRememberInfo(val: any) {
            localStorage.setItem(id, val);
        },
        getRememberInfo(): string {
            return localStorage.getItem(id) || '';
        },
        removeRememberInfo() {
            localStorage.removeItem(id);
        },
    },
})