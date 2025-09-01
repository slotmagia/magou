# BasicModal 模态框组件

BasicModal 是基于 Naive UI Modal 组件封装的高级模态框，提供了拖拽、全屏、加载状态等增强功能，支持灵活的内容和操作配置。

## 组件特性

### 🎯 核心功能
- **拖拽支持** - 支持拖拽移动模态框
- **多种尺寸** - 预设多种常用尺寸
- **全屏模式** - 支持全屏显示
- **加载状态** - 内置提交按钮加载状态
- **灵活配置** - 支持自定义头部、内容、底部
- **键盘交互** - 支持 ESC 关闭、Enter 确认

### 📐 内置尺寸
- **small**: 400px
- **medium**: 600px  
- **large**: 800px
- **extra-large**: 1000px

## 基础用法

### 简单模态框

```vue
<template>
  <div class="p-4">
    <n-button type="primary" @click="openModal">
      打开模态框
    </n-button>
    
    <BasicModal
      @register="register"
      title="基础模态框"
      @ok="handleOk"
      @close="handleClose"
    >
      <div class="py-4">
        <p>这是模态框的内容区域。</p>
        <p>您可以在这里放置任何内容。</p>
      </div>
    </BasicModal>
  </div>
</template>

<script setup lang="ts">
import { BasicModal, useModal } from '@/components/Modal';

const [register, { openModal, closeModal, setModalProps }] = useModal();

// 打开模态框时的处理
const handleOk = () => {
  console.log('用户点击了确认按钮');
  // 这里可以进行表单验证、数据提交等操作
  closeModal();
};

// 关闭模态框时的处理
const handleClose = () => {
  console.log('模态框已关闭');
};
</script>
```

### 带表单的模态框

```vue
<template>
  <div class="p-4">
    <n-button type="primary" @click="handleCreate">
      新增用户
    </n-button>
    <n-button @click="handleEdit" class="ml-2">
      编辑用户
    </n-button>
    
    <BasicModal
      @register="register"
      :title="modalTitle"
      :width="600"
      @ok="handleSubmit"
    >
      <BasicForm @register="formRegister" />
    </BasicModal>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import { BasicModal, useModal } from '@/components/Modal';
import { BasicForm, useForm, FormSchema } from '@/components/Form';

const isEdit = ref(false);
const currentUser = ref(null);

const modalTitle = computed(() => {
  return isEdit.value ? '编辑用户' : '新增用户';
});

// 表单配置
const schemas: FormSchema[] = [
  {
    field: 'username',
    label: '用户名',
    component: 'NInput',
    componentProps: {
      placeholder: '请输入用户名',
    },
    rules: [{ required: true, message: '请输入用户名' }],
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
      { type: 'email', message: '请输入正确的邮箱格式' },
    ],
  },
  {
    field: 'role',
    label: '角色',
    component: 'NSelect',
    componentProps: {
      options: [
        { label: '普通用户', value: 'user' },
        { label: '管理员', value: 'admin' },
      ],
    },
    rules: [{ required: true, message: '请选择角色' }],
  },
];

const [register, { openModal, closeModal, setModalProps }] = useModal();
const [formRegister, { validate, setFieldsValue, resetFields }] = useForm({
  schemas,
  showActionButtonGroup: false, // 在模态框中不显示表单的操作按钮
});

// 新增用户
const handleCreate = () => {
  isEdit.value = false;
  currentUser.value = null;
  resetFields();
  openModal();
};

// 编辑用户
const handleEdit = () => {
  isEdit.value = true;
  currentUser.value = {
    id: 1,
    username: 'admin',
    email: 'admin@example.com',
    role: 'admin',
  };
  
  setFieldsValue(currentUser.value);
  openModal();
};

// 提交表单
const handleSubmit = async () => {
  try {
    const values = await validate();
    
    // 模拟提交
    console.log('提交数据:', values);
    
    if (isEdit.value) {
      // 更新用户
      await updateUser(currentUser.value.id, values);
      window.$message.success('用户更新成功');
    } else {
      // 创建用户
      await createUser(values);
      window.$message.success('用户创建成功');
    }
    
    closeModal();
  } catch (error) {
    console.error('提交失败:', error);
    window.$message.error('操作失败');
  }
};

// 模拟 API 调用
const createUser = async (userData: any) => {
  return new Promise(resolve => setTimeout(resolve, 1000));
};

const updateUser = async (id: number, userData: any) => {
  return new Promise(resolve => setTimeout(resolve, 1000));
};
</script>
```

## 高级用法

### 自定义操作按钮

