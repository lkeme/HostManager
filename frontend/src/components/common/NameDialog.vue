<script setup>
import {computed, reactive, ref, watchEffect} from 'vue'
import {isEmpty} from 'lodash-es'

/**
 * Dialog for create or rename group
 */

const props = defineProps({
  modelValue: {
    type: Boolean,
    required: true,
    default: false
  },
})

const emit = defineEmits(['update:modelValue', 'confirm']);

const visible = computed({
  get() {
    return props.modelValue;
  },
  set(visible) {
    emit('update:modelValue', visible);
  }
});

const editGroup = ref('')
const groupForm = reactive({
  name: '',
})
const groupFormRef = ref(null)

const formRules = computed(() => {
  return {
    name: [
      {required: true, message: '请输入群组名', trigger: 'blur'},
      {max: 20, message: '群组名最多20个字符', trigger: 'blur'}
    ]
  }
})

const isRenameMode = computed(() => !isEmpty(editGroup.value))

watchEffect(() => {

})

const onConfirm = async () => {

}

const onClose = () => {
  visible.value = false
}


</script>

<template>
  <n-modal
      v-model:show="visible"
      :closable="false"
      :mask-closable="false"
      :negative-button-props="{ size: 'medium' }"
      negative-text="取消"
      :positive-button-props="{ size: 'medium' }"
      positive-text="确认"
      :show-icon="false"
      :title="isRenameMode ? '重命名' : '新建'"
      close-on-esc
      preset="dialog"
      transform-origin="center"
      @esc="onClose"
      @positive-click="onConfirm"
      @negative-click="onClose">
    <n-form
        ref="groupFormRef"
        :model="groupForm"
        :rules="formRules"
        :show-label="false"
        :show-require-mark="false"
        label-placement="top">
      <n-form-item label="群组名" path="name" required>
        <n-input v-model:value="groupForm.name" placeholder=""/>
      </n-form-item>
    </n-form>
  </n-modal>
</template>

<style lang="scss" scoped></style>