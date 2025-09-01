# Hooks 概览

HotGo 2.0 提供了丰富的组合式函数 (Hooks)，基于 Vue 3 Composition API 设计，用于封装可复用的业务逻辑和状态管理。

## Hooks 分类

### 🧰 通用 Hooks (`/src/hooks/common`)

提供常用的工具函数和状态管理：

```
common/
├── useBoolean.ts         # 布尔值状态管理
├── useLoading.ts         # 加载状态管理
├── useLoadingEmpty.ts    # 空数据加载状态
├── useContext.ts         # 上下文管理
├── useSendCode.ts        # 验证码发送
└── useSorter.ts          # 排序管理
```

### 🌐 Web Hooks (`/src/hooks/web`)

Web 相关的业务逻辑封装：

```
web/
├── usePermission.ts      # 权限管理
├── usePage.ts           # 页面管理
└── useECharts.ts        # 图表管理
```

### ⚙️ 设置 Hooks (`/src/hooks/setting`)

系统设置相关的状态管理：

```
setting/
├── useDesignSetting.ts   # 设计设置
└── useProjectSetting.ts  # 项目设置
```

### 🎯 事件 Hooks (`/src/hooks/event`)

事件监听和处理：

```
event/
├── useBreakpoint.ts      # 响应式断点
├── useEventListener.ts   # 事件监听器
└── useWindowSizeFn.ts    # 窗口大小变化
```

## 核心 Hooks 详解

### useBoolean - 布尔值管理

用于管理布尔值状态的常见操作：

```typescript
import { useBoolean } from '@/hooks/common/useBoolean';

export default function useBoolean(defaultValue = false) {
  const bool = ref(defaultValue);

  const setTrue = () => {
    bool.value = true;
  };

  const setFalse = () => {
    bool.value = false;
  };

  const toggle = () => {
    bool.value = !bool.value;
  };

  const setBool = (value: boolean) => {
    bool.value = value;
  };

  return {
    bool: readonly(bool),
    setTrue,
    setFalse,
    toggle,
    setBool,
  };
}
```

**使用示例**：
```vue
<script setup lang="ts">
import { useBoolean } from '@/hooks/common/useBoolean';

// 管理模态框显示状态
const { bool: modalVisible, setTrue: showModal, setFalse: hideModal } = useBoolean();

// 管理加载状态
const { bool: loading, setTrue: startLoading, setFalse: stopLoading } = useBoolean();
</script>

<template>
  <div>
    <NButton @click="showModal">打开模态框</NButton>
    <NModal v-model:show="modalVisible">
      <div>模态框内容</div>
    </NModal>
  </div>
</template>
```

### useLoading - 加载状态管理

封装常见的加载状态管理逻辑：

```typescript
import { useLoading } from '@/hooks/common/useLoading';

export function useLoading(defaultValue = false) {
  const loading = ref(defaultValue);
  const loadingText = ref('');

  const setLoading = (value: boolean, text = '') => {
    loading.value = value;
    loadingText.value = text;
  };

  const startLoading = (text = '加载中...') => {
    setLoading(true, text);
  };

  const stopLoading = () => {
    setLoading(false);
  };

  const withLoading = async <T>(
    fn: () => Promise<T>,
    text = '加载中...'
  ): Promise<T> => {
    try {
      startLoading(text);
      return await fn();
    } finally {
      stopLoading();
    }
  };

  return {
    loading: readonly(loading),
    loadingText: readonly(loadingText),
    setLoading,
    startLoading,
    stopLoading,
    withLoading,
  };
}
```

**使用示例**：
```vue
<script setup lang="ts">
import { useLoading } from '@/hooks/common/useLoading';
import { getUserList } from '@/api/system/user';

const { loading, withLoading } = useLoading();

const handleLoadData = async () => {
  await withLoading(async () => {
    const result = await getUserList();
    // 处理数据
  }, '正在加载用户数据...');
};
</script>

<template>
  <div>
    <NButton :loading="loading" @click="handleLoadData">
      加载数据
    </NButton>
  </div>
</template>
```

### usePermission - 权限管理

提供权限检查和控制功能：

```typescript
import { usePermission } from '@/hooks/web/usePermission';

export function usePermission() {
  const userStore = useUserStore();

  // 检查单个权限
  function hasPermission(permission: string): boolean {
    if (!permission) return true;
    return userStore.permissions.includes(permission);
  }

  // 检查多个权限（满足其中之一）
  function hasSomePermission(permissions: string[]): boolean {
    if (!permissions || !permissions.length) return true;
    return permissions.some(permission => hasPermission(permission));
  }

  // 检查多个权限（必须全部满足）
  function hasEveryPermission(permissions: string[]): boolean {
    if (!permissions || !permissions.length) return true;
    return permissions.every(permission => hasPermission(permission));
  }

  // 检查角色权限
  function hasRole(role: string): boolean {
    if (!role) return true;
    return userStore.roles.includes(role);
  }

  return {
    hasPermission,
    hasSomePermission,
    hasEveryPermission,
    hasRole,
  };
}
```

