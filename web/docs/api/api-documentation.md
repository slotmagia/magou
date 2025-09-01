# 多租户管理系统 API 文档

## 概述

本文档描述了多租户管理系统的所有API接口，包括租户管理、用户管理、角色管理和菜单管理等功能。

### 基础信息

- **Base URL**: `http://localhost:8000/api`
- **认证方式**: Bearer Token
- **Content-Type**: `application/json`
- **租户识别**: 通过 `X-Tenant-Id` Header 或域名识别

### 通用响应格式

```json
{
  "code": 0,
  "message": "success",
  "data": {},
  "timestamp": 1703123456789,
  "traceId": "abc123def456"
}
```

### 错误响应格式

```json
{
  "code": 400,
  "message": "参数错误",
  "error": "详细错误信息",
  "timestamp": 1703123456789,
  "traceId": "abc123def456"
}
```

---

## 1. 租户管理

### 1.1 获取租户列表

**接口地址**: `GET /tenant/list`

**请求头**:
```
Authorization: Bearer <token>
```

**请求参数**:
```json
{
  "page": 1,
  "pageSize": 20,
  "name": "测试租户",
  "code": "test",
  "domain": "test.example.com",
  "status": 1
}
```

**参数说明**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| page | int | 否 | 页码，默认1 |
| pageSize | int | 否 | 每页数量，默认20，最大100 |
| name | string | 否 | 租户名称，支持模糊查询 |
| code | string | 否 | 租户编码，支持模糊查询 |
| domain | string | 否 | 租户域名，支持模糊查询 |
| status | int | 否 | 状态：0=全部 1=正常 2=锁定 3=禁用 |

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "name": "系统租户",
        "code": "system",
        "domain": "admin.example.com",
        "status": 1,
        "statusName": "正常",
        "maxUsers": 999999,
        "storageLimit": 10737418240,
        "expireAt": null,
        "adminUserId": 1,
        "adminName": "admin",
        "config": {
          "features": {
            "advancedReports": true,
            "apiAccess": true,
            "customBranding": true
          },
          "limitations": {
            "maxApiCalls": 100000,
            "maxStorage": 10737418240
          },
          "settings": {
            "theme": "light",
            "language": "zh-CN",
            "timezone": "Asia/Shanghai"
          }
        },
        "remark": "系统默认租户，不可删除",
        "createdAt": "2024-01-01T00:00:00+08:00",
        "updatedAt": "2024-01-01T00:00:00+08:00"
      },
      {
        "id": 2,
        "name": "测试租户",
        "code": "test",
        "domain": "test.example.com",
        "status": 1,
        "statusName": "正常",
        "maxUsers": 100,
        "storageLimit": 1073741824,
        "expireAt": "2024-12-31T23:59:59+08:00",
        "adminUserId": 2,
        "adminName": "test_admin",
        "config": {
          "features": {
            "advancedReports": false,
            "apiAccess": false,
            "customBranding": true
          },
          "limitations": {
            "maxApiCalls": 10000,
            "maxStorage": 1073741824
          },
          "settings": {
            "theme": "dark",
            "language": "zh-CN",
            "timezone": "Asia/Shanghai"
          }
        },
        "remark": "测试环境租户",
        "createdAt": "2024-01-02T00:00:00+08:00",
        "updatedAt": "2024-01-02T00:00:00+08:00"
      },
      {
        "id": 3,
        "name": "企业租户A",
        "code": "company_a",
        "domain": "companya.example.com",
        "status": 1,
        "statusName": "正常",
        "maxUsers": 500,
        "storageLimit": 5368709120,
        "expireAt": "2025-01-01T23:59:59+08:00",
        "adminUserId": 3,
        "adminName": "company_admin",
        "config": {
          "features": {
            "advancedReports": true,
            "apiAccess": true,
            "customBranding": true
          },
          "limitations": {
            "maxApiCalls": 50000,
            "maxStorage": 5368709120
          },
          "settings": {
            "theme": "light",
            "language": "en-US",
            "timezone": "America/New_York"
          }
        },
        "remark": "企业客户A",
        "createdAt": "2024-01-03T00:00:00+08:00",
        "updatedAt": "2024-01-03T00:00:00+08:00"
      }
    ],
    "total": 3,
    "page": 1,
    "pageSize": 20
  },
  "timestamp": 1703123456789,
  "traceId": "abc123def456"
}
```

### 1.2 创建租户

**接口地址**: `POST /tenant/create`

**请求头**:
```
Authorization: Bearer <token>
Content-Type: application/json
```

**请求参数**:
```json
{
  "name": "新企业租户",
  "code": "new_company",
  "domain": "newcompany.example.com",
  "maxUsers": 200,
  "storageLimit": 2147483648,
  "expireAt": "2025-12-31T23:59:59+08:00",
  "adminName": "company_admin",
  "adminEmail": "admin@newcompany.com",
  "adminPassword": "MTIzNDU2",
  "remark": "新企业客户"
}
```

**参数说明**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| name | string | 是 | 租户名称，长度1-100字符 |
| code | string | 是 | 租户编码，长度1-50字符，全局唯一 |
| domain | string | 否 | 租户域名，长度0-100字符 |
| maxUsers | int | 是 | 最大用户数，1-10000 |
| storageLimit | int64 | 是 | 存储限制（字节），最小0 |
| expireAt | string | 否 | 过期时间，ISO 8601格式 |
| adminName | string | 是 | 管理员用户名，长度1-50字符 |
| adminEmail | string | 是 | 管理员邮箱，需符合邮箱格式 |
| adminPassword | string | 是 | 管理员密码（Base64编码），长度6-32位 |
| remark | string | 否 | 备注，长度0-500字符 |

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 4,
    "name": "新企业租户",
    "code": "new_company",
    "domain": "newcompany.example.com",
    "status": 1,
    "statusName": "正常",
    "maxUsers": 200,
    "storageLimit": 2147483648,
    "expireAt": "2025-12-31T23:59:59+08:00",
    "adminUserId": 4,
    "adminName": "company_admin",
    "config": null,
    "remark": "新企业客户",
    "createdAt": "2024-01-04T10:30:00+08:00",
    "updatedAt": "2024-01-04T10:30:00+08:00"
  },
  "timestamp": 1703123456789,
  "traceId": "abc123def456"
}
```

