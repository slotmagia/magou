# 其他组件

除了核心的表格、表单、模态框和上传组件外，HotGo 2.0 还提供了一系列实用的辅助组件，涵盖应用容器、图标选择、内容加载等功能。

## Application 应用容器

### 组件介绍

Application 组件是整个应用的容器组件，负责提供全局的 Naive UI 组件上下文和主题配置。

```vue
<template>
  <n-config-provider
    :theme="darkTheme"
    :theme-overrides="themeOverrides"
    :locale="zhCN"
    :date-locale="dateZhCN"
  >
    <n-loading-bar-provider>
      <n-dialog-provider>
        <n-notification-provider>
          <n-message-provider>
            <slot />
          </n-message-provider>
        </n-notification-provider>
      </n-dialog-provider>
    </n-loading-bar-provider>
  </n-config-provider>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import {
  darkTheme,
  zhCN,
  dateZhCN,
  NConfigProvider,
  NLoadingBarProvider,
  NDialogProvider,
  NNotificationProvider,
  NMessageProvider,
} from 'naive-ui';
import { useAppStore } from '@/store/modules/app';

const appStore = useAppStore();

// 主题配置
const themeOverrides = computed(() => ({
  common: {
    primaryColor: '#1890ff',
    primaryColorHover: '#40a9ff',
    primaryColorPressed: '#096dd9',
    borderRadius: '6px',
  },
  Button: {
    borderRadiusMedium: '6px',
  },
  Input: {
    borderRadius: '6px',
  },
  Card: {
    borderRadius: '8px',
  },
}));
</script>
```

### 使用方式

```typescript
// main.ts
import { createApp } from 'vue';
import App from './App.vue';
import { AppProvider } from '@/components/Application';

const appProvider = createApp(AppProvider);
const app = createApp(App);

// 先挂载 Provider
appProvider.mount('#appProvider', true);

// 再挂载主应用
app.mount('#app', true);
```

## CitySelector 城市选择器

### 组件特性

- **三级联动** - 省份、城市、区县三级联动选择
- **数据懒加载** - 按需加载城市数据
- **搜索功能** - 支持城市名称搜索
- **定位功能** - 支持获取当前位置

### 基础用法

```vue
<template>
  <div class="p-4">
    <div class="mb-4">
      <h3>基础城市选择</h3>
      <CitySelector
        v-model:value="selectedCity"
        placeholder="请选择城市"
        @change="handleCityChange"
      />
    </div>
    
    <div class="mb-4">
      <h3>带搜索的城市选择</h3>
      <CitySelector
        v-model:value="searchCity"
        :searchable="true"
        :show-path="true"
        placeholder="搜索或选择城市"
      />
    </div>
    
    <div class="mb-4">
      <h3>多选城市</h3>
      <CitySelector
        v-model:value="multipleCities"
        :multiple="true"
        :max-tag-count="3"
        placeholder="选择多个城市"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { CitySelector } from '@/components/CitySelector';

const selectedCity = ref('');
const searchCity = ref('');
const multipleCities = ref<string[]>([]);

const handleCityChange = (value: string, option: any) => {
  console.log('选择的城市:', value, option);
};
</script>
```

### API 接口

| 属性 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| modelValue | `string \| string[]` | - | 绑定值 |
| placeholder | `string` | `'请选择'` | 占位文本 |
| searchable | `boolean` | `false` | 是否可搜索 |
| multiple | `boolean` | `false` | 是否多选 |
| showPath | `boolean` | `false` | 是否显示完整路径 |
| maxTagCount | `number` | - | 多选时最大显示标签数 |
| level | `1 \| 2 \| 3` | `3` | 选择级别 |
| disabled | `boolean` | `false` | 是否禁用 |

## IconSelector 图标选择器

### 组件特性

- **多图标库支持** - 支持 Ant Design、Ionicons 等图标库
- **分类浏览** - 按分类组织图标
- **搜索功能** - 支持图标名称搜索
- **预览功能** - 实时预览选中图标

### 基础用法

