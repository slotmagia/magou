# BasicForm 表单组件

BasicForm 是一个基于 Naive UI 封装的高级表单组件，提供了配置化的表单构建、验证、布局等功能，支持动态表单和复杂业务场景。

## 组件特性

### 🎯 核心功能
- **配置化构建** - 通过 JSON 配置快速构建表单
- **丰富的组件类型** - 支持所有 Naive UI 表单组件
- **灵活的布局** - 支持栅格布局和响应式设计
- **动态表单** - 支持动态添加/删除字段
- **表单验证** - 内置验证规则和自定义验证
- **插槽支持** - 支持自定义组件和复杂布局

### 📋 支持的组件类型
- **输入类**: NInput, NInputNumber, NInputPassword
- **选择类**: NSelect, NRadioGroup, NCheckboxGroup
- **日期类**: NDatePicker, NTimePicker
- **上传类**: NUpload
- **其他**: NSwitch, NSlider, NRate, NColorPicker

## 基础用法

### 简单表单

```vue
<template>
  <div class="p-4">
    <BasicForm
      @register="register"
      @submit="handleSubmit"
      @reset="handleReset"
    />
  </div>
</template>

<script setup lang="ts">
import { BasicForm, useForm, FormSchema } from '@/components/Form';

const schemas: FormSchema[] = [
  {
    field: 'username',
    label: '用户名',
    component: 'NInput',
    componentProps: {
      placeholder: '请输入用户名',
      clearable: true,
    },
    rules: [
      { required: true, message: '请输入用户名' },
      { min: 3, max: 20, message: '用户名长度在 3 到 20 个字符' },
    ],
  },
  {
    field: 'email',
    label: '邮箱',
    component: 'NInput',
    componentProps: {
      placeholder: '请输入邮箱',
      clearable: true,
    },
    rules: [
      { required: true, message: '请输入邮箱' },
      { type: 'email', message: '请输入正确的邮箱格式' },
    ],
  },
  {
    field: 'age',
    label: '年龄',
    component: 'NInputNumber',
    componentProps: {
      placeholder: '请输入年龄',
      min: 1,
      max: 120,
    },
  },
  {
    field: 'gender',
    label: '性别',
    component: 'NRadioGroup',
    componentProps: {
      options: [
        { label: '男', value: 1 },
        { label: '女', value: 2 },
      ],
    },
  },
];

const [register, { setFieldsValue, validate, resetFields, getFieldsValue }] = useForm({
  schemas,
  labelWidth: 100,
  showActionButtonGroup: true,
  submitButtonText: '提交',
  resetButtonText: '重置',
});

// 表单提交
const handleSubmit = async (values: any) => {
  console.log('表单数据:', values);
  try {
    // 这里处理提交逻辑
    window.$message.success('提交成功');
  } catch (error) {
    window.$message.error('提交失败');
  }
};

// 表单重置
const handleReset = () => {
  resetFields();
  window.$message.info('表单已重置');
};

// 设置表单默认值
onMounted(() => {
  setFieldsValue({
    username: 'admin',
    email: 'admin@example.com',
    age: 25,
    gender: 1,
  });
});
</script>
```

## 高级用法

### 动态表单

```vue
<template>
  <div class="p-4">
    <div class="mb-4">
      <n-button type="primary" @click="addField">添加字段</n-button>
      <n-button @click="removeField" class="ml-2">删除最后一个字段</n-button>
    </div>
    
    <BasicForm @register="register" />
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { BasicForm, useForm, FormSchema } from '@/components/Form';

const dynamicFields = ref<FormSchema[]>([]);
let fieldIndex = 0;

const baseSchemas: FormSchema[] = [
  {
    field: 'name',
    label: '基本信息',
    component: 'NInput',
    componentProps: {
      placeholder: '请输入名称',
    },
  },
];

const [register, { updateSchema, removeSchemaByField }] = useForm({
  schemas: baseSchemas,
  labelWidth: 120,
});

// 添加动态字段
const addField = () => {
  fieldIndex++;
  const newField: FormSchema = {
    field: `dynamic_${fieldIndex}`,
    label: `动态字段 ${fieldIndex}`,
    component: 'NInput',
    componentProps: {
      placeholder: `请输入动态字段 ${fieldIndex}`,
    },
  };
  
  dynamicFields.value.push(newField);
  updateSchema([...baseSchemas, ...dynamicFields.value]);
};

// 删除字段
const removeField = () => {
  if (dynamicFields.value.length > 0) {
    const removedField = dynamicFields.value.pop();
    if (removedField) {
      removeSchemaByField(removedField.field);
    }
  }
};
</script>
```