### 1.3 更新租户

**接口地址**: `PUT /tenant/update`

**请求头**:
```
Authorization: Bearer <token>
Content-Type: application/json
```

**请求参数**:
```json
{
  "id": 2,
  "name": "更新后的测试租户",
  "domain": "updated-test.example.com",
  "maxUsers": 150,
  "storageLimit": 2147483648,
  "expireAt": "2025-06-30T23:59:59+08:00",
  "remark": "更新后的测试环境租户"
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 2,
    "name": "更新后的测试租户",
    "code": "test",
    "domain": "updated-test.example.com",
    "status": 1,
    "statusName": "正常",
    "maxUsers": 150,
    "storageLimit": 2147483648,
    "expireAt": "2025-06-30T23:59:59+08:00",
    "adminUserId": 2,
    "adminName": "test_admin",
    "config": null,
    "remark": "更新后的测试环境租户",
    "createdAt": "2024-01-02T00:00:00+08:00",
    "updatedAt": "2024-01-04T10:35:00+08:00"
  },
  "timestamp": 1703123456789,
  "traceId": "abc123def456"
}
```

### 1.4 删除租户

**接口地址**: `DELETE /tenant/delete`

**请求头**:
```
Authorization: Bearer <token>
Content-Type: application/json
```

**请求参数**:
```json
{
  "id": 4
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": null,
  "timestamp": 1703123456789,
  "traceId": "abc123def456"
}
```

### 1.5 获取租户详情

**接口地址**: `GET /tenant/detail`

**请求头**:
```
Authorization: Bearer <token>
```

**请求参数**:
```json
{
  "id": 2
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 2,
    "name": "测试租户",
    "code": "test",
    "domain": "test.example.com",
    "status": 1,
    "statusName": "正常",
    "maxUsers": 100,
    "storageLimit": 1073741824,
    "expireAt": "2024-12-31T23:59:59+08:00",
    "adminUserId": 2,
    "adminName": "test_admin",
    "config": {
      "features": {
        "advancedReports": false,
        "apiAccess": false,
        "customBranding": true
      },
      "limitations": {
        "maxApiCalls": 10000,
        "maxStorage": 1073741824
      },
      "settings": {
        "theme": "dark",
        "language": "zh-CN",
        "timezone": "Asia/Shanghai"
      }
    },
    "remark": "测试环境租户",
    "createdAt": "2024-01-02T00:00:00+08:00",
    "updatedAt": "2024-01-02T00:00:00+08:00",
    "stats": {
      "userCount": 25,
      "roleCount": 5,
      "menuCount": 12,
      "storageUsed": 268435456,
      "storageUsedText": "256.0 MB",
      "lastActiveTime": "2024-01-04T09:45:30+08:00"
    }
  },
  "timestamp": 1703123456789,
  "traceId": "abc123def456"
}
```

### 1.6 更新租户状态

**接口地址**: `PUT /tenant/status`

**请求头**:
```
Authorization: Bearer <token>
Content-Type: application/json
```

**请求参数**:
```json
{
  "id": 2,
  "status": 2
}
```

**参数说明**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 租户ID |
| status | int | 是 | 状态：1=正常 2=锁定 3=禁用 |

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": null,
  "timestamp": 1703123456789,
  "traceId": "abc123def456"
}
```

### 1.7 获取租户统计

**接口地址**: `GET /tenant/stats`

**请求头**:
```
Authorization: Bearer <token>
```

**请求参数**:
```json
{
  "id": 2
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "userCount": 25,
    "roleCount": 5,
    "menuCount": 12,
    "storageUsed": 268435456,
    "storageUsedText": "256.0 MB",
    "lastActiveTime": "2024-01-04T09:45:30+08:00"
  },
  "timestamp": 1703123456789,
  "traceId": "abc123def456"
}
```

### 1.8 更新租户配置

**接口地址**: `PUT /tenant/config`

**请求头**:
```
Authorization: Bearer <token>
Content-Type: application/json
```

**请求参数**:
```json
{
  "id": 2,
  "config": {
    "features": {
      "advancedReports": true,
      "apiAccess": true,
      "customBranding": true
    },
    "limitations": {
      "maxApiCalls": 20000,
      "maxStorage": 2147483648
    },
    "settings": {
      "theme": "light",
      "language": "zh-CN",
      "timezone": "Asia/Shanghai"
    }
  }
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": null,
  "timestamp": 1703123456789,
  "traceId": "abc123def456"
}
```

### 1.9 获取租户选项

**接口地址**: `GET /tenant/options`

**请求头**:
```
Authorization: Bearer <token>
```

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "value": 1,
        "label": "系统租户",
        "code": "system"
      },
      {
        "value": 2,
        "label": "测试租户",
        "code": "test"
      },
      {
        "value": 3,
        "label": "企业租户A",
        "code": "company_a"
      }
    ]
  },
  "timestamp": 1703123456789,
  "traceId": "abc123def456"
}
```

