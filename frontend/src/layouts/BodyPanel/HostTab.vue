<template>
  <div class="flex flex-col h-full" >
    <div class="flex h-1/6 justify-between ">
      <n-button quaternary>
        <template #icon>
          <n-icon
              :size="30"
              :component="SlideAdd48Filled"
              class="hover:cursor-pointer"
              :depth="2"
          />
        </template>
      </n-button>

      <button :disabled="addFolderDisabled" @click="folderNameDialogVisible = true"
              class="text-white px-2 my-2 rounded-full" :class="{
          'bg-purple-500': !addFolderDisabled,
          'bg-gray-300': addFolderDisabled,
          'hover:bg-purple-700': !addFolderDisabled
        }">
        添加
        .vue
      </button>
      <button :disabled="addDocumentDisabled" @click="documentNameDialogVisible = true"
              class="text-white px-2 my-2 rounded-full" :class="{
          'bg-purple-500': !addDocumentDisabled,
          'bg-gray-300': addDocumentDisabled,
          'hover:bg-purple-700': !addDocumentDisabled
        }">
        测试
      </button>
      <button :disabled="deleteDisabled" @click="deleteItem" class="text-white px-2 my-2 rounded-full" :class="{
        'bg-purple-500': !deleteDisabled,
        'bg-gray-300': deleteDisabled,
        'hover:bg-purple-700': !deleteDisabled
      }">
        删除
      </button>
      <button :disabled="renameDisabled" @click="renameDialogVisible = true" class="text-white px-2 my-2 rounded-full"
              :class="{
          'bg-purple-500': !renameDisabled,
          'bg-gray-300': renameDisabled,
          'hover:bg-purple-700': !renameDisabled
        }">
        重命名
      </button>
      <NameDialog v-model="folderNameDialogVisible" title="フォルダの追加" placeholder="フォルダ名"
                  description="フォルダ名を入力してください。"
                  @confirm="handleFolderAdded"/>
      <NameDialog v-model="documentNameDialogVisible" title="文書の追加" placeholder="文書名"
                  description="文書名を入力してください。"
                  @confirm="handleDocumentAdded"/>
      <NameDialog v-model="renameDialogVisible" title="名前の変更" placeholder="新しい名前"
                  description="新しい名前を入力してください。"
                  @confirm="handleRename"/>
    </div>
    <GTree v-model:data="data" @selected="selectedHandler"/>

  </div>
</template>

<script setup lang="ts">
import {computed, onMounted, reactive, ref, toRaw, watch} from 'vue';
import GTree from '@/src/components/GTree/GTree.vue';
import {
  addDocumentDataById,
  addFolderDataById,
  createData,
  removeDataById,
  renameDataById
} from '@/src/components/GTree/createData';

import NameDialog from '@/src/components/common/NameDialog.vue';

import {NIcon, TreeOption} from 'naive-ui';
import SlideAdd48Filled from "@/src/components/Icons/SlideAdd48Filled.vue";


const folderNameDialogVisible = ref(false);
const documentNameDialogVisible = ref(false);
const renameDialogVisible = ref(false);
const currentId = ref('');
const currentLabel = ref('');
const isFolder = ref(false);

let data: TreeOption[] = reactive([]);

const selectedHandler = ({id, label, isFolder: isFolderValue}) => {
  currentId.value = id;
  currentLabel.value = label;
  isFolder.value = isFolderValue;
};

const handleFolderAdded = folderName => {
  addFolderDataById(data, currentId.value, folderName);
};

const handleDocumentAdded = documentName => {
  addDocumentDataById(data, currentId.value, documentName);
};

const handleRename = newName => {
  renameDataById(data, currentId.value, newName);
};

const deleteItem = () => {
  removeDataById(data, currentId.value);
};

// [フォルダを追加]が使用不可
const addFolderDisabled = computed(() => {
  if (isFolder.value) return false;
  return !!currentId.value;
});

// [文書を追加]が使用不可
const addDocumentDisabled = computed(() => !isFolder.value);

// [削除]が使用不可
const deleteDisabled = computed(() => {
  if (isFolder.value && currentLabel.value === '未分類') return true;
  return !!!currentId.value;
});

// [改名]が使用不可
const renameDisabled = computed(() => {
  if (isFolder.value && currentLabel.value === '未分類') return true;
  return !!!currentId.value;
});

const loadFromLocalStorage = () => {
  const item = localStorage.getItem('data');
  if (item) {
    const parseData = reactive(JSON.parse(item));
    data.splice(0, data.length, ...parseData);
  } else {
    // data.splice(0, data.length, ...createData());
    data.splice(0, data.length, ...createData());
  }
};

const saveToLocalStorage = () => {
  const rawData = toRaw(data);
  localStorage.setItem('data', JSON.stringify(rawData));
};

onMounted(() => {
  loadFromLocalStorage();
});

watch(data, async data => {
  saveToLocalStorage();
});
</script>

<style lang="scss" scoped>
.n-card {
  ::v-deep(.n-card-header) {
    padding: 2px;
  }
}
</style>