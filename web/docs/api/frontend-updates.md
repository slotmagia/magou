# 前端页面交互逻辑更新说明

本文档说明了为支持多租户API而对前端页面交互逻辑所做的更新。

## 更新概览

根据 `docs/api/api-documentation.md` 中定义的多租户API接口，我们对以下前端页面进行了更新：

### 1. 登录页面更新 (`src/views/login/login/form.vue`)

#### 主要变更：
- **添加租户编码字段**：在用户名输入框上方增加了租户编码输入字段
- **更新登录逻辑**：调用新的多租户登录API (`/api/auth/index.ts`)
- **验证码获取**：优先使用多租户验证码API，失败时降级到原有API
- **表单验证**：增加租户编码的必填验证

#### 新增字段：
```typescript
interface FormState {
  tenantCode: string; // 新增：租户编码
  username: string;
  pass: string;
  cid: string;
  code: string;
  password: string;
}
```

#### API调用变更：
```typescript
// 旧版本
await userStore.login(params);

// 新版本
await userStore.tenantLogin(params);
```

### 2. 用户状态管理更新 (`src/store/modules/user.ts`)

#### 主要变更：
- **新增多租户登录方法**：`tenantLogin()`
- **新增多租户登录处理**：`handleTenantLogin()`
- **扩展存储信息**：存储租户信息、权限和菜单ID

#### 新增方法：
```typescript
// 多租户登录
async tenantLogin(userInfo) {
  return await this.handleTenantLogin(tenantLoginApi(userInfo));
}

// 多租户登录处理
async handleTenantLogin(request: Promise<any>) {
  // 存储多租户相关信息
  storage.set('TENANT_INFO', { 
    tenantId: data.userInfo.tenantId, 
    tenantCode: data.userInfo.tenantCode 
  }, ex);
  storage.set('USER_PERMISSIONS', data.permissions, ex);
  storage.set('USER_MENU_IDS', data.menuIds, ex);
  // ...
}
```

### 3. 角色管理页面更新

#### 3.1 角色列表页面 (`src/views/permission/role/role.vue`)
- **API集成**：引入新的多租户角色管理API
- **数据适配**：兼容新API响应格式和旧API格式
- **错误处理**：增加错误处理和降级机制

#### 3.2 角色编辑页面 (`src/views/permission/role/editRole.vue`)
- **创建/编辑逻辑**：使用新的 `createRole()` 和 `updateRole()` API
- **参数映射**：将表单数据映射到新API的参数格式
- **降级处理**：API失败时自动降级到旧API

#### 3.3 菜单权限编辑 (`src/views/permission/role/editMenuAuth.vue`)
- **权限更新**：使用新的 `updateRoleMenus()` API
- **参数格式**：调整为新API要求的参数格式

### 4. 用户管理页面更新 (`src/views/org/user/list.vue`)

#### 主要变更：
- **API集成**：引入多租户用户管理API
- **数据加载**：优先使用 `getTenantUserList()` API
- **用户操作**：支持 `createTenantUser()` 和 `updateTenantUser()`
- **响应适配**：适配新API的响应数据格式

#### 新增功能：
- 多租户用户创建和编辑
- 密码Base64编码处理
- 租户级别的用户管理

### 5. 菜单管理页面更新 (`src/views/permission/menu/menu.vue`)

#### 主要变更：
- **菜单树获取**：使用新的 `getMenuTree()` API
- **响应适配**：兼容新旧API响应格式
- **降级机制**：API失败时自动使用旧API

### 6. 新增租户管理页面 (`src/views/tenant/index.vue`)

#### 全新功能页面：
- **租户列表**：显示所有租户信息和状态
- **租户创建**：创建新租户及其管理员账号
- **租户编辑**：编辑租户基本信息
- **租户详情**：查看租户详细信息和统计数据
- **状态管理**：启用/禁用租户
- **删除操作**：删除租户（系统租户受保护）

