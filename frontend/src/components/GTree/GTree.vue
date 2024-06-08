<script setup lang="ts">
import {computed, h, onMounted, PropType, ref, watch, withModifiers} from 'vue';
import {NButton, NDivider, NIcon, TreeOption, useOsTheme} from 'naive-ui';
import {AddOutline, FolderOpenOutline, FolderOutline, PencilOutline, TrashOutline} from '@vicons/ionicons5';
import SearchOutline from '@/src/components/Icons/SearchOutline.vue';
import Loader from '@/src/components/Icons/Loader.vue';
import {findTreeOptionLevel} from '@/src/components/GTree/createData';


const props = defineProps({
  data: {
    type: Array as PropType<TreeOption[]>,
    required: true
  }
});

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

onMounted(() => {
  loadFromLocalStorage();
});

watch(expandedKeys, async (newKeys) => {
  saveToLocalStorage();
});

const handleDragEnd = (event) => {
  $message.info('节点拖拽');
  console.log('Drag End:', event);
};

// 批量删除按钮
const deleteButtonLoading = ref(false);

const deleteCheckAbleEnable = ref(false);

const handleDeleteButton = () => {
  deleteCheckAbleEnable.value = !deleteCheckAbleEnable.value;
  deleteButtonLoading.value = !deleteButtonLoading.value;
};


// 处理节点事件
const emit = defineEmits(['selected']);
const nodeClickDefaultNodeOption = ref<TreeOption | null>(null)
const nodeProps = ({option}: { option: TreeOption }) => {
  return {
    onClick() {
      // 保存当前点击的节点
      nodeClickDefaultNodeOption.value = option;

      const isFolder = typeof option.children !== 'undefined';
      emit('selected', {id: option.key, label: option.label, isFolder});
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
const nodeHoverDefaultKey = ref<string | number>(null);
const renderNodeSuffix = ({option}: { option: TreeOption }) => {
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
const renderSwitcherIconWithExpand = ({expanded}: { expanded: boolean }) =>
    h(NIcon, null, {
      default: () => h(expanded ? FolderOpenOutline : FolderOutline)
    });


// 新增会话处理
const handleAddSessionStatus = computed(() => {
  return props.data.length === 0;
})

const handleAddSession = () => {
  if (!nodeClickDefaultNodeOption.value) {
    $message.warning('请先选择一个分类');
    return;
  }
  $message.info('新增会话');
};
</script>


<template>
  <div class="flex flex-col flex-grow p-2 h-full gap-3 g-tree">
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
      <n-button type="info" size="tiny">
        <template #icon>
          <AddOutline/>
        </template>
        新增分类
      </n-button>
      <n-button type="primary" size="tiny" @click="handleAddSession" :disabled="handleAddSessionStatus">
        <template #icon>
          <AddOutline/>
        </template>
        新增会话
      </n-button>
      <n-button type="error" size="tiny" @click="handleDeleteButton">
        <template #icon>
          <Loader v-if="deleteButtonLoading"/>
          <TrashOutline v-else/>
        </template>
        批量删除
      </n-button>
    </n-space>
    <!--分割线-->
    <n-divider/>
    <!--详情-->
    <n-space vertical>
      <n-tag disabled>
        当前Label: {{ nodeClickDefaultNodeOption?.label ?? '无' }}
      </n-tag>
      <n-tag disabled>
        当前Key: {{ nodeClickDefaultNodeOption?.key ?? '无' }}
      </n-tag>
      <n-tag disabled>
        当前Level: {{
          nodeClickDefaultNodeOption?.key ? findTreeOptionLevel(props.data, nodeClickDefaultNodeOption) : '无'
        }}
      </n-tag>
      <n-tag disabled>
        是否是文件夹: {{ typeof nodeClickDefaultNodeOption?.children !== 'undefined' ? '是' : '否' }}
      </n-tag>
    </n-space>


    <!--分割线-->
    <n-divider/>
    <!--列树-->
    <n-scrollbar class="max-h-full">
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
          :data="props.data"
          :default-expanded-keys="expandedKeys"
          :on-update:expanded-keys="handleExpandedKeysUpdate"
          :render-switcher-icon="renderSwitcherIconWithExpand"
          checkbox-placement="right"
          :selectable="true"
          :node-props="nodeProps"
          :class="{ 'is-dark': isDark }"
          virtual-scroll
          :style="styleObject"
          :render-suffix="renderNodeSuffix"
          :node-prop="nodeProps"
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