---

## 2. 用户管理（多租户版本）

### 2.1 用户登录

**接口地址**: `POST /login`

**请求头**:
```
Content-Type: application/json
```

**请求参数**:
```json
{
  "tenantCode": "test",
  "username": "test_user",
  "password": "MTIzNDU2",
  "captchaId": "abc123",
  "captcha": "1234",
  "rememberMe": false
}
```

**参数说明**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| tenantCode | string | 是 | 租户编码，长度1-50字符 |
| username | string | 是 | 用户名，长度3-50字符 |
| password | string | 是 | 密码（Base64编码），长度6-32位 |
| captchaId | string | 是 | 验证码ID |
| captcha | string | 是 | 验证码，长度4-6位 |
| rememberMe | boolean | 否 | 记住我，默认false |

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refreshToken": "abc123def456ghi789",
    "tokenType": "Bearer",
    "expiresIn": 86400,
    "userInfo": {
      "id": 5,
      "tenantId": 2,
      "tenantCode": "test",
      "username": "test_user",
      "email": "user@test.com",
      "realName": "测试用户",
      "nickname": "测试",
      "avatar": "/upload/avatar/user.jpg",
      "gender": 1,
      "birthday": "1990-01-01T00:00:00+08:00",
      "deptId": 1,
      "position": "开发工程师",
      "status": 1,
      "loginAt": "2024-01-04T10:30:00+08:00",
      "loginCount": 25,
      "twoFactorEnabled": false,
      "emailVerifiedAt": "2024-01-02T10:00:00+08:00",
      "phoneVerifiedAt": null,
      "remark": "测试用户账号",
      "createdAt": "2024-01-02T00:00:00+08:00",
      "updatedAt": "2024-01-04T10:30:00+08:00"
    },
    "permissions": [
      "dashboard",
      "user:list",
      "user:create",
      "user:update",
      "role:list",
      "role:create"
    ],
    "menuIds": [1, 2, 3, 4, 5, 10, 11, 12]
  },
  "timestamp": 1703123456789,
  "traceId": "abc123def456"
}
```

### 2.2 生成验证码

**接口地址**: `GET /captcha`

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "captchaId": "abc123def456",
    "captchaImage": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAH..."
  },
  "timestamp": 1703123456789,
  "traceId": "abc123def456"
}
```

---

## 3. 角色管理（多租户版本）

### 3.1 获取角色列表

**接口地址**: `GET /role/list`

**请求头**:
```
Authorization: Bearer <token>
X-Tenant-Id: 2
```

**请求参数**:
```json
{
  "page": 1,
  "pageSize": 20,
  "name": "管理员",
  "code": "admin",
  "status": 1
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 5,
        "tenantId": 2,
        "name": "租户管理员",
        "code": "tenant_admin",
        "description": "租户管理员角色",
        "status": 1,
        "sort": 1,
        "dataScope": 1,
        "remark": "租户管理员，拥有租户内所有权限",
        "createdAt": "2024-01-02T00:00:00+08:00",
        "updatedAt": "2024-01-02T00:00:00+08:00"
      },
      {
        "id": 6,
        "tenantId": 2,
        "name": "部门经理",
        "code": "dept_manager",
        "description": "部门经理角色",
        "status": 1,
        "sort": 2,
        "dataScope": 2,
        "remark": "部门经理，管理部门内用户和数据",
        "createdAt": "2024-01-02T01:00:00+08:00",
        "updatedAt": "2024-01-02T01:00:00+08:00"
      },
      {
        "id": 7,
        "tenantId": 2,
        "name": "普通员工",
        "code": "employee",
        "description": "普通员工角色",
        "status": 1,
        "sort": 3,
        "dataScope": 4,
        "remark": "普通员工，只能查看和操作自己的数据",
        "createdAt": "2024-01-02T02:00:00+08:00",
        "updatedAt": "2024-01-02T02:00:00+08:00"
      }
    ],
    "total": 3,
    "page": 1,
    "pageSize": 20
  },
  "timestamp": 1703123456789,
  "traceId": "abc123def456"
}
```

### 3.2 创建角色

**接口地址**: `POST /role`

**请求头**:
```
Authorization: Bearer <token>
X-Tenant-Id: 2
Content-Type: application/json
```

