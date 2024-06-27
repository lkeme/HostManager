<script setup lang="ts">
import {ref, watch} from 'vue';
import {type FormSchema, ModalForm, type ModalFormInstance, type Recordable} from "naive-ui-form";
import {NewTreeOption} from "@/src/layouts/DashboardPanel/createData";
import {CreateFolderResponse, DeleteFolderResponse, RenameFolderResponse} from "@/src/types/response/filesystem";
import {CreateFolder, DeleteFolder, RenameFolder} from "@/wailsjs/go/controller/FileSystem";

const props = defineProps<{
  title?: '新增分类' | '删除分类' | '重命名分类';
  show: boolean;
  option: NewTreeOption | null;
}>();

const modalRef = ref<ModalFormInstance | null>(null);

const internalShow = ref(false);

const folderSchemas = ref<Array<FormSchema>>([
  {
    // label: '分类名称',
    // required: true,
    field: 'folderName',
    type: 'input',
    componentProps: {
      placeholder: '',
      clearable: true
    },
    rules: [
      {required: true, message: '请输入分类名称', trigger: 'blur'}
    ]
  }
]);

const emits = defineEmits(['update:show', 'submit']);

const loading = ref(false);

const handleSubmit = async (values: Recordable) => {
  console.log(values.folderName, props.option)
  //
  loading.value = true;
  try {
    let response: CreateFolderResponse | RenameFolderResponse | DeleteFolderResponse;
    switch (props.title) {
      case '新增分类':
        response = await CreateFolder(values.folderName, props.option?.id ?? 0);
        break;
      case '删除分类':
        response = await DeleteFolder(props.option!.id);
        break;
      case '重命名分类':
        response = await RenameFolder(values.folderName, props.option!.id);
        break;
    }
    //
    const {data} = response;
    if (data.status) {
      $message.success(`${props.title}成功`);
    } else {
      $message.error(`${props.title}失败, ${data.message}` || '${props.title}失败');
    }
    emits('submit', data.status);
  } finally {
    loading.value = false;
  }
};

// 监听来自父组件的 show 变化，更新内部的显示状态
watch(() => props.show, (newVal) => {
  internalShow.value = newVal;
  if (newVal) {
    // folderSchemas.value[0].defaultValue = newVal.label;
    if (props.title === '重命名分类' || props.title === '删除分类') {
      folderSchemas.value[0].defaultValue = props.option.label;
      folderSchemas.value[0].componentProps.disabled = false;
    }
    if (props.title === '删除分类') {
      folderSchemas.value[0].componentProps.disabled = true;
    }
    if (props.title === '新增分类') {
      folderSchemas.value[0].defaultValue = '';
      folderSchemas.value[0].componentProps.disabled = false;
    }
  }
});

// 监听内部的显示状态，更新父组件的 show 状态
watch(() => internalShow.value, (newVal) => {
  emits('update:show', newVal);
});
</script>

<template>
  <ModalForm
      ref="modalRef"
      v-model:show="internalShow"
      @submit="handleSubmit"
      :title="title"
      :schemas="folderSchemas"
      :loading="loading"
      :closable="false"
      :show-icon="false"
      :negative-button-props="{ focusable: false, size: 'medium' }"
      :positive-button-props="{ focusable: false, size: 'medium', type: 'error' }"
  />
</template>