### 条件显示

```vue
<template>
  <BasicForm @register="register" />
</template>

<script setup lang="ts">
import { BasicForm, useForm, FormSchema } from '@/components/Form';

const schemas: FormSchema[] = [
  {
    field: 'userType',
    label: '用户类型',
    component: 'NSelect',
    componentProps: {
      options: [
        { label: '普通用户', value: 'normal' },
        { label: '管理员', value: 'admin' },
        { label: '超级管理员', value: 'super' },
      ],
    },
    rules: [{ required: true, message: '请选择用户类型' }],
  },
  {
    field: 'adminLevel',
    label: '管理员级别',
    component: 'NSelect',
    componentProps: {
      options: [
        { label: '初级管理员', value: 1 },
        { label: '中级管理员', value: 2 },
        { label: '高级管理员', value: 3 },
      ],
    },
    // 条件显示：只有当用户类型为管理员或超级管理员时才显示
    ifShow: ({ model }) => {
      return model.userType === 'admin' || model.userType === 'super';
    },
  },
  {
    field: 'permissions',
    label: '特殊权限',
    component: 'NCheckboxGroup',
    componentProps: {
      options: [
        { label: '系统配置', value: 'system:config' },
        { label: '用户管理', value: 'user:manage' },
        { label: '数据导出', value: 'data:export' },
      ],
    },
    // 只有超级管理员才能设置特殊权限
    ifShow: ({ model }) => model.userType === 'super',
  },
];

const [register] = useForm({
  schemas,
  labelWidth: 120,
});
</script>
```

### 自定义组件

```vue
<template>
  <BasicForm @register="register">
    <!-- 自定义组件插槽 -->
    <template #customUpload="{ model, field }">
      <div class="custom-upload">
        <n-upload
          action="/api/upload"
          :max="3"
          multiple
          directory-dnd
          @finish="handleUploadFinish"
        >
          <n-upload-dragger>
            <div style="margin-bottom: 12px">
              <n-icon size="48" :depth="3">
                <archive-icon />
              </n-icon>
            </div>
            <n-text style="font-size: 16px">
              点击或者拖动文件到该区域来上传
            </n-text>
            <n-p depth="3" style="margin: 8px 0 0 0">
              请不要上传敏感数据，比如你的银行卡号和密码，信用卡号有效期和安全码
            </n-p>
          </n-upload-dragger>
        </n-upload>
      </div>
    </template>
    
    <!-- 自定义验证组件 -->
    <template #customValidator="{ model, field }">
      <div class="flex items-center space-x-2">
        <n-input
          v-model:value="model[field]"
          placeholder="请输入自定义内容"
          @input="handleCustomValidation"
        />
        <n-button @click="validateCustomField">验证</n-button>
      </div>
    </template>
  </BasicForm>
</template>

<script setup lang="ts">
import { BasicForm, useForm, FormSchema } from '@/components/Form';
import { ArchiveOutlined as ArchiveIcon } from '@vicons/antd';

const schemas: FormSchema[] = [
  {
    field: 'title',
    label: '标题',
    component: 'NInput',
    componentProps: {
      placeholder: '请输入标题',
    },
  },
  {
    field: 'files',
    label: '文件上传',
    slot: 'customUpload', // 使用自定义插槽
  },
  {
    field: 'customField',
    label: '自定义字段',
    slot: 'customValidator',
  },
];

const [register, { getFieldsValue, setFieldsValue }] = useForm({
  schemas,
  labelWidth: 120,
});

// 处理上传完成
const handleUploadFinish = ({ file, event }) => {
  console.log('上传完成:', file, event);
  // 处理上传结果
};

// 自定义验证
const handleCustomValidation = (value: string) => {
  // 实时验证逻辑
  console.log('自定义验证:', value);
};

// 验证自定义字段
const validateCustomField = () => {
  const values = getFieldsValue();
  console.log('当前表单值:', values);
  // 执行验证逻辑
};
</script>
```