**使用示例**：
```vue
<script setup lang="ts">
import { usePermission } from '@/hooks/web/usePermission';

const { hasPermission, hasRole } = usePermission();

const canEdit = computed(() => hasPermission('system:user:edit'));
const canDelete = computed(() => hasPermission('system:user:delete'));
const isAdmin = computed(() => hasRole('admin'));
</script>

<template>
  <div>
    <NButton v-if="canEdit" @click="handleEdit">编辑</NButton>
    <NButton v-if="canDelete" @click="handleDelete">删除</NButton>
    <NButton v-if="isAdmin" @click="handleAdmin">管理员操作</NButton>
  </div>
</template>
```

### useECharts - 图表管理

封装 ECharts 图表的初始化和管理：

```typescript
import { useECharts } from '@/hooks/web/useECharts';

export function useECharts(elRef: Ref<HTMLDivElement>) {
  const chartInstance = shallowRef<echarts.ECharts>();
  const { resize, screenEnum } = useBreakpoint();

  // 初始化图表
  const initCharts = () => {
    const el = unref(elRef);
    if (!el || !unref(el)) return;

    chartInstance.value = echarts.init(el);
  };

  // 设置选项
  const setOption = (option: echarts.EChartsOption, clear = true) => {
    if (unref(chartInstance)) {
      clear && unref(chartInstance)?.clear();
      unref(chartInstance)?.setOption(option);
    }
  };

  // 响应式调整大小
  const resizeHandler = useDebounceFn(() => {
    unref(chartInstance)?.resize();
  }, 100);

  // 监听窗口大小变化
  const { removeEvent } = useEventListener({
    el: window,
    name: 'resize',
    listener: resizeHandler,
  });

  tryOnUnmounted(() => {
    if (!chartInstance.value) return;
    removeEvent();
    chartInstance.value.dispose();
    chartInstance.value = null;
  });

  return {
    initCharts,
    setOption,
    resize: resizeHandler,
    echarts,
    chartInstance,
  };
}
```

**使用示例**：
```vue
<script setup lang="ts">
import { useECharts } from '@/hooks/web/useECharts';

const chartRef = ref<HTMLDivElement>();
const { initCharts, setOption } = useECharts(chartRef);

const option = {
  title: { text: '销售统计' },
  xAxis: { type: 'category', data: ['一月', '二月', '三月'] },
  yAxis: { type: 'value' },
  series: [{
    type: 'bar',
    data: [120, 200, 150]
  }]
};

onMounted(() => {
  initCharts();
  setOption(option);
});
</script>

<template>
  <div ref="chartRef" style="width: 100%; height: 400px;"></div>
</template>
```

## 自定义 Hooks 开发

### 开发规范

1. **命名规范**: 使用 `use` 前缀，采用 camelCase 命名
2. **文件结构**: 每个 Hook 独立文件，提供类型定义
3. **返回值**: 统一返回对象，提供响应式数据和方法
4. **副作用清理**: 使用 `tryOnUnmounted` 清理副作用

### Hook 模板

```typescript
// hooks/useExample.ts
import { ref, reactive, computed, readonly } from 'vue';
import { tryOnUnmounted } from '@vueuse/core';

export interface UseExampleOptions {
  immediate?: boolean;
  // 其他选项
}

export function useExample(options: UseExampleOptions = {}) {
  const { immediate = true } = options;
  
  // 响应式状态
  const state = reactive({
    loading: false,
    data: null as any,
    error: null as Error | null,
  });
  
  const count = ref(0);
  
  // 计算属性
  const doubleCount = computed(() => count.value * 2);
  
  // 方法
  const increment = () => {
    count.value++;
  };
  
  const decrement = () => {
    count.value--;
  };
  
  const asyncOperation = async () => {
    try {
      state.loading = true;
      state.error = null;
      
      // 异步操作
      const result = await fetchData();
      state.data = result;
      
    } catch (error) {
      state.error = error as Error;
    } finally {
      state.loading = false;
    }
  };
  
  // 初始化
  if (immediate) {
    asyncOperation();
  }
  
  // 清理副作用
  tryOnUnmounted(() => {
    // 清理逻辑
  });
  
  return {
    // 只读状态
    ...toRefs(readonly(state)),
    count: readonly(count),
    doubleCount,
    
    // 方法
    increment,
    decrement,
    asyncOperation,
  };
}

// 类型导出
export type UseExampleReturn = ReturnType<typeof useExample>;
```

