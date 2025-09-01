# 多租户API接口使用示例

本文档提供了多租户管理系统API接口的使用示例和最佳实践。

## 概述

所有API接口都已经集成到前端项目中，位于 `src/api/` 目录下：

- `src/api/tenant/index.ts` - 租户管理接口
- `src/api/auth/index.ts` - 多租户登录认证接口
- `src/api/system/role.ts` - 多租户角色管理接口（已更新）
- `src/api/system/menu.ts` - 多租户菜单管理接口（已更新）
- `src/api/system/user.ts` - 多租户用户管理接口（已更新）

## 1. 租户管理使用示例

### 1.1 获取租户列表

```typescript
import { getTenantList, TenantListParams } from '@/api/tenant';

// 组件中使用
const fetchTenantList = async () => {
  try {
    const params: TenantListParams = {
      page: 1,
      pageSize: 20,
      name: '测试', // 可选：按名称过滤
      status: 1    // 可选：按状态过滤
    };
    
    const response = await getTenantList(params);
    
    if (response.code === 0) {
      console.log('租户列表:', response.data.list);
      console.log('总数:', response.data.total);
    }
  } catch (error) {
    console.error('获取租户列表失败:', error);
  }
};
```

### 1.2 创建租户

```typescript
import { createTenant, CreateTenantParams } from '@/api/tenant';

const createNewTenant = async () => {
  try {
    const params: CreateTenantParams = {
      name: '新企业租户',
      code: 'new_company',
      domain: 'newcompany.example.com',
      maxUsers: 200,
      storageLimit: 2147483648, // 2GB
      expireAt: '2025-12-31T23:59:59+08:00',
      adminName: 'company_admin',
      adminEmail: 'admin@newcompany.com',
      adminPassword: btoa('123456'), // Base64编码
      remark: '新企业客户'
    };
    
    const response = await createTenant(params);
    
    if (response.code === 0) {
      console.log('租户创建成功:', response.data);
    }
  } catch (error) {
    console.error('创建租户失败:', error);
  }
};
```

### 1.3 更新租户配置

```typescript
import { updateTenantConfig, TenantConfig } from '@/api/tenant';

const updateConfig = async (tenantId: number) => {
  try {
    const config: TenantConfig = {
      features: {
        advancedReports: true,
        apiAccess: true,
        customBranding: true
      },
      limitations: {
        maxApiCalls: 20000,
        maxStorage: 2147483648
      },
      settings: {
        theme: 'light',
        language: 'zh-CN',
        timezone: 'Asia/Shanghai'
      }
    };
    
    const response = await updateTenantConfig({
      id: tenantId,
      config
    });
    
    if (response.code === 0) {
      console.log('配置更新成功');
    }
  } catch (error) {
    console.error('更新配置失败:', error);
  }
};
```

## 2. 多租户认证使用示例

### 2.1 用户登录

```typescript
import { login, getCaptcha, LoginParams } from '@/api/auth';

// 获取验证码
const fetchCaptcha = async () => {
  try {
    const response = await getCaptcha();
    if (response.code === 0) {
      return {
        captchaId: response.data.captchaId,
        captchaImage: response.data.captchaImage
      };
    }
  } catch (error) {
    console.error('获取验证码失败:', error);
  }
};

// 用户登录
const userLogin = async (captchaData: any) => {
  try {
    const params: LoginParams = {
      tenantCode: 'test',
      username: 'test_user',
      password: btoa('123456'), // Base64编码
      captchaId: captchaData.captchaId,
      captcha: '1234',
      rememberMe: false
    };
    
    const response = await login(params);
    
    if (response.code === 0) {
      // 保存Token和用户信息
      const { accessToken, userInfo, permissions, menuIds } = response.data;
      
      localStorage.setItem('token', accessToken);
      localStorage.setItem('userInfo', JSON.stringify(userInfo));
      localStorage.setItem('permissions', JSON.stringify(permissions));
      
      console.log('登录成功:', userInfo);
    }
  } catch (error) {
    console.error('登录失败:', error);
  }
};
```

### 2.2 Token验证和刷新