**请求参数**:
```json
{
  "name": "项目经理",
  "code": "project_manager",
  "description": "项目经理角色，负责项目管理",
  "status": 1,
  "sort": 5,
  "dataScope": 2,
  "remark": "项目经理角色"
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 8,
    "tenantId": 2,
    "name": "项目经理",
    "code": "project_manager",
    "description": "项目经理角色，负责项目管理",
    "status": 1,
    "sort": 5,
    "dataScope": 2,
    "remark": "项目经理角色",
    "createdAt": "2024-01-04T10:30:00+08:00",
    "updatedAt": "2024-01-04T10:30:00+08:00"
  },
  "timestamp": 1703123456789,
  "traceId": "abc123def456"
}
```

### 3.3 获取角色菜单权限

**接口地址**: `GET /role/{id}/menus`

**请求头**:
```
Authorization: Bearer <token>
X-Tenant-Id: 2
```

**路径参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|———|
| id | int | 是 | 角色ID |

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "roleId": 5,
    "roleName": "租户管理员",
    "menuIds": [1, 2, 3, 4, 5, 6, 7, 8, 9, 10],
    "menus": [
      {
        "id": 1,
        "tenantId": 2,
        "parentId": 0,
        "title": "系统管理",
        "name": "System",
        "path": "/system",
        "component": "Layout",
        "icon": "system",
        "type": 1,
        "sort": 1,
        "status": 1,
        "visible": 1,
        "permission": "system",
        "redirect": "/system/user",
        "alwaysShow": 1,
        "breadcrumb": 1,
        "activeMenu": "",
        "remark": "系统管理目录",
        "createdAt": "2024-01-02T00:00:00+08:00",
        "updatedAt": "2024-01-02T00:00:00+08:00"
      },
      {
        "id": 2,
        "tenantId": 2,
        "parentId": 1,
        "title": "用户管理",
        "name": "User",
        "path": "/system/user",
        "component": "system/user/index",
        "icon": "user",
        "type": 2,
        "sort": 1,
        "status": 1,
        "visible": 1,
        "permission": "system:user:list",
        "redirect": "",
        "alwaysShow": 0,
        "breadcrumb": 1,
        "activeMenu": "",
        "remark": "用户管理菜单",
        "createdAt": "2024-01-02T00:00:00+08:00",
        "updatedAt": "2024-01-02T00:00:00+08:00"
      }
    ]
  },
  "timestamp": 1703123456789,
  "traceId": "abc123def456"
}
```

### 3.4 更新角色菜单权限

**接口地址**: `PUT /role/{id}/menus`

**请求头**:
```
Authorization: Bearer <token>
X-Tenant-Id: 2
Content-Type: application/json
```

**路径参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 角色ID |

**请求参数**:
```json
{
  "menuIds": [1, 2, 3, 4, 5, 10, 11, 12]
}
```

**参数说明**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| menuIds | array | 是 | 菜单ID数组 |

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": null,
  "timestamp": 1703123456789,
  "traceId": "abc123def456"
}
```

### 3.5 获取角色权限详情

**接口地址**: `GET /role/{id}/permissions`

**请求头**:
```
Authorization: Bearer <token>
X-Tenant-Id: 2
```

**路径参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 角色ID |

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "roleId": 5,
    "roleName": "租户管理员",
    "permissions": [
      "system:user:list",
      "system:user:create",
      "system:user:update",
      "system:user:delete",
      "system:role:list",
      "system:role:create",
      "system:menu:list"
    ]
  },
  "timestamp": 1703123456789,
  "traceId": "abc123def456"
}
```

### 3.6 获取角色详情

**接口地址**: `GET /role/{id}`

**请求头**:
```
Authorization: Bearer <token>
X-Tenant-Id: 2
```

**路径参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 角色ID |

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 5,
    "tenantId": 2,
    "name": "租户管理员",
    "code": "tenant_admin",
    "description": "租户管理员角色",
    "status": 1,
    "sort": 1,
    "dataScope": 1,
    "remark": "租户管理员，拥有租户内所有权限",
    "createdAt": "2024-01-02T00:00:00+08:00",
    "updatedAt": "2024-01-02T00:00:00+08:00"
  },
  "timestamp": 1703123456789,
  "traceId": "abc123def456"
}
```

### 3.7 更新角色

**接口地址**: `PUT /role/{id}`

**请求头**:
```
Authorization: Bearer <token>
X-Tenant-Id: 2
Content-Type: application/json
```

**路径参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 角色ID |

**请求参数**:
```json
{
  "name": "高级管理员",
  "description": "高级管理员角色，权限更广泛",
  "status": 1,
  "sort": 2,
  "dataScope": 1,
  "remark": "高级管理员角色"
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 5,
    "tenantId": 2,
    "name": "高级管理员",
    "code": "tenant_admin",
    "description": "高级管理员角色，权限更广泛",
    "status": 1,
    "sort": 2,
    "dataScope": 1,
    "remark": "高级管理员角色",
    "createdAt": "2024-01-02T00:00:00+08:00",
    "updatedAt": "2024-01-04T11:00:00+08:00"
  },
  "timestamp": 1703123456789,
  "traceId": "abc123def456"
}
```

### 3.8 删除角色

**接口地址**: `DELETE /role/{id}`

**请求头**:
```
Authorization: Bearer <token>
X-Tenant-Id: 2
```