```vue
<template>
  <BasicModal @register="register" title="自定义操作">
    <div class="py-4">
      <p>这个模态框有自定义的操作按钮。</p>
    </div>
    
    <!-- 自定义操作区域 -->
    <template #action>
      <n-space>
        <n-button @click="handleCancel">
          取消
        </n-button>
        <n-button type="warning" @click="handleSaveAsDraft">
          保存草稿
        </n-button>
        <n-button type="primary" @click="handlePublish" :loading="publishing">
          发布
        </n-button>
      </n-space>
    </template>
  </BasicModal>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { BasicModal, useModal } from '@/components/Modal';

const publishing = ref(false);

const [register, { closeModal }] = useModal();

const handleCancel = () => {
  closeModal();
};

const handleSaveAsDraft = async () => {
  try {
    // 保存草稿逻辑
    await saveDraft();
    window.$message.success('草稿保存成功');
    closeModal();
  } catch (error) {
    window.$message.error('保存失败');
  }
};

const handlePublish = async () => {
  try {
    publishing.value = true;
    // 发布逻辑
    await publish();
    window.$message.success('发布成功');
    closeModal();
  } catch (error) {
    window.$message.error('发布失败');
  } finally {
    publishing.value = false;
  }
};

const saveDraft = () => new Promise(resolve => setTimeout(resolve, 1000));
const publish = () => new Promise(resolve => setTimeout(resolve, 2000));
</script>
```

### 全屏模态框

```vue
<template>
  <div class="p-4">
    <n-button type="primary" @click="openFullscreenModal">
      打开全屏模态框
    </n-button>
    
    <BasicModal
      @register="register"
      title="全屏模态框"
      :width="'100%'"
      :height="'100%'"
      :mask-closable="false"
      :closable="true"
    >
      <div class="h-full p-4">
        <div class="h-full border-2 border-dashed border-gray-300 rounded-lg flex items-center justify-center">
          <div class="text-center">
            <h3 class="text-xl font-semibold mb-4">全屏内容区域</h3>
            <p class="text-gray-600 mb-4">这里可以放置复杂的内容，比如图表、编辑器等</p>
            <n-button @click="toggleFullscreen">
              切换全屏状态
            </n-button>
          </div>
        </div>
      </div>
      
      <template #action>
        <n-space>
          <n-button @click="closeModal">关闭</n-button>
          <n-button type="primary">保存</n-button>
        </n-space>
      </template>
    </BasicModal>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { BasicModal, useModal } from '@/components/Modal';

const isFullscreen = ref(false);

const [register, { openModal, closeModal, setModalProps }] = useModal();

const openFullscreenModal = () => {
  isFullscreen.value = true;
  setModalProps({
    width: '100%',
    height: '100%',
    style: {
      margin: 0,
      padding: 0,
    },
  });
  openModal();
};

const toggleFullscreen = () => {
  isFullscreen.value = !isFullscreen.value;
  
  if (isFullscreen.value) {
    setModalProps({
      width: '100%',
      height: '100%',
      style: { margin: 0, padding: 0 },
    });
  } else {
    setModalProps({
      width: 800,
      height: 600,
      style: { margin: 'auto', padding: '20px' },
    });
  }
};
</script>
```

### 嵌套模态框

```vue
<template>
  <div class="p-4">
    <n-button type="primary" @click="openParentModal">
      打开父模态框
    </n-button>
    
    <!-- 父模态框 -->
    <BasicModal
      @register="parentRegister"
      title="父模态框"
      :width="600"
    >
      <div class="py-4">
        <p class="mb-4">这是父模态框的内容。</p>
        <n-button type="primary" @click="openChildModal">
          打开子模态框
        </n-button>
      </div>
    </BasicModal>
    
    <!-- 子模态框 -->
    <BasicModal
      @register="childRegister"
      title="子模态框"
      :width="400"
      :z-index="2000"
    >
      <div class="py-4">
        <p>这是子模态框的内容。</p>
        <p class="text-sm text-gray-500">子模态框的 z-index 更高</p>
      </div>
    </BasicModal>
  </div>
</template>

<script setup lang="ts">
import { BasicModal, useModal } from '@/components/Modal';

const [parentRegister, { openModal: openParentModal }] = useModal();
const [childRegister, { openModal: openChildModal }] = useModal();
</script>
```

## API 接口

### BasicModal Props

| 属性 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| title | `string` | - | 模态框标题 |
| width | `number \| string` | `600` | 模态框宽度 |
| height | `number \| string` | `'auto'` | 模态框高度 |
| minWidth | `number` | `260` | 最小宽度 |
| minHeight | `number` | `200` | 最小高度 |
| draggable | `boolean` | `true` | 是否可拖拽 |
| resizable | `boolean` | `false` | 是否可调整大小 |
| maskClosable | `boolean` | `true` | 点击遮罩是否关闭 |
| closable | `boolean` | `true` | 是否显示关闭按钮 |
| showIcon | `boolean` | `true` | 是否显示图标 |
| subBtuText | `string` | `'确定'` | 确认按钮文本 |
| canFullscreen | `boolean` | `true` | 是否支持全屏 |
| defaultFullscreen | `boolean` | `false` | 是否默认全屏 |

### ModalMethods

| 方法名 | 参数 | 说明 |
|--------|------|------|
| setProps | `(props: ModalProps)` | 设置模态框属性 |
| openModal | - | 打开模态框 |
| closeModal | - | 关闭模态框 |
| setSubLoading | `(loading: boolean)` | 设置确认按钮加载状态 |