#### 核心功能：
```typescript
// 主要操作
- 获取租户列表：getTenantList()
- 创建租户：createTenant()
- 更新租户：updateTenant()
- 删除租户：deleteTenant()
- 获取详情：getTenantDetail()
- 状态切换：updateTenantStatus()
```

## API兼容性策略

为了确保平滑过渡，所有更新都采用了**优雅降级**策略：

1. **优先使用新API**：首先尝试调用新的多租户API
2. **失败时降级**：如果新API调用失败，自动降级到原有API
3. **错误日志**：记录API调用失败的信息以便调试
4. **用户体验**：确保用户在API切换过程中不会感知到异常

### 示例降级代码：
```typescript
try {
  // 优先使用新的多租户API
  const response = await getTenantUserList(params);
  return response.data;
} catch (error) {
  // 如果新API失败，降级到旧API
  console.warn('多租户API调用失败，降级到旧API:', error);
  return await List(params);
}
```

## 数据格式变更

### 登录响应格式
```typescript
// 旧格式
{
  code: 0,
  data: {
    token: string,
    userInfo: UserInfo
  }
}

// 新格式（多租户）
{
  code: 0,
  data: {
    accessToken: string,
    refreshToken: string,
    userInfo: TenantUserInfo,
    permissions: string[],
    menuIds: number[]
  }
}
```

### 列表响应格式
```typescript
// 旧格式
{
  list: T[],
  total: number
}

// 新格式
{
  code: 0,
  data: {
    list: T[],
    total: number,
    page: number,
    pageSize: number
  }
}
```

## 类型定义更新

新增了完整的TypeScript类型定义：

### 租户相关类型 (`src/api/tenant/index.ts`)
- `TenantInfo` - 租户信息
- `TenantConfig` - 租户配置
- `TenantStats` - 租户统计
- `CreateTenantParams` - 创建租户参数
- `UpdateTenantParams` - 更新租户参数

### 认证相关类型 (`src/api/auth/index.ts`)
- `UserInfo` - 多租户用户信息
- `LoginParams` - 登录参数
- `LoginResult` - 登录结果
- `CaptchaResult` - 验证码结果

### 角色管理类型 (`src/api/system/role.ts`)
- `RoleInfo` - 角色信息
- `RoleMenuPermission` - 角色菜单权限
- `CreateRoleParams` - 创建角色参数
- `UpdateRoleParams` - 更新角色参数

## 使用说明

### 1. 多租户登录
1. 输入租户编码（必填）
2. 输入用户名和密码
3. 输入验证码
4. 点击登录

### 2. 租户管理
1. 访问租户管理页面
2. 可以查看、创建、编辑、删除租户
3. 支持租户状态管理（启用/禁用）
4. 查看租户详细统计信息

### 3. 角色和用户管理
- 现有功能保持不变
- 增加了多租户数据隔离
- 支持租户级别的权限管理

## 注意事项

1. **系统租户保护**：ID为1的系统租户受到特殊保护，不能编辑、禁用或删除
2. **密码安全**：所有密码都使用Base64编码传输
3. **数据隔离**：不同租户的数据完全隔离
4. **权限继承**：租户管理员拥有其租户内的所有权限
5. **API兼容**：新旧API完全兼容，支持平滑迁移

## 测试建议

1. **登录测试**：测试不同租户的用户登录
2. **权限测试**：验证租户间数据隔离
3. **API降级测试**：模拟API失败场景测试降级机制
4. **界面测试**：确保所有新增界面正常显示和操作
5. **兼容性测试**：验证与现有功能的兼容性

## 未来扩展

当前更新为多租户功能奠定了基础，未来可以扩展：

1. **租户主题定制**：允许租户自定义界面主题
2. **租户功能限制**：根据租户套餐限制功能使用
3. **租户监控面板**：实时监控租户使用情况
4. **租户计费系统**：根据使用情况计费
5. **子域名自动识别**：根据访问域名自动识别租户

通过以上更新，前端系统现在完全支持多租户架构，为企业级SaaS应用提供了坚实的基础。