## API 接口

### BasicForm Props

| 属性 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| schemas | `FormSchema[]` | `[]` | 表单配置项 |
| model | `object` | `{}` | 表单数据 |
| labelWidth | `number \| string` | `'auto'` | 标签宽度 |
| labelPlacement | `'left' \| 'top'` | `'left'` | 标签位置 |
| size | `'small' \| 'medium' \| 'large'` | `'medium'` | 组件大小 |
| inline | `boolean` | `false` | 是否内联模式 |
| showActionButtonGroup | `boolean` | `true` | 是否显示操作按钮组 |
| showSubmitButton | `boolean` | `true` | 是否显示提交按钮 |
| showResetButton | `boolean` | `true` | 是否显示重置按钮 |
| submitButtonText | `string` | `'提交'` | 提交按钮文本 |
| resetButtonText | `string` | `'重置'` | 重置按钮文本 |
| submitButtonOptions | `ButtonProps` | `{}` | 提交按钮配置 |
| resetButtonOptions | `ButtonProps` | `{}` | 重置按钮配置 |
| gridProps | `GridProps` | `{}` | 栅格布局配置 |
| collapsedRows | `number` | `1` | 收起时显示的行数 |

### FormSchema 配置

| 属性 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| field | `string` | - | 字段名，必填 |
| label | `string` | - | 标签文本 |
| labelMessage | `string` | - | 标签提示信息 |
| component | `ComponentType` | - | 组件类型 |
| componentProps | `object` | `{}` | 组件属性 |
| slot | `string` | - | 自定义插槽名 |
| rules | `FormItemRule[]` | `[]` | 验证规则 |
| defaultValue | `any` | - | 默认值 |
| giProps | `GridItemProps` | `{}` | 栅格项配置 |
| isFull | `boolean` | `false` | 是否占满一行 |
| suffix | `string` | - | 字段后缀 |
| ifShow | `boolean \| Function` | `true` | 是否显示 |
| auth | `string[]` | - | 权限控制 |

### useForm Hook

```typescript
const [register, methods] = useForm(props);
```

#### 返回方法

| 方法名 | 参数 | 返回值 | 说明 |
|--------|------|--------|------|
| setFieldsValue | `(values: object)` | `Promise<void>` | 设置表单值 |
| getFieldsValue | - | `object` | 获取表单值 |
| resetFields | - | `Promise<void>` | 重置表单 |
| validate | `(nameList?: string[])` | `Promise<object>` | 验证表单 |
| clearValidate | `(name?: string \| string[])` | `Promise<void>` | 清除验证 |
| setProps | `(props: FormProps)` | `Promise<void>` | 设置表单属性 |
| updateSchema | `(schema: FormSchema[])` | `void` | 更新表单配置 |
| appendSchemaByField | `(schema: FormSchema, prefixField?: string)` | `void` | 追加字段 |
| removeSchemaByField | `(field: string)` | `void` | 删除字段 |

## 组件类型支持

### 输入组件

