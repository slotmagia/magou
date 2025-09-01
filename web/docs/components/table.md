# BasicTable 表格组件

BasicTable 是一个功能丰富的数据表格组件，基于 Naive UI 的 DataTable 组件封装，提供了开箱即用的企业级表格功能。

## 基础用法

### 简单表格

```vue
<template>
  <BasicTable
    :columns="columns"
    :dataSource="loadUserList"
    :pagination="paginationReactive"
  />
</template>

<script setup lang="ts">
import { BasicTable } from '@/components/Table';
import { getUserList } from '@/api/system/user';

const columns = [
  {
    title: 'ID',
    key: 'id',
    width: 80,
  },
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
    title: '创建时间',
    key: 'created_at',
    width: 180,
  },
];

const paginationReactive = reactive({
  page: 1,
  pageSize: 10,
  showSizePicker: true,
  pageSizes: [10, 20, 50, 100],
  onChange: (page: number) => {
    paginationReactive.page = page;
  },
  onUpdatePageSize: (pageSize: number) => {
    paginationReactive.pageSize = pageSize;
    paginationReactive.page = 1;
  },
});

async function loadUserList(params) {
  const { page, pageSize } = params;
  return await getUserList({ page, pageSize });
}
</script>
```

## 高级用法

### 带操作列的表格

```vue
<template>
  <BasicTable
    :columns="columns"
    :dataSource="loadUserList"
    :pagination="paginationReactive"
    :actionColumn="actionColumn"
    :rowKey="(row) => row.id"
    @register="register"
  />
  
  <!-- 编辑模态框 -->
  <EditModal @register="editRegister" @success="reload" />
</template>

<script setup lang="ts">
import { h } from 'vue';
import { NButton, NButtonGroup, NPopconfirm, NTag } from 'naive-ui';
import { BasicTable, useTable, type BasicColumn } from '@/components/Table';
import { useModal } from '@/components/Modal';
import { usePermission } from '@/hooks/web/usePermission';
import { getUserList, deleteUser } from '@/api/system/user';
import EditModal from './components/EditModal.vue';

const { hasPermission } = usePermission();
const [editRegister, { openModal: openEditModal }] = useModal();
const [register, { reload }] = useTable();

// 列配置
const columns: BasicColumn[] = [
  {
    title: 'ID',
    key: 'id',
    width: 80,
  },
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
    width: 100,
    render: (row) => {
      return h(
        NTag,
        { type: row.status === 1 ? 'success' : 'error' },
        { default: () => (row.status === 1 ? '正常' : '禁用') }
      );
    },
  },
  {
    title: '创建时间',
    key: 'created_at',
    width: 180,
  },
];

// 操作列配置
const actionColumn = {
  width: 200,
  title: '操作',
  key: 'action',
  fixed: 'right',
  render(record) {
    return h(NButtonGroup, { size: 'small' }, {
      default: () => [
        h(
          NButton,
          {
            size: 'small',
            type: 'primary',
            secondary: true,
            onClick: () => handleEdit(record),
            style: { display: hasPermission(['system:user:edit']) ? 'inline-block' : 'none' },
          },
          { default: () => '编辑' }
        ),
        h(
          NPopconfirm,
          {
            onPositiveClick: () => handleDelete(record),
            style: { display: hasPermission(['system:user:delete']) ? 'inline-block' : 'none' },
          },
          {
            default: () => '确定删除吗？',
            trigger: () => h(
              NButton,
              { size: 'small', type: 'error', secondary: true },
              { default: () => '删除' }
            ),
          }
        ),
      ],
    });
  },
};

// 数据加载函数
async function loadUserList(params) {
  const { page, pageSize } = params;
  return await getUserList({ page, pageSize });
}

// 编辑用户
function handleEdit(record) {
  openEditModal(true, record);
}

// 删除用户
async function handleDelete(record) {
  try {
    await deleteUser(record.id);
    window.$message.success('删除成功');
    reload();
  } catch (error) {
    window.$message.error('删除失败');
  }
}
</script>
```

### 可编辑表格

```vue
<template>
  <BasicTable
    :columns="editableColumns"
    :dataSource="loadData"
    :can-resize="true"
    edit-mode="row"
    @register="register"
  />
</template>

<script setup lang="ts">
import { BasicTable, useTable, type BasicColumn } from '@/components/Table';

const [register, { reload, getDataSource }] = useTable();

const editableColumns: BasicColumn[] = [
  {
    title: '商品名称',
    key: 'name',
    width: 200,
    edit: true,
    editComponent: 'NInput',
    editRule: true,
  },
  {
    title: '价格',
    key: 'price',
    width: 120,
    edit: true,
    editComponent: 'NInputNumber',
    editComponentProps: {
      precision: 2,
      min: 0,
    },
  },
  {
    title: '分类',
    key: 'category',
    width: 150,
    edit: true,
    editComponent: 'NSelect',
    editComponentProps: {
      options: [
        { label: '电子产品', value: 'electronics' },
        { label: '服装', value: 'clothing' },
        { label: '食品', value: 'food' },
      ],
    },
  },
  {
    title: '状态',
    key: 'status',
    width: 100,
    edit: true,
    editComponent: 'NSwitch',
    editValueMap: (value) => (value ? '启用' : '禁用'),
  },
];

async function loadData() {
  // 返回数据
  return {
    items: [
      { id: 1, name: 'iPhone 13', price: 5999, category: 'electronics', status: true },
      { id: 2, name: '连衣裙', price: 299, category: 'clothing', status: true },
    ],
    total: 2,
  };
}
</script>
```