### 复杂 Hook 示例

```typescript
// hooks/useTableCrud.ts
import { reactive, ref } from 'vue';
import { useModal } from '@/components/Modal';
import { useTable } from '@/components/Table';

export interface UseTableCrudOptions<T = any> {
  api: {
    list: (params: any) => Promise<any>;
    create: (data: T) => Promise<any>;
    update: (id: string | number, data: T) => Promise<any>;
    delete: (id: string | number) => Promise<any>;
  };
  columns: any[];
  searchSchemas?: any[];
}

export function useTableCrud<T = any>(options: UseTableCrudOptions<T>) {
  const { api, columns, searchSchemas = [] } = options;
  
  // 状态管理
  const state = reactive({
    selectedRowKeys: [] as string[],
    selectedRows: [] as T[],
  });
  
  const currentRecord = ref<T | null>(null);
  
  // 表格 Hook
  const [tableRegister, tableActions] = useTable({
    columns,
    dataSource: api.list,
    rowSelection: {
      type: 'checkbox',
      onChange: (keys: string[], rows: T[]) => {
        state.selectedRowKeys = keys;
        state.selectedRows = rows;
      },
    },
  });
  
  // 模态框 Hook
  const [modalRegister, modalActions] = useModal();
  
  // CRUD 操作
  const handleCreate = () => {
    currentRecord.value = null;
    modalActions.openModal();
  };
  
  const handleEdit = (record: T) => {
    currentRecord.value = record;
    modalActions.openModal();
  };
  
  const handleDelete = async (record: T) => {
    try {
      await api.delete((record as any).id);
      window.$message.success('删除成功');
      tableActions.reload();
    } catch (error) {
      window.$message.error('删除失败');
    }
  };
  
  const handleBatchDelete = async () => {
    if (!state.selectedRowKeys.length) {
      window.$message.warning('请选择要删除的记录');
      return;
    }
    
    try {
      await Promise.all(
        state.selectedRowKeys.map(id => api.delete(id))
      );
      window.$message.success('批量删除成功');
      tableActions.reload();
      state.selectedRowKeys = [];
      state.selectedRows = [];
    } catch (error) {
      window.$message.error('批量删除失败');
    }
  };
  
  const handleSave = async (values: T) => {
    try {
      if (currentRecord.value) {
        await api.update((currentRecord.value as any).id, values);
        window.$message.success('更新成功');
      } else {
        await api.create(values);
        window.$message.success('创建成功');
      }
      
      modalActions.closeModal();
      tableActions.reload();
    } catch (error) {
      window.$message.error('保存失败');
    }
  };
  
  return {
    // 状态
    state: readonly(state),
    currentRecord: readonly(currentRecord),
    
    // 注册函数
    tableRegister,
    modalRegister,
    
    // 操作方法
    handleCreate,
    handleEdit,
    handleDelete,
    handleBatchDelete,
    handleSave,
    
    // 表格操作
    ...tableActions,
    
    // 模态框操作
    ...modalActions,
  };
}
```

## 最佳实践

### 1. 合理拆分逻辑

```typescript
// 单一职责：只处理表单逻辑
export function useFormLogic() {
  // 表单相关逻辑
}

// 单一职责：只处理数据请求
export function useDataFetching() {
  // 数据请求逻辑
}

// 组合使用
export function useUserManagement() {
  const formLogic = useFormLogic();
  const dataFetching = useDataFetching();
  
  return {
    ...formLogic,
    ...dataFetching,
  };
}
```

### 2. 提供类型支持

```typescript
export interface UseApiOptions<T> {
  immediate?: boolean;
  onSuccess?: (data: T) => void;
  onError?: (error: Error) => void;
}

export function useApi<T = any>(
  api: () => Promise<T>,
  options: UseApiOptions<T> = {}
) {
  // 实现逻辑
}
```

### 3. 错误处理

```typescript
export function useAsyncOperation() {
  const error = ref<Error | null>(null);
  
  const execute = async (fn: () => Promise<any>) => {
    try {
      error.value = null;
      return await fn();
    } catch (err) {
      error.value = err as Error;
      console.error('操作失败:', err);
      throw err;
    }
  };
  
  return {
    error: readonly(error),
    execute,
  };
}
```

### 4. 生命周期管理

```typescript
export function useInterval(callback: () => void, delay: number) {
  const intervalId = ref<number | null>(null);
  
  const start = () => {
    if (intervalId.value) return;
    intervalId.value = setInterval(callback, delay);
  };
  
  const stop = () => {
    if (intervalId.value) {
      clearInterval(intervalId.value);
      intervalId.value = null;
    }
  };
  
  tryOnUnmounted(stop);
  
  return { start, stop };
}
```

---

下一步：[API 接口文档](../api/design.md)