**路径参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 角色ID |

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": null,
  "timestamp": 1703123456789,
  "traceId": "abc123def456"
}
```

### 3.9 批量删除角色

**接口地址**: `DELETE /role/batch`

**请求头**:
```
Authorization: Bearer <token>
X-Tenant-Id: 2
Content-Type: application/json
```

**请求参数**:
```json
{
  "ids": [6, 7, 8]
}
```

**参数说明**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| ids | array | 是 | 角色ID数组 |

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": null,
  "timestamp": 1703123456789,
  "traceId": "abc123def456"
}
```

### 3.10 更新角色状态

**接口地址**: `PUT /role/{id}/status`

**请求头**:
```
Authorization: Bearer <token>
X-Tenant-Id: 2
Content-Type: application/json
```

**路径参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 角色ID |

**请求参数**:
```json
{
  "status": 2
}
```

**参数说明**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| status | int | 是 | 状态：1=启用 2=禁用 |

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": null,
  "timestamp": 1703123456789,
  "traceId": "abc123def456"
}
```

### 3.11 复制角色

**接口地址**: `POST /role/{id}/copy`

**请求头**:
```
Authorization: Bearer <token>
X-Tenant-Id: 2
Content-Type: application/json
```

**路径参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 源角色ID |

**请求参数**:
```json
{
  "name": "复制的项目经理",
  "code": "copied_project_manager"
}
```

**参数说明**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| name | string | 是 | 新角色名称 |
| code | string | 是 | 新角色编码 |

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 9,
    "tenantId": 2,
    "name": "复制的项目经理",
    "code": "copied_project_manager",
    "description": "项目经理角色，负责项目管理",
    "status": 1,
    "sort": 5,
    "dataScope": 2,
    "remark": "项目经理角色",
    "createdAt": "2024-01-04T11:30:00+08:00",
    "updatedAt": "2024-01-04T11:30:00+08:00"
  },
  "timestamp": 1703123456789,
  "traceId": "abc123def456"
}
```

---

## 4. 菜单管理（多租户版本）

### 4.1 获取菜单列表

**接口地址**: `GET /menu/list`

**请求头**:
```
Authorization: Bearer <token>
X-Tenant-Id: 2
```

**请求参数**:
```json
{
  "title": "用户",
  "status": 1,
  "type": 2
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "tenantId": 2,
        "parentId": 0,
        "title": "系统管理",
        "name": "System",
        "path": "/system",
        "component": "Layout",
        "icon": "system",
        "type": 1,
        "sort": 1,
        "status": 1,
        "visible": 1,
        "permission": "system",
        "redirect": "/system/user",
        "alwaysShow": 1,
        "breadcrumb": 1,
        "activeMenu": "",
        "remark": "系统管理目录",
        "createdAt": "2024-01-02T00:00:00+08:00",
        "updatedAt": "2024-01-02T00:00:00+08:00",
        "children": [
          {
            "id": 2,
            "tenantId": 2,
            "parentId": 1,
            "title": "用户管理",
            "name": "User",
            "path": "/system/user",
            "component": "system/user/index",
            "icon": "user",
            "type": 2,
            "sort": 1,
            "status": 1,
            "visible": 1,
            "permission": "system:user:list",
            "redirect": "",
            "alwaysShow": 0,
            "breadcrumb": 1,
            "activeMenu": "",
            "remark": "用户管理菜单",
            "createdAt": "2024-01-02T00:00:00+08:00",
            "updatedAt": "2024-01-02T00:00:00+08:00",
            "children": [
              {
                "id": 10,
                "tenantId": 2,
                "parentId": 2,
                "title": "新增用户",
                "name": "",
                "path": "",
                "component": "",
                "icon": "",
                "type": 3,
                "sort": 1,
                "status": 1,
                "visible": 1,
                "permission": "system:user:create",
                "redirect": "",
                "alwaysShow": 0,
                "breadcrumb": 1,
                "activeMenu": "",
                "remark": "新增用户按钮",
                "createdAt": "2024-01-02T00:00:00+08:00",
                "updatedAt": "2024-01-02T00:00:00+08:00",
                "children": []
              },
              {
                "id": 11,
                "tenantId": 2,
                "parentId": 2,
                "title": "编辑用户",
                "name": "",
                "path": "",
                "component": "",
                "icon": "",
                "type": 3,
                "sort": 2,
                "status": 1,
                "visible": 1,
                "permission": "system:user:update",
                "redirect": "",
                "alwaysShow": 0,
                "breadcrumb": 1,
                "activeMenu": "",
                "remark": "编辑用户按钮",
                "createdAt": "2024-01-02T00:00:00+08:00",
                "updatedAt": "2024-01-02T00:00:00+08:00",
                "children": []
              }
            ]
          },
          {
            "id": 3,
            "tenantId": 2,
            "parentId": 1,
            "title": "角色管理",
            "name": "Role",
            "path": "/system/role",
            "component": "system/role/index",
            "icon": "role",
            "type": 2,
            "sort": 2,
            "status": 1,
            "visible": 1,
            "permission": "system:role:list",
            "redirect": "",
            "alwaysShow": 0,
            "breadcrumb": 1,
            "activeMenu": "",
            "remark": "角色管理菜单",
            "createdAt": "2024-01-02T00:00:00+08:00",
            "updatedAt": "2024-01-02T00:00:00+08:00",
            "children": []
          }
        ]
      },
      {
        "id": 5,
        "tenantId": 2,
        "parentId": 0,
        "title": "仪表盘",
        "name": "Dashboard",
        "path": "/dashboard",
        "component": "dashboard/index",
        "icon": "dashboard",
        "type": 2,
        "sort": 0,
        "status": 1,
        "visible": 1,
        "permission": "dashboard",
        "redirect": "",
        "alwaysShow": 0,
        "breadcrumb": 1,
        "activeMenu": "",
        "remark": "仪表盘",
        "createdAt": "2024-01-02T00:00:00+08:00",
        "updatedAt": "2024-01-02T00:00:00+08:00",
        "children": []
      }
    ]
  },
  "timestamp": 1703123456789,
  "traceId": "abc123def456"
}
```

