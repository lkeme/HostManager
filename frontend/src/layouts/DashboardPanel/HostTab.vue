<script setup lang="ts">
import {computed, h, onMounted, ref, watch, withModifiers} from 'vue';
import {NButton, NDivider, NIcon, useOsTheme} from 'naive-ui';
import {
  AddOutline,
  Folder,
  FolderOpen,
  FolderOpenOutline,
  FolderOutline,
  PencilOutline,
  RefreshCircleOutline,
  TrashOutline,
} from '@vicons/ionicons5';
import SearchOutline from '@/src/components/Icons/SearchOutline.vue';
import Loader from '@/src/components/Icons/Loader.vue';
import {findTreeOptionLevel, NewTreeOption} from '@/src/layouts/DashboardPanel/createData';
import {type FormSchema, ModalForm, type ModalFormInstance, type Recordable} from 'naive-ui-form'
import {CreateFolder, GetAllFolders} from "@/wailsjs/go/controller/FileSystem";
import {CreateFolderResponse, FolderHierarchy, GetAllFoldersResponse} from "@/src/types/response/filesystem";

// const treeOptions: Ref<Array<TreeOption>> = ref<Array<TreeOption>>([]);
const treeOptions = ref([])

const showIrrelevantNodes = ref(false);
const pattern = ref('');


// Naive UI / 暗黑模式
const osThemeRef = useOsTheme();
const isDark = computed(() => (osThemeRef.value === 'dark'));

const styleObject = computed(() => ({
  // backgroundColor: isDark.value ? '#464646' : '#f2f2f2',
  '--n-node-text-color': isDark.value ? '#fff' : '#555',
  '--n-node-color-hover': isDark.value ? '#555' : '#e3e3e3',
  '--n-node-color-active': isDark.value ? '#FF8F70' : '#FFE3DB',
  'height': `100%`,
}));

// 最初的节点
const expandedKeys = ref([]);


const handleExpandedKeysUpdate = (keys, options, {node, action}) => {
  expandedKeys.value = keys;

  console.log('Updated Expanded Keys:', keys);
  console.log('Options:', options);
  console.log('Node:', node);
  console.log('Action:', action);

  // 必要に応じて追加の処理を行う
  if (action === 'expand') {
    console.log(`${node.label} has been expanded!`);
  } else if (action === 'collapse') {
    console.log(`${node.label} has been collapsed!`);
  }
};

const loadFromLocalStorage = () => {
  const item = localStorage.getItem('expandedKeys');
  if (item) {
    expandedKeys.value = JSON.parse(item);
  }
};

const saveToLocalStorage = () => {
  const keys = localStorage.setItem('expandedKeys', JSON.stringify(expandedKeys.value));
};


watch(expandedKeys, async (newKeys) => {
  saveToLocalStorage();
});

const handleDragEnd = (event) => {
  $message.info('节点拖拽');
  console.log('Drag End:', event);
};


// -----------------nodeHandle start-----------------
const nodeHoverDefaultKey = ref<string | number>(null);
const nodeClickDefaultOption = ref<NewTreeOption | null>(null)

const saveSelectedNode = (option: NewTreeOption) => {
  // 保存当前点击的节点
  nodeClickDefaultOption.value = (option.key === nodeClickDefaultOption.value?.key) ? null : option;
  console.log('选中事件', option.key);
}

// const emit = defineEmits(['selected']);
const nodeProps = ({option}: { option: NewTreeOption }) => {
  return {
    onClick() {
      saveSelectedNode(option)
      const isFolder = typeof option.children !== 'undefined';
      // emit('selected', {id: option.key, label: option.label, isFolder});
      console.log('单击事件', option.key);
    },
    ondblclick() {
      console.log('双击事件', option.key);
    },
    onMouseenter() {
      nodeHoverDefaultKey.value = option.key;
    },
    onMouseleave() {
      nodeHoverDefaultKey.value = '';
    },
  }
}

