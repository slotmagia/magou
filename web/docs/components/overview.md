# 组件概览

HotGo 2.0 提供了丰富的高质量组件库，覆盖了企业级应用开发的各种场景需求。

## 组件架构

### 设计原则

1. **高内聚低耦合** - 每个组件都是独立的功能单元
2. **配置化驱动** - 通过配置实现功能定制
3. **TypeScript 优先** - 完整的类型支持和智能提示
4. **按需导入** - 支持 Tree Shaking，减少打包体积

### 组件分类

```
组件库
├── 基础组件
│   ├── BasicTable      # 数据表格
│   ├── BasicForm       # 表单组件
│   ├── BasicModal      # 模态框
│   └── BasicUpload     # 文件上传
├── 业务组件
│   ├── CitySelector    # 城市选择器
│   ├── FileChooser     # 文件选择器
│   ├── IconSelector    # 图标选择器
│   └── Editor          # 富文本编辑器
├── 布局组件
│   ├── Application     # 应用容器
│   ├── LoadingContent  # 加载内容
│   └── MessageContent  # 消息内容
└── 工具组件
    ├── CountTo         # 数字动画
    ├── SvgIcon         # SVG 图标
    └── Lockscreen      # 锁屏组件
```

## 核心组件

### 📊 BasicTable - 数据表格

**位置**: `src/components/Table/`

**特性**:
- 📋 丰富的列配置选项
- 🔍 内置搜索和筛选
- 📄 分页和排序支持
- ✏️ 行内编辑功能
- 🎯 操作列配置
- 🎨 自定义渲染

**类型定义**:
```typescript
interface BasicColumn extends TableBaseColumn {
  edit?: boolean;              // 是否可编辑
  editComponent?: ComponentType; // 编辑组件类型
  auth?: string[];             // 权限控制
  ifShow?: boolean | Function; // 显示控制
  draggable?: boolean;         // 是否支持拖拽
}

interface BasicTableProps {
  title?: string;              // 表格标题
  dataSource: Function;        // 数据源函数
  columns: BasicColumn[];      // 列配置
  pagination: object;          // 分页配置
  actionColumn: any[];         // 操作列配置
  loading: boolean;            // 加载状态
}
```

**使用示例**:
```vue
<template>
  <BasicTable
    :columns="columns"
    :dataSource="loadData"
    :pagination="pagination"
    :actionColumn="actionColumn"
    @register="register"
  />
</template>

<script setup lang="ts">
import { BasicTable, useTable } from '@/components/Table';

const columns = [
  {
    title: '用户名',
    key: 'username',
    width: 150,
  },
  {
    title: '邮箱',
    key: 'email',
    width: 200,
  },
  {
    title: '状态',
    key: 'status',
    render: (row) => {
      return h(NTag, { type: row.status === 1 ? 'success' : 'error' }, 
        () => row.status === 1 ? '正常' : '禁用'
      );
    }
  }
];

const actionColumn = [
  {
    label: '编辑',
    key: 'edit',
    auth: ['system:user:edit'],
    onClick: handleEdit,
  },
  {
    label: '删除',
    key: 'delete',
    auth: ['system:user:delete'],
    onClick: handleDelete,
  }
];

const [register, { reload }] = useTable();

async function loadData(params) {
  return await getUserList(params);
}
</script>
```

### 📝 BasicForm - 表单组件

**位置**: `src/components/Form/`

**特性**:
- 🛠️ 丰富的表单控件
- ✅ 内置验证规则
- 🔄 动态表单配置
- 📱 响应式布局
- 🎛️ 高级搜索支持

**类型定义**:
```typescript
interface FormSchema {
  field: string;               // 字段名
  label: string;               // 标签
  component?: ComponentType;   // 组件类型
  componentProps?: object;     // 组件属性
  rules?: object | object[];   // 验证规则
  ifShow?: boolean | Function; // 显示条件
  auth?: string[];             // 权限控制
}

interface FormProps {
  model?: Recordable;          // 表单数据
  schemas?: FormSchema[];      // 表单配置
  labelWidth?: number;         // 标签宽度
  showActionButtonGroup?: boolean; // 显示操作按钮
  submitFunc?: Function;       // 提交函数
  resetFunc?: Function;        // 重置函数
}
```