### 4.2 获取菜单树

**接口地址**: `GET /menu/tree`

**请求头**:
```
Authorization: Bearer <token>
X-Tenant-Id: 2
```

**请求参数**:
```json
{
  "status": 1,
  "type": 0
}
```

**参数说明**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| status | int | 否 | 状态筛选：0=全部 1=启用 2=禁用 |
| type | int | 否 | 类型筛选：0=全部 1=目录 2=菜单 3=按钮 |

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "tenantId": 2,
        "parentId": 0,
        "title": "系统管理",
        "name": "System",
        "path": "/system",
        "component": "Layout",
        "icon": "system",
        "type": 1,
        "sort": 1,
        "status": 1,
        "visible": 1,
        "permission": "system",
        "redirect": "/system/user",
        "alwaysShow": 1,
        "breadcrumb": 1,
        "activeMenu": "",
        "remark": "系统管理目录",
        "createdAt": "2024-01-02T00:00:00+08:00",
        "updatedAt": "2024-01-02T00:00:00+08:00",
        "children": [
          {
            "id": 2,
            "tenantId": 2,
            "parentId": 1,
            "title": "用户管理",
            "name": "User",
            "path": "/system/user",
            "component": "system/user/index",
            "icon": "user",
            "type": 2,
            "sort": 1,
            "status": 1,
            "visible": 1,
            "permission": "system:user:list",
            "redirect": "",
            "alwaysShow": 0,
            "breadcrumb": 1,
            "activeMenu": "",
            "remark": "用户管理菜单",
            "createdAt": "2024-01-02T00:00:00+08:00",
            "updatedAt": "2024-01-02T00:00:00+08:00",
            "children": []
          }
        ]
      }
    ]
  },
  "timestamp": 1703123456789,
  "traceId": "abc123def456"
}
```

### 4.3 获取菜单详情

**接口地址**: `GET /menu/{id}`

**请求头**:
```
Authorization: Bearer <token>
X-Tenant-Id: 2
```

**路径参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 菜单ID |

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 2,
    "tenantId": 2,
    "parentId": 1,
    "title": "用户管理",
    "name": "User",
    "path": "/system/user",
    "component": "system/user/index",
    "icon": "user",
    "type": 2,
    "sort": 1,
    "status": 1,
    "visible": 1,
    "permission": "system:user:list",
    "redirect": "",
    "alwaysShow": 0,
    "breadcrumb": 1,
    "activeMenu": "",
    "remark": "用户管理菜单",
    "createdAt": "2024-01-02T00:00:00+08:00",
    "updatedAt": "2024-01-02T00:00:00+08:00"
  },
  "timestamp": 1703123456789,
  "traceId": "abc123def456"
}
```

### 4.4 创建菜单

**接口地址**: `POST /menu`

**请求头**:
```
Authorization: Bearer <token>
X-Tenant-Id: 2
Content-Type: application/json
```

**请求参数**:
```json
{
  "parentId": 1,
  "title": "部门管理",
  "name": "Dept",
  "path": "/system/dept",
  "component": "system/dept/index",
  "icon": "dept",
  "type": 2,
  "sort": 3,
  "status": 1,
  "visible": 1,
  "permission": "system:dept:list",
  "redirect": "",
  "alwaysShow": 0,
  "breadcrumb": 1,
  "activeMenu": "",
  "remark": "部门管理菜单"
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 15,
    "tenantId": 2,
    "parentId": 1,
    "title": "部门管理",
    "name": "Dept",
    "path": "/system/dept",
    "component": "system/dept/index",
    "icon": "dept",
    "type": 2,
    "sort": 3,
    "status": 1,
    "visible": 1,
    "permission": "system:dept:list",
    "redirect": "",
    "alwaysShow": 0,
    "breadcrumb": 1,
    "activeMenu": "",
    "remark": "部门管理菜单",
    "createdAt": "2024-01-04T12:00:00+08:00",
    "updatedAt": "2024-01-04T12:00:00+08:00"
  },
  "timestamp": 1703123456789,
  "traceId": "abc123def456"
}
```

### 4.5 更新菜单

**接口地址**: `PUT /menu/{id}`

**请求头**:
```
Authorization: Bearer <token>
X-Tenant-Id: 2
Content-Type: application/json
```

**路径参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 菜单ID |

