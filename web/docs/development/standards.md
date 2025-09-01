# 开发规范

本文档定义了 HotGo 2.0 项目的开发规范和最佳实践，确保代码质量和团队协作效率。

## 代码规范

### TypeScript 规范

#### 1. 类型定义

```typescript
// ✅ 推荐：使用 interface 定义对象类型
interface UserInfo {
  id: number;
  username: string;
  email: string;
  status: 1 | 0; // 使用字面量类型
}

// ✅ 推荐：使用 type 定义联合类型
type ButtonType = 'primary' | 'secondary' | 'danger';
type Status = 'loading' | 'success' | 'error';

// ❌ 避免：使用 any 类型
const userData: any = {};

// ✅ 推荐：使用具体类型
const userData: UserInfo = {
  id: 1,
  username: 'admin',
  email: 'admin@example.com',
  status: 1,
};
```

#### 2. 函数定义

```typescript
// ✅ 推荐：明确参数和返回值类型
async function getUserList(
  params: UserListParams
): Promise<PageResponse<UserInfo>> {
  return await userApi.getList(params);
}

// ✅ 推荐：使用可选参数和默认值
function createUser(
  userInfo: CreateUserRequest,
  options: CreateUserOptions = {}
): Promise<UserInfo> {
  const { validateEmail = true, sendWelcome = false } = options;
  // 实现逻辑
}

// ✅ 推荐：使用泛型提高复用性
function createApiResponse<T>(data: T, message = '操作成功'): ApiResponse<T> {
  return {
    code: 200,
    message,
    data,
  };
}
```

#### 3. 类和接口

```typescript
// ✅ 推荐：使用访问修饰符
class UserService {
  private readonly apiClient: ApiClient;
  
  constructor(apiClient: ApiClient) {
    this.apiClient = apiClient;
  }
  
  public async getUser(id: number): Promise<UserInfo> {
    return this.apiClient.get(`/users/${id}`);
  }
  
  private validateUserId(id: number): boolean {
    return id > 0;
  }
}

// ✅ 推荐：继承和实现
interface UserRepository {
  findById(id: number): Promise<UserInfo | null>;
  create(user: CreateUserRequest): Promise<UserInfo>;
}

class HttpUserRepository implements UserRepository {
  async findById(id: number): Promise<UserInfo | null> {
    // 实现
  }
  
  async create(user: CreateUserRequest): Promise<UserInfo> {
    // 实现
  }
}
```

### Vue 组件规范

#### 1. 组件结构

```vue
<template>
  <!-- 模板内容 -->
</template>

<script setup lang="ts">
// 1. 导入依赖
import { ref, reactive, computed, watch, onMounted } from 'vue';
import { NButton, NInput } from 'naive-ui';

// 2. 类型定义
interface Props {
  modelValue: string;
  disabled?: boolean;
}

interface Emits {
  (e: 'update:modelValue', value: string): void;
  (e: 'change', value: string): void;
}

// 3. Props 和 Emits
const props = withDefaults(defineProps<Props>(), {
  disabled: false,
});

const emits = defineEmits<Emits>();

// 4. 响应式数据
const inputValue = ref('');
const formData = reactive({
  username: '',
  email: '',
});

// 5. 计算属性
const isValid = computed(() => {
  return formData.username.length > 0 && formData.email.length > 0;
});

// 6. 方法
const handleSubmit = () => {
  if (!isValid.value) return;
  // 提交逻辑
};

// 7. 监听器
watch(
  () => props.modelValue,
  (newValue) => {
    inputValue.value = newValue;
  },
  { immediate: true }
);

// 8. 生命周期
onMounted(() => {
  // 初始化逻辑
});
</script>

<style scoped lang="less">
/* 样式内容 */
</style>
```

#### 2. 组件命名

```typescript
// ✅ 推荐：使用 PascalCase
export default defineComponent({
  name: 'UserManagementTable',
});

// ✅ 推荐：组件文件名
UserManagementTable.vue
UserEditModal.vue
DashboardChart.vue

// ❌ 避免：小写或短横线命名
userTable.vue
user-modal.vue
chart.vue
```

