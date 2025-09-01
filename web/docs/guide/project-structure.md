# 项目结构

本文档详细介绍 HotGo 2.0 前端项目的目录结构和文件组织方式。

## 总体结构

```
web/
├── build/                  # 构建配置文件
├── docs/                   # 项目文档
├── public/                 # 静态资源文件
├── src/                    # 源代码目录
├── types/                  # 全局类型定义
├── .env.example           # 环境变量示例
├── .eslintrc.js           # ESLint 配置
├── .gitignore             # Git 忽略文件
├── .prettierrc            # Prettier 配置
├── index.html             # HTML 入口文件
├── package.json           # 项目依赖配置
├── tsconfig.json          # TypeScript 配置
├── vite.config.ts         # Vite 配置文件
└── README.md              # 项目说明文档
```

## 核心目录详解

### `/src` 源代码目录

```
src/
├── api/                    # API 接口定义
├── assets/                 # 静态资源
├── components/             # 公共组件
├── directives/             # 自定义指令
├── enums/                  # 枚举类型定义
├── hooks/                  # 组合式函数
├── layout/                 # 布局组件
├── plugins/                # 插件配置
├── router/                 # 路由配置
├── settings/               # 项目配置
├── store/                  # 状态管理
├── styles/                 # 样式文件
├── utils/                  # 工具函数
├── views/                  # 页面组件
├── App.vue                 # 根组件
└── main.ts                 # 应用入口
```

## 详细目录说明

### 📡 API 接口层 (`/src/api`)

```
api/
├── addons/                 # 插件相关接口
│   └── hgexample/         # 示例插件接口
├── apply/                  # 申请相关接口
│   ├── attachment.ts      # 附件管理
│   ├── notice.ts          # 通知公告
│   └── provinces.ts       # 省市区数据
├── system/                 # 系统管理接口
│   ├── menu.ts            # 菜单管理
│   ├── role.ts            # 角色管理
│   └── user.ts            # 用户管理
├── org/                    # 组织架构接口
│   ├── dept.ts            # 部门管理
│   ├── post.ts            # 岗位管理
│   └── user.ts            # 员工管理
└── base/                   # 基础接口
    └── index.ts           # 通用接口
```

**特点**：
- 按业务模块划分
- 统一的接口定义规范
- 完整的 TypeScript 类型支持
- 支持请求/响应拦截器

### 🧩 组件库 (`/src/components`)

```
components/
├── Application/            # 应用容器组件
├── Table/                  # 表格组件
│   ├── index.ts           # 组件导出
│   └── src/               # 组件源码
│       ├── Table.vue      # 主表格组件
│       ├── components/    # 子组件
│       ├── hooks/         # 表格相关 hooks
│       ├── props.ts       # 属性定义
│       └── types/         # 类型定义
├── Form/                   # 表单组件
│   ├── index.ts           # 组件导出
│   └── src/               # 组件源码
│       ├── BasicForm.vue  # 基础表单
│       ├── hooks/         # 表单 hooks
│       ├── props.ts       # 属性定义
│       └── types/         # 类型定义
├── Modal/                  # 模态框组件
├── Upload/                 # 上传组件
├── Editor/                 # 富文本编辑器
├── FileChooser/           # 文件选择器
└── ...                    # 其他组件
```

**设计原则**：
- 高内聚低耦合
- 支持按需导入
- 完整的 TypeScript 支持
- 丰富的配置选项

### 🎣 Hooks 层 (`/src/hooks`)

```
hooks/
├── common/                 # 通用 hooks
│   ├── useBoolean.ts      # 布尔值管理
│   ├── useLoading.ts      # 加载状态管理
│   ├── useContext.ts      # 上下文管理
│   └── useSorter.ts       # 排序管理
├── web/                    # Web 相关 hooks
│   ├── usePermission.ts   # 权限管理
│   ├── usePage.ts         # 页面管理
│   └── useECharts.ts      # 图表管理
├── setting/                # 设置相关 hooks
│   ├── useDesignSetting.ts # 设计设置
│   └── useProjectSetting.ts # 项目设置
└── index.ts               # hooks 导出
```