```typescript
// 文本输入
{
  field: 'text',
  label: '文本输入',
  component: 'NInput',
  componentProps: {
    type: 'text',
    placeholder: '请输入文本',
    clearable: true,
    maxlength: 100,
    showCount: true,
  },
}

// 数字输入
{
  field: 'number',
  label: '数字输入',
  component: 'NInputNumber',
  componentProps: {
    min: 0,
    max: 999999,
    step: 1,
    precision: 2,
    placeholder: '请输入数字',
  },
}

// 密码输入
{
  field: 'password',
  label: '密码',
  component: 'NInput',
  componentProps: {
    type: 'password',
    showPasswordOn: 'mousedown',
    placeholder: '请输入密码',
  },
}

// 文本域
{
  field: 'textarea',
  label: '多行文本',
  component: 'NInput',
  componentProps: {
    type: 'textarea',
    placeholder: '请输入多行文本',
    rows: 4,
    maxlength: 500,
    showCount: true,
  },
}
```

### 选择组件

```typescript
// 下拉选择
{
  field: 'select',
  label: '下拉选择',
  component: 'NSelect',
  componentProps: {
    placeholder: '请选择',
    options: [
      { label: '选项1', value: 1 },
      { label: '选项2', value: 2 },
    ],
    clearable: true,
    filterable: true,
  },
}

// 多选
{
  field: 'multiSelect',
  label: '多选',
  component: 'NSelect',
  componentProps: {
    placeholder: '请选择',
    multiple: true,
    options: [
      { label: '选项1', value: 1 },
      { label: '选项2', value: 2 },
      { label: '选项3', value: 3 },
    ],
  },
}

// 单选组
{
  field: 'radio',
  label: '单选',
  component: 'NRadioGroup',
  componentProps: {
    options: [
      { label: '选项1', value: 1 },
      { label: '选项2', value: 2 },
    ],
  },
}

// 复选组
{
  field: 'checkbox',
  label: '多选',
  component: 'NCheckboxGroup',
  componentProps: {
    options: [
      { label: '选项1', value: 1 },
      { label: '选项2', value: 2 },
      { label: '选项3', value: 3 },
    ],
  },
}
```

### 日期时间组件

```typescript
// 日期选择
{
  field: 'date',
  label: '日期',
  component: 'NDatePicker',
  componentProps: {
    type: 'date',
    clearable: true,
    placeholder: '请选择日期',
  },
}

// 日期时间选择
{
  field: 'datetime',
  label: '日期时间',
  component: 'NDatePicker',
  componentProps: {
    type: 'datetime',
    clearable: true,
    placeholder: '请选择日期时间',
  },
}

// 日期范围选择
{
  field: 'daterange',
  label: '日期范围',
  component: 'NDatePicker',
  componentProps: {
    type: 'daterange',
    clearable: true,
    placeholder: ['开始日期', '结束日期'],
  },
}

// 时间选择
{
  field: 'time',
  label: '时间',
  component: 'NTimePicker',
  componentProps: {
    clearable: true,
    placeholder: '请选择时间',
  },
}
```

### 其他组件

```typescript
// 开关
{
  field: 'switch',
  label: '开关',
  component: 'NSwitch',
  componentProps: {
    checkedValue: true,
    uncheckedValue: false,
  },
}

// 滑块
{
  field: 'slider',
  label: '滑块',
  component: 'NSlider',
  componentProps: {
    min: 0,
    max: 100,
    step: 1,
    marks: {
      0: '0',
      50: '50',
      100: '100',
    },
  },
}

// 评分
{
  field: 'rate',
  label: '评分',
  component: 'NRate',
  componentProps: {
    allowHalf: true,
    count: 5,
  },
}

// 颜色选择
{
  field: 'color',
  label: '颜色',
  component: 'NColorPicker',
  componentProps: {
    showAlpha: false,
  },
}
```

## 表单验证

### 内置验证规则

```typescript
const schemas: FormSchema[] = [
  {
    field: 'required',
    label: '必填字段',
    component: 'NInput',
    rules: [
      { required: true, message: '这是必填字段' },
    ],
  },
  {
    field: 'email',
    label: '邮箱',
    component: 'NInput',
    rules: [
      { required: true, message: '请输入邮箱' },
      { type: 'email', message: '请输入正确的邮箱格式' },
    ],
  },
  {
    field: 'phone',
    label: '手机号',
    component: 'NInput',
    rules: [
      { required: true, message: '请输入手机号' },
      { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号格式' },
    ],
  },
  {
    field: 'range',
    label: '数值范围',
    component: 'NInputNumber',
    rules: [
      { required: true, message: '请输入数值' },
      { type: 'number', min: 1, max: 100, message: '数值范围在 1 到 100 之间' },
    ],
  },
];
```