**请求参数**:
```json
{
  "title": "部门管理(更新)",
  "sort": 5,
  "remark": "部门管理菜单(已更新)"
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 15,
    "tenantId": 2,
    "parentId": 1,
    "title": "部门管理(更新)",
    "name": "Dept",
    "path": "/system/dept",
    "component": "system/dept/index",
    "icon": "dept",
    "type": 2,
    "sort": 5,
    "status": 1,
    "visible": 1,
    "permission": "system:dept:list",
    "redirect": "",
    "alwaysShow": 0,
    "breadcrumb": 1,
    "activeMenu": "",
    "remark": "部门管理菜单(已更新)",
    "createdAt": "2024-01-04T12:00:00+08:00",
    "updatedAt": "2024-01-04T12:30:00+08:00"
  },
  "timestamp": 1703123456789,
  "traceId": "abc123def456"
}
```

### 4.6 删除菜单

**接口地址**: `DELETE /menu/{id}`

**请求头**:
```
Authorization: Bearer <token>
X-Tenant-Id: 2
```

**路径参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 菜单ID |

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": null,
  "timestamp": 1703123456789,
  "traceId": "abc123def456"
}
```

### 4.7 更新菜单状态

**接口地址**: `PUT /menu/{id}/status`

**请求头**:
```
Authorization: Bearer <token>
X-Tenant-Id: 2
Content-Type: application/json
```

**路径参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int | 是 | 菜单ID |

**请求参数**:
```json
{
  "status": 2
}
```

**参数说明**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| status | int | 是 | 状态：1=启用 2=禁用 |

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": null,
  "timestamp": 1703123456789,
  "traceId": "abc123def456"
}
```

### 4.8 批量删除菜单

**接口地址**: `DELETE /menu/batch/delete`

**请求头**:
```
Authorization: Bearer <token>
X-Tenant-Id: 2
Content-Type: application/json
```

**请求参数**:
```json
{
  "ids": [10, 11, 12]
}
```

**参数说明**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| ids | array | 是 | 菜单ID数组 |

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": null,
  "timestamp": 1703123456789,
  "traceId": "abc123def456"
}
```

### 4.9 获取菜单选项

**接口地址**: `GET /menu/options`

**请求头**:
```
Authorization: Bearer <token>
X-Tenant-Id: 2
```

**请求参数**:
```json
{
  "type": 1
}
```

**参数说明**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| type | int | 否 | 类型筛选：0=全部 1=目录 2=菜单 3=按钮 |

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "value": 0,
        "label": "根目录",
        "type": 0
      },
      {
        "value": 1,
        "label": "系统管理",
        "type": 1
      },
      {
        "value": 5,
        "label": "仪表盘",
        "type": 2
      }
    ]
  },
  "timestamp": 1703123456789,
  "traceId": "abc123def456"
}
```

### 4.10 获取前端路由

**接口地址**: `GET /menu/routers`

**请求头**:
```
Authorization: Bearer <token>
X-Tenant-Id: 2
```

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "path": "/dashboard",
        "name": "Dashboard",
        "component": "dashboard/index",
        "redirect": "",
        "meta": {
          "title": "仪表盘",
          "icon": "dashboard",
          "roles": ["admin", "user"],
          "noCache": false,
          "affix": true
        }
      },
      {
        "path": "/system",
        "name": "System",
        "component": "Layout",
        "redirect": "/system/user",
        "meta": {
          "title": "系统管理",
          "icon": "system",
          "roles": ["admin"],
          "alwaysShow": true
        },
        "children": [
          {
            "path": "user",
            "name": "User",
            "component": "system/user/index",
            "meta": {
              "title": "用户管理",
              "icon": "user",
              "roles": ["admin"]
            }
          }
        ]
      }
    ]
  },
  "timestamp": 1703123456789,
  "traceId": "abc123def456"
}
```

---

## 5. 用户管理接口补充

### 5.1 用户退出

**接口地址**: `POST /logout`

**请求头**:
```
Authorization: Bearer <token>
```

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": null,
  "timestamp": 1703123456789,
  "traceId": "abc123def456"
}
```

### 5.2 获取用户信息

**接口地址**: `GET /profile`

**请求头**:
```
Authorization: Bearer <token>
```

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 5,
    "tenantId": 2,
    "tenantCode": "test",
    "username": "test_user",
    "email": "user@test.com",
    "realName": "测试用户",
    "nickname": "测试",
    "avatar": "/upload/avatar/user.jpg",
    "gender": 1,
    "birthday": "1990-01-01T00:00:00+08:00",
    "deptId": 1,
    "position": "开发工程师",
    "status": 1,
    "loginAt": "2024-01-04T10:30:00+08:00",
    "loginCount": 25,
    "twoFactorEnabled": false,
    "emailVerifiedAt": "2024-01-02T10:00:00+08:00",
    "phoneVerifiedAt": null,
    "remark": "测试用户账号",
    "createdAt": "2024-01-02T00:00:00+08:00",
    "updatedAt": "2024-01-04T10:30:00+08:00"
  },
  "timestamp": 1703123456789,
  "traceId": "abc123def456"
}
```

### 5.3 刷新访问令牌

**接口地址**: `POST /refresh-token`

**请求头**:
```
Content-Type: application/json
```

**请求参数**:
```json
{
  "refreshToken": "abc123def456ghi789"
}
```

**参数说明**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| refreshToken | string | 是 | 刷新令牌 |

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refreshToken": "new_refresh_token_here",
    "expiresIn": 86400
  },
  "timestamp": 1703123456789,
  "traceId": "abc123def456"
}
```

