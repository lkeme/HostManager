<script setup lang="ts">
import {computed, h, onMounted, ref, watch, withModifiers} from 'vue';
import {NButton, NDivider, NIcon, TreeOption, useOsTheme} from 'naive-ui';
import {findTreeOptionLevel, NewTreeDropInfo, NewTreeOption} from '@/src/layouts/DashboardPanel/createData';
import {GetStorage} from "@/wailsjs/go/controller/FileSystem";
import {Folder, GetStorageResponse} from "@/src/types/response/filesystem";
import IconButton from "@/src/components/common/IconButton.vue";
import AddGroup from "@/src/components/Icons/AddGroup.vue";
import AddLink from "@/src/components/Icons/AddLink.vue";
import SearchOutline from '@/src/components/Icons/SearchOutline.vue';
import Refresh from '@/src/components/Icons/Refresh.vue';
import {
  FileTrayFullOutline,
  Folder as FolderClose,
  FolderOpen,
  FolderOpenOutline,
  FolderOutline as FolderCloseOutline,
  Link,
  SettingsOutline,
  TrashOutline,
} from '@vicons/ionicons5';
import FolderDialog from "@/src/components/dialogs/FolderDialog.vue";

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


const handleExpandedKeysUpdate = (
    _keys: Array<string | number>,
    _option: Array<TreeOption | null>,
    meta: {
      node: TreeOption | null
      action: 'expand' | 'collapse' | 'filter'
    }
) => {
  if (!meta.node) return;
  //
  expandedKeys.value = _keys;

  console.log('Updated Expanded Keys:', _keys);
  console.log('Options:', _option);
  console.log('Node:', meta.node);
  console.log('Action:', meta.action);

  //
  if (meta.action === 'expand') {
    console.log(`${meta.node.label} has been expanded!`);
  } else if (meta.action === 'collapse') {
    console.log(`${meta.node.label} has been collapsed!`);
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

const getNodeMenu = (option: NewTreeOption, isEdit: boolean = false, isDelete: boolean = false) => {
  nodeHoverDefaultOption.value = option

  const style = 'display: flex; align-items: center; justify-content: center;';
  const buttons = [];

  if (isEdit) {
    buttons.push(createButton(SettingsOutline, '#2ed6e9', '编辑', () => {
      handleRenameFolder()
    }));
  }

  if (isDelete) {
    buttons.push(createButton(TrashOutline, '#fb0e0e', '删除', () => {
      handleDeleteFolder()
    }));
  }

  return h('div', {style}, buttons);
};

const renderNodeSuffix = ({option}: { option: NewTreeOption }) => {
  // 如果不是当前hover的节点，不渲染
  if (option.key !== nodeHoverDefaultKey.value) return;
  // 如果是文件夹，且有子节点，只渲染编辑图标按钮
  if (option.children && option.children.length > 0) {
    return getNodeMenu(option, true, false);
  } else {
    return getNodeMenu(option, true, true);
  }
}

// 渲染展开开关的图标
const renderSwitcherIconWithExpand = ({option, expanded}: { option: NewTreeOption, expanded: boolean }) => {
  return h(NIcon, null, {
    default: () => {
      if (option.type !== 'folder') {
        return h(FileTrayFullOutline)
      }
      if (typeof option.children !== 'undefined' && option.children.length > 0) {
        return h(expanded ? FolderOpen : FolderClose)
      }
      return h(expanded ? FolderOpenOutline : FolderCloseOutline)
    }
  });
}
// -----------------nodeHandle end-----------------

// -----------refreshData start----------------
const refreshDataLoading = ref(false);

// 将 Folder 转换为 NewTreeOption 的函数
const convertFolderToNewTreeOption = (folder: Folder): NewTreeOption => {
  let newTreeOption: NewTreeOption = {
    id: folder.ID,
    level: folder.level,
    parentID: folder.parentID,
    type: folder.type,
    // isLeaf: !children || children.length === 0, // 添加对children为null或undefined的检查
    key: folder.key,
    label: folder.name,
    children: []
  };
  if (folder.subfolders?.length > 0) {
    newTreeOption.children = folder.subfolders.map(subfolder => convertFolderToNewTreeOption(subfolder));
  }

  return newTreeOption;
}

// 将文件夹层级结构转换为树形结构
const convertStorageToNewTreeOption = (storage: Array<Folder>): NewTreeOption[] => {
  return storage.map(folder => convertFolderToNewTreeOption(folder));
}

const handleRefreshData = async () => {
  refreshDataLoading.value = true
  try {
    const response1 = await GetStorage();
    console.log(response1)

    const response: GetStorageResponse = await GetStorage();
    const {data} = response;
    if (data.status) {
      treeOptions.value = convertStorageToNewTreeOption(data.storage);
      $message.success('加载数据成功');
    } else {
      $message.error(`加载数据失败, ${data.message}` || '加载数据失败，请刷新重试');
    }
  } finally {
    refreshDataLoading.value = false;
  }
};
// -----------refreshData end----------------


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


// -----------分类操作----------------
const nodeHoverDefaultOption = ref<NewTreeOption | null>(null);
const showFolderModal = ref(false);
const folderModalTitle = ref(<'新增分类' | '删除分类' | '重命名分类'>'新增分类');

const addFolderBtnLoading = ref(false);
const addFolderBtnStatus = computed(() => {
  return nodeClickDefaultOption.value?.level > 1;
})

const handleAddFolder = () => {
  nodeHoverDefaultOption.value = nodeClickDefaultOption.value;
  folderModalTitle.value = '新增分类';
  showFolderModal.value = true;
}

const handleRenameFolder = () => {
  folderModalTitle.value = '重命名分类';
  showFolderModal.value = true;
}

const handleDeleteFolder = () => {
  folderModalTitle.value = '删除分类';
  showFolderModal.value = true;
}

const handleFolder = async (status: boolean) => {
  if (status) {
    await handleRefreshData();
  }
  showFolderModal.value = false;
};

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
    <!--      @dragend="handleDragEnd"-->
    <!--    expand-on-click-->
    <n-tree
        :indent="20"
        class-backup="max-h-full overflow-y-auto"
        class="tree max-h-full"
        block-line
        cascade
        draggable
        :checkable="deleteCheckAbleEnable"
        @drop="handleDrop"
        :show-irrelevant-nodes="showIrrelevantNodes"
        :pattern="pattern"
        :data="treeOptions"
        :default-expanded-keys="expandedKeys"
        :render-switcher-icon="renderSwitcherIconWithExpand"
        checkbox-placement="right"
        :selectable="true"
        :class="{ 'is-dark': isDark }"
        virtual-scroll
        :style="styleObject"
        @update:expanded-keys="handleExpandedKeysUpdate"
        @update:checked-keys="updateCheckedKeys"
        :render-suffix="renderNodeSuffix"
        :node-props="nodeProps"
    >
      <template #empty>
        <n-empty size="large" class="h-full justify-center" description="还未添加分类数据"/>
      </template>
    </n-tree>
    <!--分割线-->
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
          @click="handleAddFolder"
          :disabled="addFolderBtnStatus"
          :loading="addFolderBtnLoading"
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
            IsFolder: {{ nodeClickDefaultOption?.type === 'folder' ? '☑' : '☒' }}
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
  <FolderDialog
      v-model:show="showFolderModal"
      :title="folderModalTitle"
      :option="nodeHoverDefaultOption"
      @submit="handleFolder"
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