```vue
<template>
  <div class="p-4">
    <div class="mb-4">
      <h3>Ant Design 图标选择</h3>
      <IconSelector
        v-model:value="antdIcon"
        type="antd"
        @change="handleIconChange"
      />
      <div class="mt-2">
        选中的图标: {{ antdIcon }}
        <n-icon v-if="antdIcon" :component="getAntdIcon(antdIcon)" size="20" class="ml-2" />
      </div>
    </div>
    
    <div class="mb-4">
      <h3>Ionicons 图标选择</h3>
      <IconSelector
        v-model:value="ionIcon"
        type="ionicons"
        :show-preview="true"
      />
    </div>
    
    <div class="mb-4">
      <h3>自定义图标列表</h3>
      <IconSelector
        v-model:value="customIcon"
        :icons="customIcons"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { IconSelector } from '@/components/IconSelector';
import * as AntdIcons from '@vicons/antd';

const antdIcon = ref('');
const ionIcon = ref('');
const customIcon = ref('');

const customIcons = [
  { name: 'custom-icon-1', label: '自定义图标1' },
  { name: 'custom-icon-2', label: '自定义图标2' },
];

const handleIconChange = (iconName: string) => {
  console.log('选择的图标:', iconName);
};

const getAntdIcon = (iconName: string) => {
  return AntdIcons[iconName as keyof typeof AntdIcons];
};
</script>
```

### API 接口

| 属性 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| modelValue | `string` | - | 绑定的图标名称 |
| type | `'antd' \| 'ionicons' \| 'custom'` | `'antd'` | 图标库类型 |
| icons | `IconItem[]` | - | 自定义图标列表 |
| showPreview | `boolean` | `true` | 是否显示预览 |
| searchable | `boolean` | `true` | 是否可搜索 |
| size | `number` | `24` | 图标显示大小 |

## CountTo 数字动画

### 组件特性

- **平滑动画** - 数字递增动画效果
- **自定义格式** - 支持千分位、小数点等格式化
- **动画控制** - 可控制动画时长和缓动函数
- **前缀后缀** - 支持添加前缀和后缀

### 基础用法

```vue
<template>
  <div class="p-4">
    <div class="grid grid-cols-2 gap-6">
      <n-card title="基础数字动画">
        <CountTo
          :start-val="0"
          :end-val="2023"
          :duration="2000"
          class="text-2xl font-bold text-blue-500"
        />
      </n-card>
      
      <n-card title="金额动画">
        <CountTo
          :start-val="0"
          :end-val="123456.78"
          :duration="3000"
          :decimals="2"
          prefix="¥"
          :separator="true"
          class="text-2xl font-bold text-green-500"
        />
      </n-card>
      
      <n-card title="百分比动画">
        <CountTo
          :start-val="0"
          :end-val="85.5"
          :duration="2500"
          :decimals="1"
          suffix="%"
          class="text-2xl font-bold text-orange-500"
        />
      </n-card>
      
      <n-card title="自定义格式">
        <CountTo
          :start-val="0"
          :end-val="999999"
          :duration="3000"
          prefix="用户数: "
          suffix=" 人"
          :separator="true"
          class="text-2xl font-bold text-purple-500"
        />
      </n-card>
    </div>
    
    <div class="mt-6">
      <n-button @click="startAnimation" type="primary">
        重新开始动画
      </n-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { CountTo } from '@/components/CountTo';

const animationKey = ref(0);

const startAnimation = () => {
  animationKey.value++;
};
</script>
```

### API 接口

| 属性 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| startVal | `number` | `0` | 起始值 |
| endVal | `number` | - | 结束值 |
| duration | `number` | `2000` | 动画时长(ms) |
| autoplay | `boolean` | `true` | 是否自动播放 |
| decimals | `number` | `0` | 小数位数 |
| prefix | `string` | - | 前缀 |
| suffix | `string` | - | 后缀 |
| separator | `boolean` | `false` | 是否使用千分位分隔符 |
| useEasing | `boolean` | `true` | 是否使用缓动效果 |

## LoadingContent 加载内容

### 组件特性

- **多种加载状态** - 支持不同的加载样式
- **骨架屏** - 提供内容骨架屏效果
- **空状态** - 统一的空数据展示
- **错误状态** - 错误信息展示

### 基础用法

```vue
<template>
  <div class="p-4">
    <div class="mb-4">
      <n-space>
        <n-button @click="loadingState = 'loading'">加载中</n-button>
        <n-button @click="loadingState = 'empty'">空数据</n-button>
        <n-button @click="loadingState = 'error'">错误状态</n-button>
        <n-button @click="loadingState = 'success'">成功状态</n-button>
      </n-space>
    </div>
    
    <LoadingContent
      :loading="loadingState === 'loading'"
      :empty="loadingState === 'empty'"
      :error="loadingState === 'error'"
      :skeleton="true"
      :skeleton-rows="5"
      @retry="handleRetry"
    >
      <div v-if="loadingState === 'success'">
        <n-card title="内容区域">
          <p>这里是正常的内容展示区域。</p>
          <p>当数据加载完成后会显示这些内容。</p>
        </n-card>
      </div>
    </LoadingContent>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { LoadingContent } from '@/components/LoadingContent';

const loadingState = ref<'loading' | 'empty' | 'error' | 'success'>('loading');

const handleRetry = () => {
  loadingState.value = 'loading';
  
  // 模拟重新加载
  setTimeout(() => {
    loadingState.value = 'success';
  }, 2000);
};
</script>
```