### 自定义验证

```typescript
const customValidator = (rule: any, value: any) => {
  return new Promise((resolve, reject) => {
    if (!value) {
      reject(new Error('请输入内容'));
    } else if (value.length < 6) {
      reject(new Error('长度不能少于6位'));
    } else {
      resolve();
    }
  });
};

const schemas: FormSchema[] = [
  {
    field: 'custom',
    label: '自定义验证',
    component: 'NInput',
    rules: [
      { validator: customValidator, trigger: 'blur' },
    ],
  },
];
```

### 异步验证

```typescript
const asyncValidator = async (rule: any, value: any) => {
  if (!value) {
    throw new Error('请输入用户名');
  }
  
  // 模拟异步验证
  const response = await checkUsername(value);
  if (!response.available) {
    throw new Error('用户名已存在');
  }
};

const schemas: FormSchema[] = [
  {
    field: 'username',
    label: '用户名',
    component: 'NInput',
    rules: [
      { validator: asyncValidator, trigger: 'blur' },
    ],
  },
];
```

## 布局配置

### 栅格布局

```typescript
const [register] = useForm({
  schemas,
  gridProps: {
    cols: 24, // 总列数
    xGap: 16, // 水平间距
    yGap: 16, // 垂直间距
  },
  // 每个字段的栅格配置
  giProps: {
    span: 12, // 默认占12列（一半宽度）
  },
});

// 单独配置某个字段的栅格
const schemas: FormSchema[] = [
  {
    field: 'title',
    label: '标题',
    component: 'NInput',
    giProps: {
      span: 24, // 占满整行
    },
  },
  {
    field: 'name',
    label: '姓名',
    component: 'NInput',
    giProps: {
      span: 8, // 占1/3宽度
    },
  },
  {
    field: 'age',
    label: '年龄',
    component: 'NInputNumber',
    giProps: {
      span: 8, // 占1/3宽度
    },
  },
  {
    field: 'gender',
    label: '性别',
    component: 'NSelect',
    giProps: {
      span: 8, // 占1/3宽度
    },
  },
];
```

### 响应式布局

```typescript
const [register] = useForm({
  schemas,
  gridProps: {
    cols: '1 s:1 m:2 l:3 xl:4 2xl:4',
    responsive: 'screen',
    xGap: 16,
    yGap: 16,
  },
});
```

## 最佳实践

### 1. 表单性能优化

```typescript
// 使用计算属性优化大量选项
const userOptions = computed(() => {
  return userList.value.map(user => ({
    label: user.name,
    value: user.id,
  }));
});

// 异步加载选项
const loadUserOptions = async () => {
  const response = await getUserList();
  return response.data.map(user => ({
    label: user.name,
    value: user.id,
  }));
};
```

### 2. 表单状态管理

```typescript
// 使用 Pinia 管理表单状态
const useFormStore = defineStore('form', () => {
  const formData = ref({});
  
  const setFormData = (data: any) => {
    formData.value = { ...formData.value, ...data };
  };
  
  const resetFormData = () => {
    formData.value = {};
  };
  
  return {
    formData,
    setFormData,
    resetFormData,
  };
});
```

### 3. 表单复用

```typescript
// 创建可复用的表单配置
export const createUserFormSchemas = (type: 'create' | 'edit'): FormSchema[] => [
  {
    field: 'username',
    label: '用户名',
    component: 'NInput',
    componentProps: {
      disabled: type === 'edit', // 编辑时禁用用户名
    },
    rules: [{ required: true, message: '请输入用户名' }],
  },
  // ... 其他字段
];
```

---

下一步：[模态框组件](./modal.md)