**使用示例**:
```vue
<template>
  <BasicForm
    @register="register"
    @submit="handleSubmit"
    @reset="handleReset"
  />
</template>

<script setup lang="ts">
import { BasicForm, useForm } from '@/components/Form';

const schemas = [
  {
    field: 'username',
    label: '用户名',
    component: 'NInput',
    componentProps: {
      placeholder: '请输入用户名',
    },
    rules: [
      { required: true, message: '请输入用户名' },
      { min: 3, max: 20, message: '用户名长度在 3 到 20 个字符' }
    ],
  },
  {
    field: 'email',
    label: '邮箱',
    component: 'NInput',
    componentProps: {
      placeholder: '请输入邮箱',
    },
    rules: [
      { required: true, message: '请输入邮箱' },
      { type: 'email', message: '请输入正确的邮箱格式' }
    ],
  },
  {
    field: 'status',
    label: '状态',
    component: 'NSelect',
    componentProps: {
      options: [
        { label: '正常', value: 1 },
        { label: '禁用', value: 0 }
      ]
    },
  }
];

const [register, { setFieldsValue, validate, resetFields }] = useForm({
  schemas,
  showActionButtonGroup: true,
  submitButtonText: '保存',
  resetButtonText: '重置',
});

async function handleSubmit(values) {
  console.log('表单数据:', values);
  // 提交逻辑
}

function handleReset() {
  resetFields();
}
</script>
```

### 📦 BasicModal - 模态框组件

**位置**: `src/components/Modal/`

**特性**:
- 🪟 多种尺寸和位置
- 🖱️ 拖拽功能
- 🖥️ 全屏模式
- ⌨️ 键盘交互
- 🎨 自定义头部和底部

**使用示例**:
```vue
<template>
  <BasicModal
    @register="register"
    title="用户编辑"
    :width="600"
    @ok="handleOk"
  >
    <BasicForm @register="formRegister" />
  </BasicModal>
</template>

<script setup lang="ts">
import { BasicModal, useModal } from '@/components/Modal';
import { BasicForm, useForm } from '@/components/Form';

const [register, { openModal, closeModal, setModalProps }] = useModal();
const [formRegister, { validate }] = useForm();

async function handleOk() {
  try {
    const values = await validate();
    // 保存逻辑
    closeModal();
  } catch (error) {
    console.error('验证失败:', error);
  }
}

// 打开模态框
function openEditModal(record) {
  setModalProps({ title: record.id ? '编辑用户' : '新增用户' });
  openModal();
}
</script>
```

### 📤 BasicUpload - 文件上传

**位置**: `src/components/Upload/`

**特性**:
- 📁 多种上传方式
- 🖼️ 图片预览
- 📋 文件列表管理
- ⏸️ 暂停/恢复上传
- 🔄 断点续传
- 📊 上传进度显示

## 组件开发规范

### 1. 目录结构

```
ComponentName/
├── index.ts              # 组件导出
├── src/                  # 源码目录
│   ├── ComponentName.vue # 主组件
│   ├── components/       # 子组件
│   ├── hooks/           # 组件 hooks
│   ├── props.ts         # 属性定义
│   └── types/           # 类型定义
└── README.md            # 组件文档
```

### 2. 导出规范

```typescript
// index.ts
export { default as ComponentName } from './src/ComponentName.vue';
export { useComponentName } from './src/hooks/useComponentName';
export * from './src/types';
```

### 3. 属性定义

```typescript
// props.ts
import { ExtractPropTypes } from 'vue';

export const componentProps = {
  // 基础属性
  value: {
    type: [String, Number, Array] as PropType<any>,
    default: undefined,
  },
  
  // 配置属性
  disabled: {
    type: Boolean,
    default: false,
  },
  
  // 函数属性
  onChange: {
    type: Function as PropType<(value: any) => void>,
    default: undefined,
  },
} as const;

export type ComponentProps = ExtractPropTypes<typeof componentProps>;
```