// 渲染节点后缀
const renderNodeSuffix = ({option}: { option: NewTreeOption }) => {
  // 如果不是当前hover的节点，不渲染
  if (option.key !== nodeHoverDefaultKey.value) return;
  // 如果是文件夹，且有子节点，只渲染编辑图标按钮
  if (option.children) {
    return h('div', {
      style: 'display: flex; align-items: center; justify-content: center;',
    }, [
      h(NButton, {
        style: 'margin-left: 4px',
        text: true,
        size: 'medium',
        onClick: withModifiers(() => {
          console.log('编辑', option.key);
        }, ['stop', 'prevent'])

      }, {default: () => h(NIcon, {color: '#2ed6e9'}, {default: () => h(PencilOutline)}),})
    ])
  } else {
    return h('div', {
      style: 'display: flex; align-items: center; justify-content: center;',
    }, [
      h(NButton, {
        style: 'margin-left: 4px',
        text: true,
        size: 'medium',
        onClick: withModifiers(() => {
          console.log('编辑', option.key);
        }, ['stop', 'prevent'])

      }, {default: () => h(NIcon, {color: '#2ed6e9'}, {default: () => h(PencilOutline)}),}),
      h(NButton, {
        text: true,
        style: 'margin-left: 4px',
        size: 'medium',
        onClick: withModifiers(() => {
          console.log('删除', option.key);
        }, ['stop', 'prevent'])
      }, {default: () => h(NIcon, {color: '#fb0e0e'}, {default: () => h(TrashOutline)}),}),
    ])
  }

}

// 渲染展开开关的图标
const renderSwitcherIconWithExpand = ({option, expanded}: { option: NewTreeOption, expanded: boolean }) => {
  return h(NIcon, null, {
    default: () => {
      if (typeof option.children !== 'undefined' && option.children.length > 0) {
        return h(expanded ? FolderOpen : Folder)
      }
      return h(expanded ? FolderOpenOutline : FolderOutline)
    }
  });
}
// -----------------nodeHandle end-----------------

// -----------refreshData start----------------
const refreshDataLoading = ref(false);

const convertToTreeOption = (folderHierarchy: FolderHierarchy): NewTreeOption => {
  const {folder, children} = folderHierarchy;
  return {
    id: folder.ID,
    level: folder.level,
    parentID: folder.parentID,
    // isLeaf: !children || children.length === 0, // 添加对children为null或undefined的检查
    key: folder.key,
    label: folder.name,
    children: children ? children.map(convertToTreeOption) : [] // 添加对children为null或undefined的检查
  };
}

const convertFolderHierarchyToTreeOptions = (folderHierarchies: Array<FolderHierarchy>): Array<NewTreeOption> => {
  return folderHierarchies.map(convertToTreeOption);
}

const handleRefreshData = async () => {
  refreshDataLoading.value = true
  try {
    const response: GetAllFoldersResponse = await GetAllFolders();
    const {data} = response;
    console.log(data)
    if (data.status) {
      treeOptions.value = convertFolderHierarchyToTreeOptions(data.folderHierarchy);
      $message.success('加载数据成功');
    } else {
      $message.error(`加载数据失败, ${data.message}` || '加载数据失败，请刷新重试');
    }
  } finally {
    refreshDataLoading.value = false;
  }
};
// -----------refreshData end----------------

// -----------addFolder start----------------
const showAddFolderModal = ref(false);
const addFolderModalRef = ref<ModalFormInstance | null>(null);
const addFolderLoading = ref(false);
const addFolderSchemas: FormSchema[] = [
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
];
const handleAddFolder = async (values: Recordable) => {
  addFolderLoading.value = true;
  try {
    const response: CreateFolderResponse = await CreateFolder(values.folderName, nodeClickDefaultOption.value?.id || 0, nodeClickDefaultOption.value?.key.toString() || 'ROOT');
    const {data} = response;
    //
    if (data.status) {
      $message.success('新增分类成功');
    } else {
      $message.error(`新增分类失败, ${data.message}` || '新增分类失败');
    }
  } finally {
    addFolderLoading.value = false;
    showAddFolderModal.value = false;
    await handleRefreshData();
  }
};
// -----------addFolder end----------------


// -----------addFile start----------------
const handleAddFileStatus = computed(() => {
  return treeOptions.value.length === 0;
})

const handleAddFile = () => {
  if (!nodeClickDefaultOption.value) {
    $message.warning('请先选择一个分类');
    return;
  }
  $message.info('新增会话');
};
// -----------addFile end----------------


// -----------batchDelete start----------------
// 批量删除按钮
const deleteButtonLoading = ref(false);
const deleteCheckAbleEnable = ref(false);
const handleDeleteButton = () => {
  deleteCheckAbleEnable.value = !deleteCheckAbleEnable.value;
  deleteButtonLoading.value = !deleteButtonLoading.value;
};
// -----------batchDelete end----------------

onMounted(async () => {
  await handleRefreshData();
});
</script>