#### 3. Props 定义

```typescript
// ✅ 推荐：完整的 Props 定义
interface UserTableProps {
  /** 数据源 */
  dataSource: UserInfo[];
  /** 是否显示分页 */
  showPagination?: boolean;
  /** 每页数量 */
  pageSize?: number;
  /** 加载状态 */
  loading?: boolean;
  /** 选择模式 */
  selectionMode?: 'single' | 'multiple' | 'none';
}

const props = withDefaults(defineProps<UserTableProps>(), {
  showPagination: true,
  pageSize: 10,
  loading: false,
  selectionMode: 'multiple',
});
```

### CSS/Less 规范

#### 1. 类名命名

```less
// ✅ 推荐：使用 BEM 命名法
.user-table {
  &__header {
    display: flex;
    justify-content: space-between;
    
    &--sticky {
      position: sticky;
      top: 0;
    }
  }
  
  &__content {
    padding: 16px;
  }
  
  &__row {
    &:hover {
      background-color: #f5f5f5;
    }
    
    &--selected {
      background-color: #e6f7ff;
    }
  }
}

// ✅ 推荐：使用语义化类名
.primary-button { }
.success-message { }
.error-text { }

// ❌ 避免：使用表现形式类名
.red-text { }
.big-font { }
.center-div { }
```

#### 2. 样式组织

```less
// 1. 变量定义
@primary-color: #1890ff;
@success-color: #52c41a;
@error-color: #ff4d4f;
@border-radius: 6px;

// 2. 混入定义
.flex-center {
  display: flex;
  justify-content: center;
  align-items: center;
}

.text-ellipsis {
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
}

// 3. 组件样式
.component-name {
  // 布局属性
  display: flex;
  position: relative;
  
  // 盒模型属性
  width: 100%;
  height: auto;
  margin: 0;
  padding: 16px;
  
  // 边框和背景
  border: 1px solid #d9d9d9;
  border-radius: @border-radius;
  background-color: #fff;
  
  // 字体属性
  font-size: 14px;
  font-weight: normal;
  color: #333;
  
  // 其他属性
  transition: all 0.3s ease;
  
  // 伪类
  &:hover {
    border-color: @primary-color;
  }
  
  // 嵌套选择器
  .child-element {
    // 子元素样式
  }
}
```

## 文件组织规范

### 目录结构

```
src/
├── components/           # 公共组件
│   ├── BasicTable/      # 组件目录
│   │   ├── index.ts     # 导出文件
│   │   ├── src/         # 源码目录
│   │   │   ├── Table.vue
│   │   │   ├── props.ts
│   │   │   └── types.ts
│   │   └── README.md    # 组件文档
│   └── index.ts         # 统一导出
├── views/               # 页面组件
│   ├── user/           # 功能模块
│   │   ├── list/       # 页面目录
│   │   │   ├── index.vue
│   │   │   ├── components/
│   │   │   └── hooks/
│   │   └── detail/
│   └── dashboard/
└── api/                # 接口定义
    ├── types/          # 类型定义
    ├── user.ts         # 用户接口
    └── index.ts        # 统一导出
```

### 文件命名

```
// ✅ 推荐的文件命名
UserManagement.vue       # 页面组件：PascalCase
BasicTable.vue          # 公共组件：PascalCase
userApi.ts              # 接口文件：camelCase
userTypes.ts            # 类型文件：camelCase
useUserManagement.ts    # Hook文件：camelCase with use前缀

// ❌ 避免的文件命名
user_management.vue     # 下划线命名
user-api.ts            # 短横线命名
User.vue               # 过于简单的命名
```

## Git 提交规范

### 提交信息格式

```
<type>(<scope>): <subject>

<body>

<footer>
```

### 类型说明