```typescript
import { validateToken, refreshToken } from '@/api/auth';

// 验证Token有效性
const checkToken = async () => {
  try {
    const response = await validateToken();
    if (response.code === 0) {
      return response.data.valid;
    }
  } catch (error) {
    console.error('Token验证失败:', error);
    return false;
  }
};

// 刷新Token
const refreshUserToken = async () => {
  try {
    const refreshTokenValue = localStorage.getItem('refreshToken');
    if (!refreshTokenValue) return false;
    
    const response = await refreshToken(refreshTokenValue);
    
    if (response.code === 0) {
      localStorage.setItem('token', response.data.accessToken);
      localStorage.setItem('refreshToken', response.data.refreshToken);
      return true;
    }
  } catch (error) {
    console.error('刷新Token失败:', error);
    return false;
  }
};
```

## 3. 多租户角色管理使用示例

### 3.1 获取角色列表

```typescript
import { getRoleList, RoleListParams } from '@/api/system/role';

const fetchRoles = async () => {
  try {
    const params: RoleListParams = {
      page: 1,
      pageSize: 20,
      name: '管理员', // 可选：按名称过滤
      status: 1      // 可选：按状态过滤
    };
    
    const response = await getRoleList(params);
    
    if (response.code === 0) {
      console.log('角色列表:', response.data.list);
    }
  } catch (error) {
    console.error('获取角色列表失败:', error);
  }
};
```

### 3.2 创建角色

```typescript
import { createRole, CreateRoleParams } from '@/api/system/role';

const createNewRole = async () => {
  try {
    const params: CreateRoleParams = {
      name: '项目经理',
      code: 'project_manager',
      description: '项目经理角色，负责项目管理',
      status: 1,
      sort: 5,
      dataScope: 2,
      remark: '项目经理角色'
    };
    
    const response = await createRole(params);
    
    if (response.code === 0) {
      console.log('角色创建成功:', response.data);
    }
  } catch (error) {
    console.error('创建角色失败:', error);
  }
};
```

### 3.3 获取和更新角色菜单权限

```typescript
import { getRoleMenus, updateRoleMenus } from '@/api/system/role';

// 获取角色菜单权限
const fetchRoleMenus = async (roleId: number) => {
  try {
    const response = await getRoleMenus({ roleId });
    
    if (response.code === 0) {
      console.log('角色菜单权限:', response.data);
      return response.data.menuIds;
    }
  } catch (error) {
    console.error('获取角色菜单权限失败:', error);
  }
};

// 更新角色菜单权限
const updateRoleMenuPermissions = async (roleId: number, menuIds: number[]) => {
  try {
    const response = await updateRoleMenus({ roleId, menuIds });
    
    if (response.code === 0) {
      console.log('角色菜单权限更新成功');
    }
  } catch (error) {
    console.error('更新角色菜单权限失败:', error);
  }
};
```

## 4. 多租户菜单管理使用示例

### 4.1 获取菜单树

```typescript
import { getMenuTree, getMenuList } from '@/api/system/menu';

// 获取菜单树形结构
const fetchMenuTree = async () => {
  try {
    const response = await getMenuTree({ status: 1 });
    
    if (response.code === 0) {
      console.log('菜单树:', response.data);
      return response.data;
    }
  } catch (error) {
    console.error('获取菜单树失败:', error);
  }
};

// 获取用户可访问的菜单
const fetchUserMenus = async () => {
  try {
    const response = await getUserMenus();
    
    if (response.code === 0) {
      console.log('用户菜单:', response.data);
      return response.data;
    }
  } catch (error) {
    console.error('获取用户菜单失败:', error);
  }
};
```

### 4.2 创建菜单

```typescript
import { createMenu, CreateMenuParams } from '@/api/system/menu';

const createNewMenu = async () => {
  try {
    const params: CreateMenuParams = {
      parentId: 1,
      title: '新功能',
      name: 'NewFeature',
      path: '/new-feature',
      component: 'new-feature/index',
      icon: 'feature',
      type: 2, // 菜单类型
      sort: 10,
      status: 1,
      visible: 1,
      permission: 'new:feature:list',
      redirect: '',
      alwaysShow: 0,
      breadcrumb: 1,
      activeMenu: '',
      remark: '新功能菜单'
    };
    
    const response = await createMenu(params);
    
    if (response.code === 0) {
      console.log('菜单创建成功:', response.data);
    }
  } catch (error) {
    console.error('创建菜单失败:', error);
  }
};
```

## 5. 多租户用户管理使用示例

### 5.1 获取租户用户列表