### 5.4 修改密码

**接口地址**: `POST /change-password`

**请求头**:
```
Authorization: Bearer <token>
Content-Type: application/json
```

**请求参数**:
```json
{
  "oldPassword": "MTIzNDU2",
  "newPassword": "YWJjZGVm",
  "confirmPassword": "YWJjZGVm"
}
```

**参数说明**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| oldPassword | string | 是 | 原密码（Base64编码） |
| newPassword | string | 是 | 新密码（Base64编码） |
| confirmPassword | string | 是 | 确认新密码（Base64编码） |

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": null,
  "timestamp": 1703123456789,
  "traceId": "abc123def456"
}
```

---

## 6. 错误码说明

| 错误码 | HTTP状态码 | 说明 |
|--------|-----------|------|
| 0 | 200 | 成功 |
| 400 | 400 | 请求参数错误 |
| 401 | 401 | 未授权，需要登录 |
| 403 | 403 | 禁止访问，权限不足 |
| 404 | 404 | 资源不存在 |
| 500 | 500 | 服务器内部错误 |

### 租户相关错误

| 错误信息 | 说明 |
|----------|------|
| "租户不存在" | 指定的租户编码或ID不存在 |
| "租户已被禁用或锁定" | 租户状态不正常 |
| "租户已过期" | 租户已超过有效期 |
| "租户编码已存在" | 创建租户时编码重复 |
| "系统租户不能删除" | 尝试删除系统租户 |
| "无权限访问该租户" | 用户无权访问指定租户 |
| "租户验证失败" | 登录时租户编码验证失败 |

### 用户登录相关错误

| 错误信息 | 说明 |
|----------|------|
| "用户名或密码错误" | 在指定租户下用户名或密码不正确 |
| "验证码错误" | 验证码不正确或已过期 |
| "用户已被锁定" | 用户账号被锁定 |
| "用户已被禁用" | 用户账号被禁用 |
| "用户状态异常" | 用户状态不正常 |

---

## 6. 使用示例

### 6.1 完整的租户创建流程

```bash
# 1. 系统管理员登录
curl -X POST "http://localhost:8000/api/login" \
  -H "Content-Type: application/json" \
  -d '{
    "tenantCode": "system",
    "username": "admin",
    "password": "MTIzNDU2",
    "captchaId": "abc123",
    "captcha": "1234"
  }'

# 2. 创建新租户
curl -X POST "http://localhost:8000/api/tenant/create" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "name": "新企业租户",
    "code": "new_company",
    "domain": "newcompany.example.com",
    "maxUsers": 200,
    "storageLimit": 2147483648,
    "expireAt": "2025-12-31T23:59:59+08:00",
    "adminName": "company_admin",
    "adminEmail": "admin@newcompany.com",
    "adminPassword": "MTIzNDU2",
    "remark": "新企业客户"
  }'

# 3. 租户管理员登录
curl -X POST "http://localhost:8000/api/login" \
  -H "Content-Type: application/json" \
  -d '{
    "tenantCode": "new_company",
    "username": "company_admin",
    "password": "MTIzNDU2",
    "captchaId": "def456",
    "captcha": "5678"
  }'
```

### 6.2 租户切换示例

```bash
# 通过Header指定租户
curl -X GET "http://localhost:8000/api/user/list" \
  -H "Authorization: Bearer <token>" \
  -H "X-Tenant-Id: 2"

# 通过域名自动识别租户
curl -X GET "http://test.example.com:8000/api/user/list" \
  -H "Authorization: Bearer <token>"
```

---

## 7. 注意事项

1. **租户隔离**: 所有API都会自动根据租户进行数据隔离
2. **权限验证**: 除了基本的登录验证，还会验证用户是否属于当前租户
3. **系统管理员**: 系统管理员（租户ID=1）可以跨租户访问数据
4. **数据一致性**: 创建、更新、删除操作都会自动带上租户信息
5. **菜单共享**: 菜单是全局共享的，租户通过角色-菜单关联控制访问权限
6. **用户登录**: 必须提供租户编码，支持不同租户的相同用户名
7. **缓存策略**: 租户信息和权限信息会被缓存以提高性能
8. **监控告警**: 建议对租户资源使用情况进行监控

## 8. 多租户架构设计

### 数据隔离策略

- **用户表**: 通过 `tenant_id` 隔离，同一租户内用户名唯一
- **角色表**: 通过 `tenant_id` 隔离，同一租户内角色编码唯一
- **菜单表**: 全局共享，所有租户使用同一套菜单
- **关联表**: 通过 `tenant_id` 隔离，确保权限分配的租户独立性

### 权限控制流程

```
用户登录 → 验证租户 → 获取租户内角色 → 查询角色菜单权限 → 返回可访问菜单
```

### 租户识别方式

1. **登录时指定**: 通过 `tenantCode` 参数
2. **Header指定**: 通过 `X-Tenant-Id` Header
3. **域名识别**: 通过子域名自动识别租户