| 类型 | 说明 |
|------|------|
| feat | 新功能 |
| fix | 修复bug |
| docs | 文档变更 |
| style | 代码格式调整（不影响代码运行） |
| refactor | 重构（既不是新增功能，也不是修复bug） |
| perf | 性能优化 |
| test | 增加测试 |
| chore | 构建过程或辅助工具的变动 |
| revert | 回滚 |

### 提交示例

```bash
# 新功能
feat(user): 添加用户管理功能

增加用户列表展示、创建、编辑、删除功能
- 实现用户列表页面
- 添加用户编辑模态框
- 集成用户API接口

Closes #123

# 修复bug
fix(table): 修复表格分页显示异常

修复当数据总数为0时分页组件显示错误的问题

# 文档更新
docs: 更新API文档

添加用户管理接口文档说明

# 重构
refactor(auth): 重构权限验证逻辑

将权限验证逻辑从组件中抽离到hooks中，提高代码复用性
```

## 代码注释规范

### JSDoc 注释

```typescript
/**
 * 用户管理服务类
 * @description 提供用户相关的业务操作方法
 * @author 开发者姓名
 * @since 1.0.0
 */
class UserService {
  /**
   * 获取用户列表
   * @param params - 查询参数
   * @param params.page - 页码
   * @param params.pageSize - 每页数量
   * @param params.keyword - 搜索关键词
   * @returns 用户列表数据
   * @throws {BusinessError} 当参数无效时抛出业务错误
   * @example
   * ```typescript
   * const users = await userService.getUsers({
   *   page: 1,
   *   pageSize: 10,
   *   keyword: 'admin'
   * });
   * ```
   */
  async getUsers(params: UserListParams): Promise<PageResponse<UserInfo>> {
    // 实现逻辑
  }
}
```

### 行内注释

```typescript
// ✅ 推荐：解释业务逻辑
function calculateUserScore(user: UserInfo): number {
  // 基础分数：根据用户等级计算
  let score = user.level * 10;
  
  // 活跃度加分：最近30天登录次数
  score += Math.min(user.loginCount, 50);
  
  // 完善度加分：个人信息完整度
  if (user.avatar && user.nickname && user.email) {
    score += 20;
  }
  
  return score;
}

// ✅ 推荐：标注重要信息
const API_TIMEOUT = 30000; // 30秒超时，避免长时间等待

// TODO: 需要优化性能，考虑使用虚拟滚动
const renderLargeList = (items: any[]) => {
  // 渲染逻辑
};

// FIXME: 临时解决方案，需要重构
const handleLegacyData = (data: any) => {
  // 处理逻辑
};
```

## 错误处理规范

### 错误分类

```typescript
// 业务错误
class BusinessError extends Error {
  constructor(
    public code: number,
    message: string,
    public details?: any
  ) {
    super(message);
    this.name = 'BusinessError';
  }
}

// 验证错误
class ValidationError extends BusinessError {
  constructor(message: string, public field?: string) {
    super(400, message);
    this.name = 'ValidationError';
  }
}

// 网络错误
class NetworkError extends Error {
  constructor(message: string, public status?: number) {
    super(message);
    this.name = 'NetworkError';
  }
}
```

### 错误处理

```typescript
// ✅ 推荐：具体的错误处理
async function handleUserSubmit(userData: CreateUserRequest) {
  try {
    // 前置验证
    validateUserData(userData);
    
    // 创建用户
    const user = await createUser(userData);
    
    // 成功处理
    window.$message.success('用户创建成功');
    return user;
    
  } catch (error) {
    if (error instanceof ValidationError) {
      // 验证错误：显示字段错误
      formRef.value?.setFieldError(error.field, error.message);
    } else if (error instanceof BusinessError) {
      // 业务错误：显示错误消息
      window.$message.error(error.message);
    } else if (error instanceof NetworkError) {
      // 网络错误：提示网络问题
      window.$message.error('网络连接失败，请稍后重试');
    } else {
      // 未知错误：记录日志并显示通用错误
      console.error('Unexpected error:', error);
      window.$message.error('操作失败，请稍后重试');
    }
    
    throw error; // 重新抛出错误供上级处理
  }
}
```