```typescript
import { getTenantUserList, TenantUserListParams } from '@/api/system/user';

const fetchTenantUsers = async () => {
  try {
    const params: TenantUserListParams = {
      pageNumber: 1,
      pageSize: 20,
      username: 'test', // 可选：按用户名过滤
      status: 1         // 可选：按状态过滤
    };
    
    const response = await getTenantUserList(params);
    
    if (response.code === 0) {
      console.log('租户用户列表:', response.data.list);
    }
  } catch (error) {
    console.error('获取租户用户列表失败:', error);
  }
};
```

### 5.2 创建租户用户

```typescript
import { createTenantUser, CreateTenantUserParams } from '@/api/system/user';

const createNewTenantUser = async () => {
  try {
    const params: CreateTenantUserParams = {
      username: 'new_user',
      email: 'newuser@test.com',
      realName: '新用户',
      nickname: '小新',
      password: btoa('123456'), // Base64编码
      gender: 1,
      birthday: '1990-01-01',
      deptId: 1,
      position: '开发工程师',
      status: 1,
      remark: '新创建的用户'
    };
    
    const response = await createTenantUser(params);
    
    if (response.code === 0) {
      console.log('用户创建成功:', response.data);
    }
  } catch (error) {
    console.error('创建用户失败:', error);
  }
};
```

## 6. 错误处理最佳实践

### 6.1 统一错误处理

```typescript
// utils/errorHandler.ts
export const handleApiError = (error: any) => {
  if (error.response) {
    const { code, message } = error.response.data;
    
    switch (code) {
      case 401:
        // 未授权，跳转到登录页
        window.location.href = '/login';
        break;
      case 403:
        // 权限不足
        console.error('权限不足:', message);
        break;
      case 404:
        // 资源不存在
        console.error('资源不存在:', message);
        break;
      default:
        console.error('API错误:', message);
    }
  } else {
    console.error('网络错误:', error.message);
  }
};

// 在组件中使用
import { handleApiError } from '@/utils/errorHandler';

try {
  const response = await getTenantList();
  // 处理成功响应
} catch (error) {
  handleApiError(error);
}
```

### 6.2 租户切换处理

```typescript
// utils/tenantManager.ts
export class TenantManager {
  private static currentTenantId: number | null = null;
  
  static setCurrentTenant(tenantId: number) {
    this.currentTenantId = tenantId;
    localStorage.setItem('currentTenantId', tenantId.toString());
  }
  
  static getCurrentTenant(): number | null {
    if (!this.currentTenantId) {
      const stored = localStorage.getItem('currentTenantId');
      this.currentTenantId = stored ? parseInt(stored) : null;
    }
    return this.currentTenantId;
  }
  
  static clearCurrentTenant() {
    this.currentTenantId = null;
    localStorage.removeItem('currentTenantId');
  }
}

// 在HTTP拦截器中自动添加租户ID
import axios from 'axios';
import { TenantManager } from '@/utils/tenantManager';

axios.interceptors.request.use(config => {
  const tenantId = TenantManager.getCurrentTenant();
  if (tenantId) {
    config.headers['X-Tenant-Id'] = tenantId;
  }
  return config;
});
```

## 7. TypeScript 类型安全

所有API接口都包含完整的TypeScript类型定义，确保类型安全：

```typescript
// 类型推断示例
const response = await getTenantList(); // 自动推断为 BasicResponseModel<{list: TenantInfo[], total: number, ...}>

// 编译时类型检查
const params: CreateTenantParams = {
  name: '测试',
  code: 'test',
  // maxUsers: '100', // ❌ TypeScript错误：类型不匹配
  maxUsers: 100,     // ✅ 正确
  // ... 其他必填字段
};
```

## 8. 总结

本文档提供了多租户管理系统API接口的完整使用示例。所有接口都已经：

1. ✅ 完整对接API文档中的所有接口
2. ✅ 提供完整的TypeScript类型定义
3. ✅ 支持多租户数据隔离
4. ✅ 包含错误处理机制
5. ✅ 遵循项目现有的代码风格

使用这些API接口时，请确保：

- 正确处理租户ID的传递
- 妥善处理登录状态和Token刷新
- 实现适当的错误处理机制
- 遵循最佳的安全实践

如有疑问，请参考具体的API接口文档或联系开发团队。
