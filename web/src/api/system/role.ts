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

// 角色信息接口
export interface RoleInfo {
  id: number;
  tenantId: number;
  name: string;
  code: string;
  description: string;
  status: number;
  sort: number;
  dataScope: number;
  remark?: string;
  createdAt: string;
  updatedAt: string;
}

// 角色列表查询参数
export interface RoleListParams extends BasicPageParams {
  name?: string;
  code?: string;
  status?: number;
  tenantId?: number;
}

// 创建角色参数
export interface CreateRoleParams {
  name: string;
  code: string;
  description?: string;
  status?: number;
  sort?: number;
  dataScope?: number;
  remark?: string;
}

// 更新角色参数
export interface UpdateRoleParams extends CreateRoleParams {
  id: number;
}

// 角色菜单权限
export interface RoleMenuPermission {
  roleId: number;
  roleName: string;
  menuIds: number[];
  menus: MenuInfo[];
}

// 菜单信息
export interface MenuInfo {
  id: number;
  tenantId: number;
  parentId: number;
  title: string;
  name: string;
  path: string;
  component: string;
  icon: string;
  type: number;
  sort: number;
  status: number;
  visible: number;
  permission: string;
  redirect: string;
  alwaysShow: number;
  breadcrumb: number;
  activeMenu: string;
  remark?: string;
  createdAt: string;
  updatedAt: string;
}

/**
 * @description: 获取角色列表（多租户版本）
 */
export function getRoleList(params?: RoleListParams) {
  return http.request<BasicResponseModel<{
    list: RoleInfo[];
    total: number;
    page: number;
    pageSize: number;
  }>>({
    url: '/role/list',
    method: 'GET',
    params,
    headers: {
      'X-Tenant-Id': params?.tenantId || undefined,
    },
  });
}

export function getRoleOption() {
  const params = { pageSize: 100 };
  return getRoleList(params);
}

/**
 * @description: 创建角色
 */
export function createRole(params: CreateRoleParams) {
  return http.request<BasicResponseModel<RoleInfo>>({
    url: '/role',
    method: 'POST',
    data: params,
  });
}

/**
 * @description: 更新角色
 */
export function updateRole(params: UpdateRoleParams) {
  const { id, ...data } = params;
  return http.request<BasicResponseModel<RoleInfo>>({
    url: `/role/${id}`,
    method: 'PUT',
    data,
  });
}

/**
 * @description: 获取角色详情
 */
export function getRoleDetail(id: number) {
  return http.request<BasicResponseModel<RoleInfo>>({
    url: `/role/${id}`,
    method: 'GET',
  });
}

/**
 * @description: 删除角色
 */
export function deleteRole(id: number) {
  return http.request<BasicResponseModel<null>>({
    url: `/role/${id}`,
    method: 'DELETE',
  });
}

/**
 * @description: 批量删除角色
 */
export function batchDeleteRoles(ids: number[]) {
  return http.request<BasicResponseModel<null>>({
    url: '/role/batch',
    method: 'DELETE',
    data: { ids },
  });
}

/**
 * @description: 更新角色状态
 */
export function updateRoleStatus(id: number, status: number) {
  return http.request<BasicResponseModel<null>>({
    url: `/role/${id}/status`,
    method: 'PUT',
    data: { status },
  });
}

/**
 * @description: 复制角色
 */
export function copyRole(id: number, params: { name: string; code: string }) {
  return http.request<BasicResponseModel<RoleInfo>>({
    url: `/role/${id}/copy`,
    method: 'POST',
    data: params,
  });
}

/**
 * @description: 编辑角色（兼容旧版本）
 */
export function Edit(params) {
  return http.request({
    url: '/role/edit',
    method: 'POST',
    params,
  });
}

/**
 * @description: 删除角色
 */
export function Delete(params) {
  return http.request({
    url: '/role/delete',
    method: 'POST',
    params,
  });
}

/**
 * @description: 更新角色权限
 */
export function UpdatePermissions(params) {
  return http.request({
    url: '/role/updatePermissions',
    method: 'POST',
    params,
  });
}

/**
 * @description: 获取角色权限
 */
export function GetPermissions(params) {
  return http.request({
    url: '/role/getPermissions',
    method: 'GET',
    params,
  });
}

/**
 * @description: 获取角色菜单权限（多租户版本）
 */
export function getRoleMenus(roleId: number) {
  return http.request<BasicResponseModel<RoleMenuPermission>>({
    url: `/role/${roleId}/menus`,
    method: 'GET',
  });
}

/**
 * @description: 更新角色菜单权限
 */
export function updateRoleMenus(roleId: number, menuIds: number[]) {
  return http.request<BasicResponseModel<null>>({
    url: `/role/${roleId}/menus`,
    method: 'PUT',
    data: { menuIds },
  });
}

/**
 * @description: 获取角色权限详情
 */
export function getRolePermissions(roleId: number) {
  return http.request<BasicResponseModel<{
    roleId: number;
    roleName: string;
    permissions: string[];
  }>>({
    url: `/role/${roleId}/permissions`,
    method: 'GET',
  });
}

/**
 * @description: 数据范围选择
 */
export function DataScopeSelect() {
  return http.request({
    url: '/role/dataScope/select',
    method: 'GET',
  });
}

/**
 * @description: 编辑数据范围
 */
export function DataScopeEdit(params) {
  return http.request({
    url: '/role/dataScope/edit',
    method: 'POST',
    params,
  });
}