**特色功能**：
- 业务逻辑复用
- 状态管理抽象
- 生命周期管理
- 副作用处理

### 🎨 布局系统 (`/src/layout`)

```
layout/
├── components/             # 布局子组件
│   ├── Header/            # 头部组件
│   │   ├── index.vue      # 主头部
│   │   ├── UserDropdown.vue # 用户下拉菜单
│   │   ├── Breadcrumb.vue # 面包屑导航
│   │   └── FullScreen.vue # 全屏切换
│   ├── Menu/              # 侧边菜单
│   │   ├── index.vue      # 主菜单
│   │   └── SubMenu.vue    # 子菜单
│   ├── TagsView/          # 标签页导航
│   │   └── index.vue      # 标签页组件
│   └── Main/              # 主内容区
│       └── index.vue      # 内容容器
├── index.vue              # 主布局组件
└── parentLayout.vue       # 父级布局
```

**布局特性**：
- 响应式设计
- 多布局模式
- 主题切换支持
- 移动端适配

### 🛣️ 路由系统 (`/src/router`)

```
router/
├── modules/               # 路由模块
│   ├── dashboard.ts       # 仪表板路由
│   ├── system.ts          # 系统管理路由
│   ├── org.ts             # 组织架构路由
│   └── ...                # 其他业务路由
├── base.ts                # 基础路由
├── constant.ts            # 路由常量
├── generator-routers.ts   # 动态路由生成
├── index.ts               # 路由配置
├── router-guards.ts       # 路由守卫
└── router-icons.ts        # 路由图标
```

**路由特性**：
- 动态路由加载
- 权限控制
- 路由守卫
- 面包屑导航自动生成

### 🗄️ 状态管理 (`/src/store`)

```
store/
├── modules/               # 状态模块
│   ├── user.ts           # 用户状态
│   ├── permission.ts     # 权限状态
│   ├── asyncRoute.ts     # 异步路由状态
│   ├── designSetting.ts  # 设计设置状态
│   ├── dict.ts           # 字典状态
│   └── tagsView.ts       # 标签页状态
├── index.ts              # Store 配置
├── mutation-types.ts     # 变更类型
└── types.ts              # 状态类型定义
```

**状态管理特性**：
- 模块化设计
- TypeScript 支持
- 持久化存储
- 开发工具支持

### 🔧 工具函数 (`/src/utils`)

```
utils/
├── http/                  # HTTP 请求相关
│   └── axios/            # Axios 封装
│       ├── Axios.ts      # 主要封装类
│       ├── axiosCancel.ts # 请求取消
│       ├── axiosTransform.ts # 请求转换
│       └── types.ts      # 类型定义
├── websocket/            # WebSocket 封装
│   ├── index.ts          # 主要封装
│   └── registerMessage.ts # 消息注册
├── is/                   # 类型判断工具
├── lib/                  # 第三方库封装
├── array.ts              # 数组工具
├── browser-type.ts       # 浏览器检测
├── charset.ts            # 字符集处理
└── ...                   # 其他工具函数
```

**工具特色**：
- 纯函数设计
- 高性能实现
- 完整单元测试
- 详细注释文档

### 📄 页面组件 (`/src/views`)

```
views/
├── dashboard/             # 仪表板页面
│   ├── console/          # 控制台
│   └── workplace/        # 工作台
├── system/               # 系统管理页面
│   ├── user/             # 用户管理
│   ├── role/             # 角色管理
│   ├── menu/             # 菜单管理
│   ├── dict/             # 字典管理
│   └── config/           # 系统配置
├── org/                  # 组织架构页面
│   ├── dept/             # 部门管理
│   ├── post/             # 岗位管理
│   └── user/             # 员工管理
├── monitor/              # 监控页面
│   ├── online/           # 在线用户
│   ├── serve-monitor/    # 服务监控
│   └── netconn/          # 网络连接
├── login/                # 登录页面
├── exception/            # 异常页面
│   ├── 403.vue          # 无权限
│   ├── 404.vue          # 页面不存在
│   └── 500.vue          # 服务器错误
└── ...                   # 其他业务页面
```