### 4. Hooks 开发

```typescript
// hooks/useComponent.ts
import { ref, unref, nextTick } from 'vue';
import type { ComponentProps } from '../types';

export function useComponent(props: ComponentProps) {
  const loading = ref(false);
  
  const setLoading = (value: boolean) => {
    loading.value = value;
  };
  
  const doSomething = async () => {
    try {
      setLoading(true);
      // 业务逻辑
      await nextTick();
    } finally {
      setLoading(false);
    }
  };
  
  return {
    loading: readonly(loading),
    doSomething,
  };
}
```

## 组件通信

### 1. Props/Emits

```vue
<script setup lang="ts">
interface Props {
  modelValue?: string;
  disabled?: boolean;
}

interface Emits {
  (e: 'update:modelValue', value: string): void;
  (e: 'change', value: string): void;
}

const props = withDefaults(defineProps<Props>(), {
  modelValue: '',
  disabled: false,
});

const emits = defineEmits<Emits>();
</script>
```

### 2. Provide/Inject

```typescript
// 提供者
import { provide, InjectionKey } from 'vue';

interface FormContext {
  addField: (field: FormField) => void;
  removeField: (field: FormField) => void;
}

export const FormContextKey: InjectionKey<FormContext> = Symbol('FormContext');

// 在父组件中提供
provide(FormContextKey, {
  addField,
  removeField,
});

// 在子组件中注入
import { inject } from 'vue';

const formContext = inject(FormContextKey);
```

### 3. Event Bus

```typescript
// utils/eventBus.ts
import mitt from 'mitt';

export const eventBus = mitt();

// 发送事件
eventBus.emit('user-updated', userInfo);

// 监听事件
eventBus.on('user-updated', (userInfo) => {
  // 处理用户更新
});
```

## 性能优化

### 1. 懒加载

```typescript
// 异步组件
const AsyncComponent = defineAsyncComponent(() => import('./HeavyComponent.vue'));

// 条件渲染
const LazyComponent = defineAsyncComponent({
  loader: () => import('./LazyComponent.vue'),
  loadingComponent: LoadingComponent,
  errorComponent: ErrorComponent,
  delay: 200,
  timeout: 3000,
});
```

### 2. 虚拟滚动

```vue
<!-- 大列表虚拟滚动 -->
<VirtualList
  :items="largeDataList"
  :item-height="50"
  :visible-count="10"
>
  <template #default="{ item, index }">
    <div>{{ item.name }}</div>
  </template>
</VirtualList>
```

### 3. 缓存优化

```typescript
// 计算属性缓存
const expensiveComputed = computed(() => {
  return heavyCalculation(props.data);
});

// 组件缓存
<KeepAlive :include="['ComponentA', 'ComponentB']">
  <component :is="currentComponent" />
</KeepAlive>
```

## 测试策略

### 1. 单元测试

```typescript
// ComponentName.test.ts
import { mount } from '@vue/test-utils';
import ComponentName from '../src/ComponentName.vue';

describe('ComponentName', () => {
  it('should render correctly', () => {
    const wrapper = mount(ComponentName, {
      props: {
        value: 'test',
      },
    });
    
    expect(wrapper.text()).toContain('test');
  });
  
  it('should emit change event', async () => {
    const wrapper = mount(ComponentName);
    
    await wrapper.find('input').setValue('new value');
    
    expect(wrapper.emitted('change')).toBeTruthy();
  });
});
```

### 2. 集成测试

```typescript
// integration.test.ts
import { mount } from '@vue/test-utils';
import { createRouter, createWebHistory } from 'vue-router';
import { createPinia } from 'pinia';

describe('Component Integration', () => {
  const router = createRouter({
    history: createWebHistory(),
    routes: [],
  });
  
  const pinia = createPinia();
  
  it('should work with router and store', () => {
    const wrapper = mount(ComponentName, {
      global: {
        plugins: [router, pinia],
      },
    });
    
    // 测试逻辑
  });
});
```

---

下一步：[表格组件详解](./table.md)