## 性能优化规范

### 组件优化

```vue
<script setup lang="ts">
// ✅ 推荐：使用 shallowRef 优化大对象
const largeData = shallowRef<any[]>([]);

// ✅ 推荐：使用 markRaw 标记不需要响应式的对象
const chartInstance = markRaw(new Chart());

// ✅ 推荐：使用 computed 缓存计算结果
const expensiveComputed = computed(() => {
  return heavyCalculation(props.data);
});

// ✅ 推荐：使用 watchEffect 代替多个 watch
watchEffect(() => {
  // 自动收集依赖
  updateChart(props.data, props.options);
});

// ✅ 推荐：条件渲染优化
const shouldRenderChart = computed(() => {
  return props.data.length > 0 && props.visible;
});
</script>

<template>
  <!-- 使用 v-show 替代频繁切换的 v-if -->
  <div v-show="shouldRenderChart" class="chart-container">
    <!-- 图表内容 -->
  </div>
  
  <!-- 使用 key 优化列表渲染 -->
  <div v-for="item in list" :key="item.id" class="list-item">
    {{ item.name }}
  </div>
</template>
```

### 代码分割

```typescript
// ✅ 推荐：路由级别代码分割
const routes = [
  {
    path: '/user',
    component: () => import('@/views/user/index.vue'),
  },
  {
    path: '/dashboard',
    component: () => import('@/views/dashboard/index.vue'),
  },
];

// ✅ 推荐：组件级别代码分割
const HeavyComponent = defineAsyncComponent({
  loader: () => import('./components/HeavyComponent.vue'),
  loadingComponent: LoadingComponent,
  errorComponent: ErrorComponent,
  delay: 200,
  timeout: 3000,
});
```

## 测试规范

### 单元测试

```typescript
// userService.test.ts
import { describe, it, expect, vi } from 'vitest';
import { UserService } from '@/services/UserService';

describe('UserService', () => {
  it('should get user list successfully', async () => {
    // 准备测试数据
    const mockUsers = [
      { id: 1, username: 'user1' },
      { id: 2, username: 'user2' },
    ];
    
    // 模拟API调用
    const mockApi = vi.fn().mockResolvedValue({
      items: mockUsers,
      total: 2,
    });
    
    const userService = new UserService(mockApi);
    
    // 执行测试
    const result = await userService.getUsers({ page: 1, pageSize: 10 });
    
    // 验证结果
    expect(result.items).toEqual(mockUsers);
    expect(result.total).toBe(2);
    expect(mockApi).toHaveBeenCalledWith({ page: 1, pageSize: 10 });
  });
});
```

### 组件测试

```typescript
// UserTable.test.ts
import { mount } from '@vue/test-utils';
import { describe, it, expect } from 'vitest';
import UserTable from '@/components/UserTable.vue';

describe('UserTable', () => {
  it('should render user list correctly', () => {
    const users = [
      { id: 1, username: 'user1', email: 'user1@test.com' },
      { id: 2, username: 'user2', email: 'user2@test.com' },
    ];
    
    const wrapper = mount(UserTable, {
      props: { users },
    });
    
    // 验证渲染结果
    expect(wrapper.findAll('.user-row')).toHaveLength(2);
    expect(wrapper.text()).toContain('user1');
    expect(wrapper.text()).toContain('user1@test.com');
  });
  
  it('should emit edit event when edit button clicked', async () => {
    const users = [{ id: 1, username: 'user1', email: 'user1@test.com' }];
    
    const wrapper = mount(UserTable, {
      props: { users },
    });
    
    // 点击编辑按钮
    await wrapper.find('.edit-button').trigger('click');
    
    // 验证事件发射
    expect(wrapper.emitted('edit')).toBeTruthy();
    expect(wrapper.emitted('edit')[0]).toEqual([users[0]]);
  });
});
```

---

下一步：[部署指南](../deployment/deployment.md)