### API 接口

| 属性 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| loading | `boolean` | `false` | 是否加载中 |
| empty | `boolean` | `false` | 是否空数据 |
| error | `boolean` | `false` | 是否错误状态 |
| skeleton | `boolean` | `false` | 是否显示骨架屏 |
| skeletonRows | `number` | `3` | 骨架屏行数 |
| emptyText | `string` | `'暂无数据'` | 空数据提示文本 |
| errorText | `string` | `'加载失败'` | 错误提示文本 |
| showRetry | `boolean` | `true` | 是否显示重试按钮 |

## SvgIcon SVG 图标

### 组件特性

- **SVG 支持** - 完整的 SVG 图标支持
- **自定义颜色** - 支持自定义图标颜色
- **响应式大小** - 支持响应式图标大小
- **旋转动画** - 支持图标旋转动画

### 基础用法

```vue
<template>
  <div class="p-4">
    <div class="mb-4">
      <h3>基础 SVG 图标</h3>
      <n-space>
        <SvgIcon name="user" size="24" />
        <SvgIcon name="setting" size="32" color="#1890ff" />
        <SvgIcon name="loading" size="28" :spin="true" />
      </n-space>
    </div>
    
    <div class="mb-4">
      <h3>不同尺寸</h3>
      <n-space>
        <SvgIcon name="heart" size="16" />
        <SvgIcon name="heart" size="24" />
        <SvgIcon name="heart" size="32" />
        <SvgIcon name="heart" size="48" />
      </n-space>
    </div>
    
    <div class="mb-4">
      <h3>不同颜色</h3>
      <n-space>
        <SvgIcon name="star" size="32" color="#ff4d4f" />
        <SvgIcon name="star" size="32" color="#52c41a" />
        <SvgIcon name="star" size="32" color="#faad14" />
        <SvgIcon name="star" size="32" color="#1890ff" />
      </n-space>
    </div>
  </div>
</template>

<script setup lang="ts">
import { SvgIcon } from '@/components/SvgIcon';
</script>
```

### API 接口

| 属性 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| name | `string` | - | 图标名称 |
| size | `number \| string` | `16` | 图标大小 |
| color | `string` | `'currentColor'` | 图标颜色 |
| spin | `boolean` | `false` | 是否旋转 |
| prefix | `string` | `'icon'` | 图标前缀 |

## ComplexMemberPicker 复合成员选择器

### 组件特性

- **多类型选择** - 支持用户、部门、角色等多种类型
- **树形结构** - 支持组织架构树形展示
- **搜索过滤** - 支持成员搜索和过滤
- **批量选择** - 支持多选和批量操作

### 基础用法

```vue
<template>
  <div class="p-4">
    <div class="mb-4">
      <h3>单选用户</h3>
      <ComplexMemberPicker
        v-model:value="selectedUser"
        type="user"
        placeholder="请选择用户"
        @change="handleUserChange"
      />
    </div>
    
    <div class="mb-4">
      <h3>多选部门</h3>
      <ComplexMemberPicker
        v-model:value="selectedDepts"
        type="department"
        :multiple="true"
        placeholder="请选择部门"
      />
    </div>
    
    <div class="mb-4">
      <h3>混合选择</h3>
      <ComplexMemberPicker
        v-model:value="selectedMembers"
        :types="['user', 'department', 'role']"
        :multiple="true"
        :show-type-tabs="true"
        placeholder="请选择成员"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { ComplexMemberPicker } from '@/components/ComplexMemberPicker';

const selectedUser = ref('');
const selectedDepts = ref<string[]>([]);
const selectedMembers = ref<string[]>([]);

const handleUserChange = (value: string, option: any) => {
  console.log('选择的用户:', value, option);
};
</script>
```

### API 接口

