<script setup lang="ts">
import {computed, h, onMounted, ref, watch, withModifiers} from 'vue';
import {NButton, NDivider, NIcon, useOsTheme} from 'naive-ui';
import {findTreeOptionLevel, NewTreeDropInfo, NewTreeOption} from '@/src/layouts/DashboardPanel/createData';
import {type FormSchema, ModalForm, type ModalFormInstance, type Recordable} from 'naive-ui-form'
import {CreateFolder, GetAllFolders} from "@/wailsjs/go/controller/FileSystem";
import {CreateFolderResponse, FolderHierarchy, GetAllFoldersResponse} from "@/src/types/response/filesystem";
import IconButton from "@/src/components/common/IconButton.vue";
import AddGroup from "@/src/components/Icons/AddGroup.vue";
import AddLink from "@/src/components/Icons/AddLink.vue";
import SearchOutline from '@/src/components/Icons/SearchOutline.vue';
import Refresh from '@/src/components/Icons/Refresh.vue';
import Loader from '@/src/components/Icons/Loader.vue';
import {
  Folder,
  FolderOpen,
  FolderOpenOutline,
  FolderOutline,
  Link,
  SettingsOutline,
  TrashOutline,
} from '@vicons/ionicons5';

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

const handleDragEnd = (node: NewTreeOption, event: DragEvent) => {
  $message.info('节点拖拽');
  console.log('Drag End:', event);
  console.log('Node:', node);
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
const createButton = (icon, color, text, onClick) => {
  return h(NButton, {
    style: 'margin-left: 4px',
    text: true,
    size: 'medium',
    onClick: withModifiers(onClick, ['stop', 'prevent'])
  }, {default: () => h(NIcon, {color}, {default: () => h(icon)})});
};

const getNodeMenu = (isEdit = false, isDelete = false) => {
  const style = 'display: flex; align-items: center; justify-content: center;';
  const buttons = [];

  if (isEdit) {
    buttons.push(createButton(SettingsOutline, '#2ed6e9', '编辑', () => {
      console.log('编辑', nodeHoverDefaultKey.value);
    }));
  }

  if (isDelete) {
    buttons.push(createButton(TrashOutline, '#fb0e0e', '删除', () => {
      console.log('删除', nodeHoverDefaultKey.value);
    }));
  }

  return h('div', {style}, buttons);
};

const renderNodeSuffix = ({option}: { option: NewTreeOption }) => {
  // 如果不是当前hover的节点，不渲染
  if (option.key !== nodeHoverDefaultKey.value) return;
  // 如果是文件夹，且有子节点，只渲染编辑图标按钮
  if (option.children && option.children.length > 0) {
    return getNodeMenu(true, false);
  } else {
    return getNodeMenu(true, true);
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
const addFolderBtnStatus = computed(() => {
  return nodeClickDefaultOption.value?.level > 1;
})

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
const addFileBtnStatus = computed(() => {
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
  // if (!deleteCheckAbleEnable.value) {
  //   $message.warning('请选择需要删除的节点');
  // }

  deleteCheckAbleEnable.value = !deleteCheckAbleEnable.value;
  deleteButtonLoading.value = !deleteButtonLoading.value;
};
const defaultCheckedKeys = ref([]);
const updateCheckedKeys = (keys: Array<any>) => {
  defaultCheckedKeys.value = keys;
  console.log('Checked Keys:', keys);
};

// -----------batchDelete end----------------

onMounted(async () => {
  await handleRefreshData();
});


const findSiblingsAndIndex = (node: NewTreeOption, nodes?: NewTreeOption[]): [NewTreeOption[], number] | [null, null] => {
  if (!nodes) return [null, null]
  for (let i = 0; i < nodes.length; ++i) {
    const siblingNode = nodes[i]
    if (siblingNode.key === node.key) return [nodes, i]
    const [siblings, index] = findSiblingsAndIndex(node, siblingNode.children)
    if (siblings && index !== null) return [siblings, index]
  }
  return [null, null]
}

const handleDrop = ({node, dragNode, dropPosition}: NewTreeDropInfo) => {
  const [dragNodeSiblings, dragNodeIndex] = findSiblingsAndIndex(dragNode, treeOptions.value)
  if (dragNodeSiblings === null || dragNodeIndex === null) return
  dragNodeSiblings.splice(dragNodeIndex, 1)
  if (dropPosition === 'inside') {
    if (node.children) {
      node.children.unshift(dragNode)
    } else {
      node.children = [dragNode]
    }
  } else if (dropPosition === 'before') {
    const [nodeSiblings, nodeIndex] = findSiblingsAndIndex(node, treeOptions.value)
    if (nodeSiblings === null || nodeIndex === null) return
    nodeSiblings.splice(nodeIndex, 0, dragNode)
  } else if (dropPosition === 'after') {
    const [nodeSiblings, nodeIndex] = findSiblingsAndIndex(node, treeOptions.value)
    if (nodeSiblings === null || nodeIndex === null) return
    nodeSiblings.splice(nodeIndex + 1, 0, dragNode)
  }
  // Vue 3 reactivity trigger
  treeOptions.value = [...treeOptions.value]
}

</script>

<template>
  <div class="flex flex-col flex-1 px-2 py-1 gap-y-1 h-full">
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
<!--    <n-space align="stretch" justify="space-between" :wrap="false">-->

<!--      <n-button type="error" size="tiny" @click="handleDeleteButton">-->
<!--        <template #icon>-->
<!--          <Loader v-if="deleteButtonLoading"/>-->
<!--          <TrashOutline v-else/>-->
<!--        </template>-->
<!--        批量-->
<!--      </n-button>-->
<!--    </n-space>-->


    <!--分割线-->
    <n-divider/>
    <!--列树-->
    <n-scrollbar class="max-h-full overflow-y-auto">
      <!--      @dragend="handleDragEnd"-->
      <n-tree
          :indent="20"
          class="tree"
          cascade
          draggable
          :checkable="deleteCheckAbleEnable"
          @drop="handleDrop"
          :show-irrelevant-nodes="showIrrelevantNodes"
          :pattern="pattern"
          block-line
          :data="treeOptions"
          :default-expanded-keys="expandedKeys"
          :render-switcher-icon="renderSwitcherIconWithExpand"
          checkbox-placement="right"
          :selectable="true"
          :class="{ 'is-dark': isDark }"
          virtual-scroll
          :style="styleObject"
          :render-suffix="renderNodeSuffix"
          :node-props="nodeProps"
          @update:expanded-keys="handleExpandedKeysUpdate"
          @update:checked-keys="updateCheckedKeys"
      />
    </n-scrollbar>
    <n-divider/>
    <!-- bottom function bar -->
    <div class="flex flex-row items-center space-x-2">
      <!--添加会话-->
      <icon-button
          :icon="AddLink"
          :stroke-width="3.5"
          size="20"
          t-tooltip="添加会话"
          @click="handleAddFile"
          :disabled="addFileBtnStatus"
      />
      <!--添加分类-->
      <icon-button
          :icon="AddGroup"
          :stroke-width="3.5"
          size="20"
          :t-tooltip="addFolderBtnStatus ? '最大支持二级分类' : '添加分类'"
          @click="showAddFolderModal = true;"
          :disabled="addFolderBtnStatus"
          :loading="addFolderLoading"
      />
      <!--分割线-->
      <n-divider vertical/>
      <!--详情-->
      <n-tooltip placement="top-start" trigger="hover">
        <template #trigger>
          <n-input :autofocus="false" :placeholder="nodeClickDefaultOption?.label ?? 'ROOT'" size="small" disabled
                   clearable>
            <template #prefix>
              <n-icon :component="Link" size="20"/>
            </template>
          </n-input>
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
      <n-divider vertical/>
      <!--刷新-->
      <icon-button
          :icon="Refresh"
          :stroke-width="3.5"
          size="20"
          t-tooltip="刷新数据"
          @click="handleRefreshData"
          :loading="refreshDataLoading"
      />
    </div>
  </div>
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