### 带搜索的表格

```vue
<template>
  <div>
    <!-- 搜索表单 -->
    <BasicForm
      @register="searchRegister"
      @submit="handleSearch"
      @reset="handleReset"
    />
    
    <!-- 数据表格 -->
    <BasicTable
      :columns="columns"
      :dataSource="loadUserList"
      @register="tableRegister"
    />
  </div>
</template>

<script setup lang="ts">
import { BasicTable, useTable } from '@/components/Table';
import { BasicForm, useForm } from '@/components/Form';
import { getUserList } from '@/api/system/user';

const [tableRegister, { reload, setTableData }] = useTable();

// 搜索表单配置
const searchSchemas = [
  {
    field: 'username',
    label: '用户名',
    component: 'NInput',
    componentProps: {
      placeholder: '请输入用户名',
      clearable: true,
    },
  },
  {
    field: 'email',
    label: '邮箱',
    component: 'NInput',
    componentProps: {
      placeholder: '请输入邮箱',
      clearable: true,
    },
  },
  {
    field: 'status',
    label: '状态',
    component: 'NSelect',
    componentProps: {
      placeholder: '请选择状态',
      options: [
        { label: '全部', value: '' },
        { label: '正常', value: 1 },
        { label: '禁用', value: 0 },
      ],
    },
  },
  {
    field: 'dateRange',
    label: '创建时间',
    component: 'NDatePicker',
    componentProps: {
      type: 'daterange',
      clearable: true,
    },
  },
];

const [searchRegister, { getFieldsValue, resetFields }] = useForm({
  schemas: searchSchemas,
  showAdvancedButton: true,
  showActionButtonGroup: true,
  submitButtonText: '搜索',
  resetButtonText: '重置',
  baseColProps: { span: 6 },
});

// 列配置
const columns = [
  // ... 列配置
];

// 数据加载
async function loadUserList(params) {
  const searchParams = getFieldsValue();
  const mergedParams = { ...params, ...searchParams };
  
  // 处理日期范围
  if (mergedParams.dateRange) {
    mergedParams.startDate = mergedParams.dateRange[0];
    mergedParams.endDate = mergedParams.dateRange[1];
    delete mergedParams.dateRange;
  }
  
  return await getUserList(mergedParams);
}

// 搜索
function handleSearch() {
  reload();
}

// 重置
function handleReset() {
  resetFields();
  reload();
}
</script>
```

## API 接口

### BasicTable Props

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| columns | `BasicColumn[]` | `[]` | 表格列配置 |
| dataSource | `Function` | - | 数据源函数 |
| pagination | `object` | `{}` | 分页配置 |
| actionColumn | `object` | - | 操作列配置 |
| rowKey | `string \| Function` | `'id'` | 行数据的 Key |
| loading | `boolean` | `false` | 加载状态 |
| canResize | `boolean` | `false` | 是否可调整列宽 |
| showPagination | `boolean` | `true` | 是否显示分页 |
| size | `'small' \| 'medium' \| 'large'` | `'medium'` | 表格尺寸 |
| bordered | `boolean` | `true` | 是否显示边框 |
| striped | `boolean` | `false` | 是否显示斑马纹 |

### BasicColumn 配置

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| title | `string` | - | 列标题 |
| key | `string` | - | 列数据字段名 |
| width | `number \| string` | - | 列宽度 |
| minWidth | `number` | - | 最小列宽 |
| maxWidth | `number` | - | 最大列宽 |
| fixed | `'left' \| 'right'` | - | 固定列 |
| align | `'left' \| 'center' \| 'right'` | `'left'` | 对齐方式 |
| sortable | `boolean` | `false` | 是否可排序 |
| filterable | `boolean` | `false` | 是否可筛选 |
| resizable | `boolean` | `false` | 是否可调整大小 |
| ellipsis | `boolean` | `false` | 是否省略 |
| render | `Function` | - | 自定义渲染函数 |
| edit | `boolean` | `false` | 是否可编辑 |
| editComponent | `ComponentType` | - | 编辑组件类型 |
| editComponentProps | `object` | - | 编辑组件属性 |
| auth | `string[]` | - | 权限控制 |
| ifShow | `boolean \| Function` | `true` | 显示条件 |