### useModal Hook

```typescript
const [register, methods] = useModal(props?);
```

返回一个注册函数和方法对象。

## 事件

| 事件名 | 参数 | 说明 |
|--------|------|------|
| register | `(instance: ModalMethods)` | 注册模态框实例 |
| ok | - | 点击确认按钮 |
| cancel | - | 点击取消按钮 |
| close | - | 模态框关闭 |

## 插槽

| 插槽名 | 说明 |
|--------|------|
| default | 模态框内容 |
| header | 自定义头部 |
| action | 自定义操作区域 |

## 样式定制

### CSS 变量

```css
.basic-modal {
  --modal-border-radius: 8px;
  --modal-header-height: 54px;
  --modal-footer-height: 60px;
  --modal-padding: 20px;
}
```

### 自定义样式

```vue
<template>
  <BasicModal
    @register="register"
    title="自定义样式"
    class="custom-modal"
  >
    <div class="custom-content">
      自定义内容
    </div>
  </BasicModal>
</template>

<style scoped>
.custom-modal {
  :deep(.n-modal) {
    border-radius: 12px;
    box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.25);
  }
  
  :deep(.n-modal__header) {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: white;
    border-radius: 12px 12px 0 0;
  }
}

.custom-content {
  padding: 20px;
  min-height: 200px;
}
</style>
```

## 性能优化

### 懒加载模态框内容

```vue
<template>
  <BasicModal @register="register" title="懒加载内容">
    <div v-if="contentLoaded">
      <HeavyComponent :data="modalData" />
    </div>
    <div v-else class="text-center py-8">
      <n-spin size="large" />
      <p class="mt-4 text-gray-500">加载中...</p>
    </div>
  </BasicModal>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { BasicModal, useModal } from '@/components/Modal';

const contentLoaded = ref(false);
const modalData = ref(null);

const [register, { openModal }] = useModal();

// 重写 openModal 方法，添加数据加载逻辑
const openModalWithData = async () => {
  openModal();
  
  try {
    // 加载数据
    modalData.value = await loadModalData();
    contentLoaded.value = true;
  } catch (error) {
    window.$message.error('数据加载失败');
  }
};

const loadModalData = () => {
  return new Promise(resolve => {
    setTimeout(() => {
      resolve({ message: '这是懒加载的数据' });
    }, 2000);
  });
};

// 定义重型组件
const HeavyComponent = defineAsyncComponent(() => import('./HeavyComponent.vue'));
</script>
```

### 缓存模态框实例

```typescript
// composables/useModalCache.ts
const modalCache = new Map();

export function useModalCache(key: string) {
  const getCachedModal = () => modalCache.get(key);
  
  const setCachedModal = (modalInstance: any) => {
    modalCache.set(key, modalInstance);
  };
  
  const clearCache = () => {
    modalCache.delete(key);
  };
  
  return {
    getCachedModal,
    setCachedModal,
    clearCache,
  };
}
```

## 最佳实践

### 1. 模态框状态管理

```typescript
// stores/modalStore.ts
export const useModalStore = defineStore('modal', () => {
  const modals = ref<Record<string, boolean>>({});
  
  const openModal = (key: string) => {
    modals.value[key] = true;
  };
  
  const closeModal = (key: string) => {
    modals.value[key] = false;
  };
  
  const isModalOpen = (key: string) => {
    return modals.value[key] || false;
  };
  
  return {
    modals,
    openModal,
    closeModal,
    isModalOpen,
  };
});
```

### 2. 模态框组合

```vue
<script setup lang="ts">
// 组合多个模态框功能
const useUserModal = () => {
  const [register, modalMethods] = useModal();
  const [formRegister, formMethods] = useForm({
    schemas: userFormSchemas,
  });
  
  const openCreateModal = () => {
    formMethods.resetFields();
    modalMethods.setModalProps({ title: '新增用户' });
    modalMethods.openModal();
  };
  
  const openEditModal = (user: User) => {
    formMethods.setFieldsValue(user);
    modalMethods.setModalProps({ title: '编辑用户' });
    modalMethods.openModal();
  };
  
  return {
    register,
    formRegister,
    openCreateModal,
    openEditModal,
    ...modalMethods,
    ...formMethods,
  };
};
</script>
```

### 3. 模态框权限控制

```vue
<template>
  <BasicModal
    @register="register"
    :title="modalTitle"
    v-if="hasPermission"
  >
    <!-- 模态框内容 -->
  </BasicModal>
</template>

<script setup lang="ts">
import { usePermission } from '@/hooks/web/usePermission';

const { hasPermission } = usePermission();

const canEdit = computed(() => {
  return hasPermission(['user:edit']);
});

const modalTitle = computed(() => {
  return canEdit.value ? '编辑用户' : '查看用户';
});
</script>
```

---

下一步：[上传组件](./upload.md)