| 属性 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| modelValue | `string \| string[]` | - | 绑定值 |
| type | `'user' \| 'department' \| 'role'` | `'user'` | 选择类型 |
| types | `string[]` | - | 支持的类型列表 |
| multiple | `boolean` | `false` | 是否多选 |
| showTypeTabs | `boolean` | `false` | 是否显示类型标签 |
| searchable | `boolean` | `true` | 是否可搜索 |
| treeMode | `boolean` | `true` | 是否树形模式 |

## DatePicker 日期选择器

### 组件特性

- **多种日期格式** - 支持日期、时间、日期时间等格式
- **快捷选择** - 内置常用日期快捷选择
- **范围选择** - 支持日期范围选择
- **自定义格式** - 支持自定义日期显示格式

### 基础用法

```vue
<template>
  <div class="p-4">
    <div class="mb-4">
      <h3>基础日期选择</h3>
      <DatePicker
        v-model:value="basicDate"
        type="date"
        placeholder="请选择日期"
      />
    </div>
    
    <div class="mb-4">
      <h3>日期时间选择</h3>
      <DatePicker
        v-model:value="datetimeValue"
        type="datetime"
        placeholder="请选择日期时间"
        :show-time="true"
      />
    </div>
    
    <div class="mb-4">
      <h3>日期范围选择</h3>
      <DatePicker
        v-model:value="dateRange"
        type="daterange"
        :shortcuts="dateShortcuts"
        placeholder="请选择日期范围"
      />
    </div>
    
    <div class="mb-4">
      <h3>带快捷选择</h3>
      <DatePicker
        v-model:value="shortcutDate"
        type="date"
        :shortcuts="singleDateShortcuts"
        placeholder="请选择日期"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { DatePicker } from '@/components/DatePicker';

const basicDate = ref(null);
const datetimeValue = ref(null);
const dateRange = ref<[number, number] | null>(null);
const shortcutDate = ref(null);

const dateShortcuts = [
  {
    label: '最近一周',
    value: () => {
      const end = new Date();
      const start = new Date();
      start.setTime(start.getTime() - 3600 * 1000 * 24 * 7);
      return [start.getTime(), end.getTime()];
    },
  },
  {
    label: '最近一个月',
    value: () => {
      const end = new Date();
      const start = new Date();
      start.setTime(start.getTime() - 3600 * 1000 * 24 * 30);
      return [start.getTime(), end.getTime()];
    },
  },
];

const singleDateShortcuts = [
  {
    label: '今天',
    value: () => new Date().getTime(),
  },
  {
    label: '昨天',
    value: () => {
      const date = new Date();
      date.setDate(date.getDate() - 1);
      return date.getTime();
    },
  },
];
</script>
```

## 组件开发规范

### 1. 组件结构

```
ComponentName/
├── index.ts              # 组件导出
├── index.vue             # 主组件
├── src/                  # 子组件
│   ├── SubComponent.vue
│   └── utils.ts
├── props.ts              # 属性定义
├── types.ts              # 类型定义
└── README.md             # 组件文档
```

### 2. 组件模板

```vue
<template>
  <div class="component-name" :class="cssClass">
    <slot name="header" />
    <div class="component-content">
      <slot />
    </div>
    <slot name="footer" />
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { componentProps } from './props';

defineOptions({
  name: 'ComponentName',
});

const props = defineProps(componentProps);
const emit = defineEmits<{
  change: [value: any];
  click: [event: MouseEvent];
}>();

const cssClass = computed(() => ({
  [`component-name--${props.size}`]: props.size,
  'component-name--disabled': props.disabled,
}));
</script>

<style scoped lang="less">
.component-name {
  &--small {
    font-size: 12px;
  }
  
  &--disabled {
    opacity: 0.6;
    pointer-events: none;
  }
}
</style>
```

### 3. 属性定义规范

```typescript
// props.ts
import type { ExtractPropTypes, PropType } from 'vue';

export const componentProps = {
  modelValue: {
    type: [String, Number, Array] as PropType<any>,
    default: undefined,
  },
  size: {
    type: String as PropType<'small' | 'medium' | 'large'>,
    default: 'medium',
  },
  disabled: {
    type: Boolean,
    default: false,
  },
  onChange: {
    type: Function as PropType<(value: any) => void>,
  },
} as const;

export type ComponentProps = ExtractPropTypes<typeof componentProps>;
```

---

这些组件文档涵盖了 HotGo 2.0 中的主要组件，为开发者提供了完整的使用指南和最佳实践。每个组件都经过精心设计，确保在企业级应用中的稳定性和可扩展性。






