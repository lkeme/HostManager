// import {h} from 'vue'
import {repeat} from 'seemly'
// import {NButton, TreeOption} from 'naive-ui'
import {TreeOption,TreeDropInfo} from 'naive-ui'
import {v4 as uuidv4} from 'uuid'

export interface NewTreeOption extends TreeOption {
    id: number,
    level: number,
    parentID: number,
    children?: NewTreeOption[]
}


export interface NewTreeDropInfo extends TreeDropInfo {
    event: DragEvent;
    node: NewTreeOption;
    dragNode: NewTreeOption;
    dropPosition: 'before' | 'inside' | 'after';
}

const createKey = (): string => {
    // 550e8400-e29b-41d4-a716-446655440000
    // 取e29b-41d4-a716
    return uuidv4().slice(9, 23)
}

const createLabel = (level: number): string => {
    if (level === 4) return 'Level 4'
    if (level === 3) return 'Level 3'
    if (level === 2) return 'Level 2'
    if (level === 1) return 'Level 1'
    return ''
}

// 创建虚拟初始数据
export const createData = (level = 3, baseKey = ''): TreeOption[] => {
    if (!level) return [];
    return repeat(6 - level, undefined).map((_, index) => {
        // const key = '' + baseKey + level + index
        const key = createKey()
        const label = createLabel(level)
        return {
            label,
            key,
            children: level > 1 ? createData(level - 1, key) : undefined,
        }
    })
}

// idで指定したフォルダ以下に新規フォルダを追加
export const addFolderDataById = (tree: TreeOption[], id: string, name: string): void => {
    const createFolderData = (name: string) => {
        return {
            label: name,
            key: uuidv4(),
            children: [],
        }
    }

    if (!id) {
        const newData = createFolderData(name);
        tree.push(newData);
        return;
    }

    tree.forEach((node) => {
        if (node.key === id) {
            const newData = createFolderData(name);

            if (!node.children) node.children = [];
            node.children.push(newData);
            return;
        }
        if (node.children) {
            addFolderDataById(node.children, id, name);
        }
    });
};


// idで指定したフォルダ以下に新規文書を追加
export const addDocumentDataById = (tree: TreeOption[], id: string, name: string): void => {
    const createDocumentData = (name: string) => {
        return {
            label: name,
            key: uuidv4(),
        }
    }

    if (!id) {
        const newData = createDocumentData(name);
        tree.push(newData);
        return;
    }

    tree.forEach((node) => {
        if (node.key === id) {
            const newData = createDocumentData(name);

            if (!node.children) node.children = [];
            node.children.push(newData);
            return;
        }
        if (node.children) {
            addDocumentDataById(node.children, id, name);
        }
    });
};


export const removeDataById = (tree: TreeOption[], id: string): void => {
    for (let i = 0; i < tree.length; i++) {
        const node = tree[i];
        if (node.key === id) {
            // IDが一致する場合、このノードを削除
            tree.splice(i, 1);
            return;
        }
        if (node.children) {
            // 子ノードがある場合、再帰的に探索
            removeDataById(node.children, id);
        }
    }
};

export const renameDataById = (tree: TreeOption[], id: string, newName: string): void => {
    for (let i = 0; i < tree.length; i++) {
        const node = tree[i];
        if (node.key === id) {
            // IDが一致する場合、このノードの名前を変更
            node.label = newName;
            return;
        }
        if (node.children) {
            // 子ノードがある場合、再帰的に探索
            renameDataById(node.children, id, newName);
        }
    }
};

// 查找树形结构中节点的层级
export const findTreeOptionLevel = (tree: TreeOption[], target: TreeOption, level = 0): number => {
    for (const option of tree) {
        if (option === target) {
            return level;
        }
        if (option.children) {
            const childLevel = findTreeOptionLevel(option.children, target, level + 1);
            if (childLevel !== -1) {
                return childLevel;
            }
        }
    }
    return -1;
}