**页面组织原则**：
- 按功能模块划分
- 统一的页面结构
- 复用的页面组件
- 标准的路由配置

### 🎨 样式系统 (`/src/styles`)

```
styles/
├── common.less           # 通用样式
├── hotgo.less           # 项目主题样式
├── index.less           # 样式入口文件
├── tailwind.css         # Tailwind CSS
├── var.less             # 样式变量
├── transition/          # 过渡动画
│   ├── base.less        # 基础过渡
│   ├── fade.less        # 淡入淡出
│   ├── slide.less       # 滑动效果
│   └── zoom.less        # 缩放效果
└── mixins.less          # 样式混合
```

**样式特性**：
- Less 预处理器
- 原子化 CSS (Tailwind)
- 主题系统
- 响应式设计

### ⚙️ 配置文件 (`/src/settings`)

```
settings/
├── animateSetting.ts     # 动画配置
├── componentSetting.ts   # 组件配置
├── designSetting.ts      # 设计配置
└── projectSetting.ts     # 项目配置
```

### 📝 类型定义 (`/types`)

```
types/
├── config.d.ts           # 配置类型
├── global.d.ts           # 全局类型
├── index.d.ts            # 类型入口
├── module.d.ts           # 模块类型
└── vue-router.d.ts       # 路由类型扩展
```

## 文件命名规范

### 组件文件
- **组件文件名**：使用 PascalCase（如：`UserList.vue`）
- **组件目录**：使用 kebab-case（如：`user-list/`）
- **组件导出**：统一通过 `index.ts` 导出

### 页面文件
- **页面文件名**：使用 kebab-case（如：`user-list.vue`）
- **页面目录**：按功能模块划分
- **路由文件**：使用 kebab-case（如：`user-management.ts`）

### 工具文件
- **工具文件名**：使用 camelCase（如：`arrayUtils.ts`）
- **常量文件名**：使用 UPPER_SNAKE_CASE（如：`API_CONSTANTS.ts`）
- **类型文件名**：使用 camelCase + `.d.ts`（如：`userTypes.d.ts`）

## 导入路径规范

### 路径别名
```typescript
// 使用 @ 别名指向 src 目录
import { BasicTable } from '@/components/Table';
import { usePermission } from '@/hooks/web/usePermission';
import { ApiEnum } from '@/enums/apiEnum';

// 使用 /# 别名指向 types 目录
import type { UserInfo } from '/#/user';
```

### 导入顺序
```typescript
// 1. Vue 相关
import { ref, reactive, computed } from 'vue';
import { useRouter, useRoute } from 'vue-router';

// 2. 第三方库
import { NButton, NDataTable } from 'naive-ui';
import { format } from 'date-fns';

// 3. 项目内部 - 组件
import { BasicTable, BasicForm } from '@/components';

// 4. 项目内部 - hooks/utils
import { usePermission } from '@/hooks/web/usePermission';
import { formatDate } from '@/utils/dateUtil';

// 5. 项目内部 - API/store
import { getUserList } from '@/api/system/user';
import { useUserStore } from '@/store/modules/user';

// 6. 类型导入
import type { FormSchema } from '@/components/Form/types';
import type { UserInfo } from '/#/user';
```

## 最佳实践

### 1. 目录组织
- 功能相关的文件放在同一目录
- 公共功能提取到 `components` 或 `utils`
- 业务逻辑抽象到 `hooks`
- 类型定义集中管理

### 2. 文件结构
- 每个组件目录包含完整的功能实现
- 提供清晰的导出接口
- 包含详细的类型定义
- 添加必要的文档注释

### 3. 依赖管理
- 避免循环依赖
- 明确依赖关系
- 合理使用路径别名
- 统一导入顺序

---

下一步：[核心架构](../architecture/core.md)

