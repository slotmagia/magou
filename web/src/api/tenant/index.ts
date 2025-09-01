import { http } from '@/utils/http/axios';

export interface BasicResponseModel<T = any> {
  code: number;
  message: string;
  data: T;
  timestamp: number;
  traceId: string;
}

export interface BasicPageParams {
  page?: number;
  pageSize?: number;
}

// 租户信息接口
export interface TenantInfo {
  id: number;
  name: string;
  code: string;
  domain?: string;
  status: number;
  statusName: string;
  maxUsers: number;
  storageLimit: number;
  expireAt?: string;
  adminUserId: number;
  adminName: string;
  config?: TenantConfig;
  remark?: string;
  createdAt: string;
  updatedAt: string;
  stats?: TenantStats;
}

// 租户配置接口
export interface TenantConfig {
  features: {
    advancedReports: boolean;
    apiAccess: boolean;
    customBranding: boolean;
  };
  limitations: {
    maxApiCalls: number;
    maxStorage: number;
  };
  settings: {
    theme: string;
    language: string;
    timezone: string;
  };
}

// 租户统计信息接口
export interface TenantStats {
  userCount: number;
  roleCount: number;
  menuCount: number;
  storageUsed: number;
  storageUsedText: string;
  lastActiveTime: string;
}

// 租户列表查询参数
export interface TenantListParams extends BasicPageParams {
  name?: string;
  code?: string;
  domain?: string;
  status?: number;
}

// 创建租户参数
export interface CreateTenantParams {
  name: string;
  code: string;
  domain?: string;
  maxUsers: number;
  storageLimit: number;
  expireAt?: string;
  adminName: string;
  adminEmail: string;
  adminPassword: string;
  remark?: string;
}

// 更新租户参数
export interface UpdateTenantParams {
  id: number;
  name?: string;
  domain?: string;
  maxUsers?: number;
  storageLimit?: number;
  expireAt?: string;
  remark?: string;
}

// 租户状态更新参数
export interface UpdateTenantStatusParams {
  id: number;
  status: number; // 1=正常 2=锁定 3=禁用
}

// 租户选项
export interface TenantOption {
  value: number;
  label: string;
  code: string;
}

/**
 * @description: 获取租户列表
 */
export function getTenantList(params?: TenantListParams) {
  return http.request<BasicResponseModel<{
    list: TenantInfo[];
    total: number;
    page: number;
    pageSize: number;
  }>>({
    url: '/tenant/list',
    method: 'GET',
    params,
  });
}

/**
 * @description: 创建租户
 */
export function createTenant(params: CreateTenantParams) {
  return http.request<BasicResponseModel<TenantInfo>>({
    url: '/tenant/create',
    method: 'POST',
    params,
  });
}

/**
 * @description: 更新租户
 */
export function updateTenant(params: UpdateTenantParams) {
  return http.request<BasicResponseModel<TenantInfo>>({
    url: '/tenant/update',
    method: 'PUT',
    params,
  });
}

/**
 * @description: 删除租户
 */
export function deleteTenant(params: { id: number }) {
  return http.request<BasicResponseModel<null>>({
    url: '/tenant/delete',
    method: 'DELETE',
    params,
  });
}

/**
 * @description: 获取租户详情
 */
export function getTenantDetail(params: { id: number }) {
  return http.request<BasicResponseModel<TenantInfo>>({
    url: '/tenant/detail',
    method: 'GET',
    params,
  });
}

/**
 * @description: 更新租户状态
 */
export function updateTenantStatus(params: UpdateTenantStatusParams) {
  return http.request<BasicResponseModel<null>>({
    url: '/tenant/status',
    method: 'PUT',
    params,
  });
}

/**
 * @description: 获取租户统计
 */
export function getTenantStats(params: { id: number }) {
  return http.request<BasicResponseModel<TenantStats>>({
    url: '/tenant/stats',
    method: 'GET',
    params,
  });
}

/**
 * @description: 更新租户配置
 */
export function updateTenantConfig(params: { id: number; config: TenantConfig }) {
  return http.request<BasicResponseModel<null>>({
    url: '/tenant/config',
    method: 'PUT',
    params,
  });
}

/**
 * @description: 获取租户选项
 */
export function getTenantOptions() {
  return http.request<BasicResponseModel<{
    list: TenantOption[];
  }>>({
    url: '/tenant/options',
    method: 'GET',
  });
}