### useTable Hook

```typescript
const [register, methods] = useTable(options);
```

#### 返回方法

| 方法名 | 参数 | 返回值 | 说明 |
|--------|------|--------|------|
| reload | `(opt?: object)` | `Promise<void>` | 重新加载数据 |
| setTableData | `(data: any[])` | `void` | 设置表格数据 |
| getTableData | - | `any[]` | 获取表格数据 |
| getColumns | `(opt?: object)` | `BasicColumn[]` | 获取列配置 |
| setColumns | `(columns: BasicColumn[])` | `void` | 设置列配置 |
| setLoading | `(loading: boolean)` | `void` | 设置加载状态 |
| getSize | - | `number` | 获取数据条数 |
| updateTableData | `(index: number, key: string, value: any)` | `void` | 更新表格数据 |
| updateTableDataRecord | `(rowKey: any, record: object)` | `void` | 更新表格记录 |
| deleteTableDataRecord | `(rowKey: any)` | `void` | 删除表格记录 |
| insertTableDataRecord | `(record: object, index?: number)` | `void` | 插入表格记录 |

## 事件

| 事件名 | 参数 | 说明 |
|--------|------|------|
| register | `(instance)` | 注册表格实例 |
| selection-change | `(keys: any[], rows: any[])` | 选择项变化 |
| row-click | `(row: any, index: number)` | 行点击 |
| row-dblclick | `(row: any, index: number)` | 行双击 |
| edit-end | `(value: any, row: any, column: BasicColumn)` | 编辑结束 |
| edit-cancel | `(row: any, column: BasicColumn)` | 编辑取消 |

## 自定义渲染

### 使用 render 函数

```typescript
const columns = [
  {
    title: '状态',
    key: 'status',
    render: (row) => {
      return h(
        NTag,
        { 
          type: row.status === 1 ? 'success' : 'error',
          size: 'small',
        },
        { default: () => row.status === 1 ? '正常' : '禁用' }
      );
    },
  },
  {
    title: '头像',
    key: 'avatar',
    render: (row) => {
      return h(
        NAvatar,
        {
          src: row.avatar,
          size: 'small',
          fallbackSrc: '/default-avatar.png',
        }
      );
    },
  },
];
```

### 使用插槽

```vue
<template>
  <BasicTable :columns="columns">
    <template #status="{ record }">
      <NTag :type="record.status === 1 ? 'success' : 'error'">
        {{ record.status === 1 ? '正常' : '禁用' }}
      </NTag>
    </template>
    
    <template #action="{ record }">
      <NButtonGroup size="small">
        <NButton @click="handleEdit(record)">编辑</NButton>
        <NButton @click="handleDelete(record)">删除</NButton>
      </NButtonGroup>
    </template>
  </BasicTable>
</template>
```

## 性能优化

### 虚拟滚动

对于大量数据，可以启用虚拟滚动：

```vue
<BasicTable
  :columns="columns"
  :dataSource="loadData"
  virtual-scroll
  :scroll-x="1200"
  :max-height="400"
/>
```

### 列宽优化

合理设置列宽，避免频繁重新计算：

```typescript
const columns = [
  {
    title: 'ID',
    key: 'id',
    width: 80,        // 固定宽度
    minWidth: 80,     // 最小宽度
  },
  {
    title: '描述',
    key: 'description',
    ellipsis: true,   // 超长省略
    minWidth: 200,
  },
];
```

### 数据缓存

```typescript
// 使用缓存避免重复请求
const dataCache = new Map();

async function loadUserList(params) {
  const cacheKey = JSON.stringify(params);
  
  if (dataCache.has(cacheKey)) {
    return dataCache.get(cacheKey);
  }
  
  const result = await getUserList(params);
  dataCache.set(cacheKey, result);
  
  return result;
}
```

## 最佳实践

### 1. 权限控制

```typescript
import { usePermission } from '@/hooks/web/usePermission';

const { hasPermission } = usePermission();

const columns = [
  {
    title: '敏感数据',
    key: 'sensitive',
    ifShow: () => hasPermission(['admin:view']),
  },
];
```

### 2. 响应式设计

```typescript
import { useBreakpoint } from '@/hooks/event/useBreakpoint';

const { screenEnum } = useBreakpoint();

const columns = computed(() => [
  {
    title: 'ID',
    key: 'id',
    width: 80,
  },
  {
    title: '详细信息',
    key: 'detail',
    ifShow: screenEnum.value > ScreenEnum.MD,
  },
]);
```

### 3. 国际化

```typescript
import { useI18n } from '@/hooks/web/useI18n';

const { t } = useI18n();

const columns = [
  {
    title: t('common.id'),
    key: 'id',
  },
  {
    title: t('user.username'),
    key: 'username',
  },
];
```

---

下一步：[表单组件详解](./form.md)