<template>
  <ModalForm
      title="新增分类"
      v-model:show="showAddFolderModal"
      ref="addFolderModalRef"
      :schemas="addFolderSchemas"
      :loading="addFolderLoading"
      @submit="handleAddFolder"
      :closable="false"
      :show-icon="false"
      :negative-button-props="{ focusable: false, size: 'medium' }"
      :positive-button-props="{ focusable: false, size: 'medium',type:'error'}"
  />

  <div class="flex flex-col flex-1 p-2 h-full gap-3 g-tree">
    <!--搜索框-->
    <n-input v-model:value="pattern" placeholder="输入搜索的内容">
      <template #prefix>
        <n-icon :component="SearchOutline" size="20"/>
      </template>
      <template #suffix>
        <n-tooltip trigger="hover" placement="top">
          <template #trigger>
            <n-switch size="small" v-model:value="showIrrelevantNodes"/>
          </template>
          <template #default>
            {{ showIrrelevantNodes ? '展示搜索无关的节点' : '隐藏搜索无关的节点' }}
          </template>
        </n-tooltip>
      </template>
    </n-input>
    <!--操作按钮-->
    <n-space align="stretch" justify="space-between" :wrap="false">
      <n-tooltip placement="top-start" trigger="hover">
        <template #trigger>
          <n-button type="success" size="tiny" @click="handleRefreshData" :loading="refreshDataLoading">
            <template #icon>
              <RefreshCircleOutline/>
            </template>
            刷新
          </n-button>
        </template>
        <span>刷新数据</span>
      </n-tooltip>
      <n-tooltip placement="top-start" trigger="hover">
        <template #trigger>
          <n-button type="info" size="tiny" @click="showAddFolderModal = true;" :disabled="nodeClickDefaultOption?.level > 1" :loading="addFolderLoading">
            <template #icon>
              <AddOutline/>
            </template>
            分类
          </n-button>
        </template>
        <span>最大支持二级分类</span>
      </n-tooltip>
      <n-button type="primary" size="tiny" @click="handleAddFile" :disabled="handleAddFileStatus">
        <template #icon>
          <AddOutline/>
        </template>
        会话
      </n-button>
      <n-button type="error" size="tiny" @click="handleDeleteButton">
        <template #icon>
          <Loader v-if="deleteButtonLoading"/>
          <TrashOutline v-else/>
        </template>
        批量
      </n-button>
    </n-space>
    <!--分割线-->
    <n-divider/>
    <!--详情-->
    <n-tooltip placement="bottom-start" trigger="hover">
      <template #trigger>
        <n-tag disabled>
          已选择分类: {{ nodeClickDefaultOption?.label ?? 'ROOT' }}
        </n-tag>
      </template>
      <n-space vertical>
        <n-tag type="success">
          Label: {{ nodeClickDefaultOption?.label ?? 'ROOT' }}
        </n-tag>
        <n-tag type="success">
          Key: {{ nodeClickDefaultOption?.key ?? 'ROOT' }}
        </n-tag>
        <n-tag type="success">
          Level: {{
            nodeClickDefaultOption?.key ? findTreeOptionLevel(treeOptions, nodeClickDefaultOption) : '0'
          }}
        </n-tag>
        <n-tag type="success">
          IsFolder: {{ typeof nodeClickDefaultOption?.children !== 'undefined' ? '☑' : '☒' }}
        </n-tag>
      </n-space>
    </n-tooltip>
    <!--分割线-->
    <n-divider/>
    <!--列树-->
    <n-scrollbar class="max-h-full overflow-y-auto">
      <n-tree
          :indent="20"
          class="tree"
          cascade
          draggable
          :checkable="deleteCheckAbleEnable"
          :on-dragend="handleDragEnd"
          :show-irrelevant-nodes="showIrrelevantNodes"
          :pattern="pattern"
          block-line
          :data="treeOptions"
          :default-expanded-keys="expandedKeys"
          :on-update:expanded-keys="handleExpandedKeysUpdate"
          :render-switcher-icon="renderSwitcherIconWithExpand"
          checkbox-placement="right"
          :selectable="true"
          :class="{ 'is-dark': isDark }"
          virtual-scroll
          :style="styleObject"
          :render-suffix="renderNodeSuffix"
          :node-props="nodeProps"

      />
    </n-scrollbar>

  </div>

</template>


<style lang="scss" scoped>
.n-divider {
  margin: 0;
}

.tree svg path {
  stroke: #555;
}

.tree.is-dark svg path {
  stroke: #fff;
}
</style>


