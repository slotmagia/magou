import { http } from '@/utils/http/axios';
import { ApiEnum } from '@/enums/apiEnum';

export interface BasicResponseModel<T = any> {
  code: number;
  message: string;
  data: T;
  timestamp: number;
  traceId: string;
}

// 菜单信息接口
export interface MenuInfo {
  id: number;
  tenantId: number;
  parentId: number;
  title: string;
  name: string;
  path: string;
  component: string;
  icon: string;
  type: number; // 1=目录 2=菜单 3=按钮
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
  children?: MenuInfo[];
}

// 菜单列表查询参数
export interface MenuListParams {
  title?: string;
  status?: number;
  type?: number;
}

// 创建菜单参数
export interface CreateMenuParams {
  parentId: number;
  title: string;
  name: string;
  path: string;
  component: string;
  icon?: string;
  type: number;
  sort?: number;
  status?: number;
  visible?: number;
  permission?: string;
  redirect?: string;
  alwaysShow?: number;
  breadcrumb?: number;
  activeMenu?: string;
  remark?: string;
}

// 更新菜单参数
export interface UpdateMenuParams extends CreateMenuParams {
  id: number;
}

/**
 * @description: 根据用户id获取用户菜单
 */
export function adminMenus() {
  return http.request({
    url: ApiEnum.RoleDynamic,
    method: 'GET',
  });
}

/**
 * @description: 获取菜单列表（多租户版本）
 */
export function getMenuList(params?: MenuListParams) {
  return http.request<BasicResponseModel<{
    list: MenuInfo[];
  }>>({
    url: '/menu/list',
    method: 'GET',
    params,
  });
}

/**
 * @description: 创建菜单
 */
export function createMenu(params: CreateMenuParams) {
  return http.request<BasicResponseModel<MenuInfo>>({
    url: '/menu',
    method: 'POST',
    data: params,
  });
}

/**
 * @description: 更新菜单
 */
export function updateMenu(params: UpdateMenuParams) {
  const { id, ...data } = params;
  return http.request<BasicResponseModel<MenuInfo>>({
    url: `/menu/${id}`,
    method: 'PUT',
    data,
  });
}

/**
 * @description: 删除菜单（新版本）
 */
export function deleteMenu(id: number) {
  return http.request<BasicResponseModel<null>>({
    url: `/menu/${id}`,
    method: 'DELETE',
  });
}

/**
 * @description: 批量删除菜单
 */
export function batchDeleteMenus(ids: number[]) {
  return http.request<BasicResponseModel<null>>({
    url: '/menu/batch/delete',
    method: 'DELETE',
    data: { ids },
  });
}

/**
 * @description: 更新菜单状态
 */
export function updateMenuStatus(id: number, status: number) {
  return http.request<BasicResponseModel<null>>({
    url: `/menu/${id}/status`,
    method: 'PUT',
    data: { status },
  });
}

/**
 * @description: 编辑菜单（兼容旧版本）
 */
export function EditMenu(params?) {
  return http.request({
    url: '/menu/edit',
    method: 'POST',
    params,
  });
}

/**
 * @description: 删除菜单
 */
export function DeleteMenu(params?) {
  return http.request({
    url: '/menu/delete',
    method: 'POST',
    params,
  });
}

/**
 * @description: 获取菜单详情
 */
export function getMenuDetail(id: number) {
  return http.request<BasicResponseModel<MenuInfo>>({
    url: `/menu/${id}`,
    method: 'GET',
  });
}

/**
 * @description: 获取菜单选项
 */
export function getMenuOptions(params?: { type?: number }) {
  return http.request<BasicResponseModel<{
    list: Array<{
      value: number;
      label: string;
      type: number;
    }>;
  }>>({
    url: '/menu/options',
    method: 'GET',
    params,
  });
}

/**
 * @description: 获取前端路由
 */
export function getMenuRouters() {
  return http.request<BasicResponseModel<{
    list: Array<{
      path: string;
      name: string;
      component: string;
      redirect?: string;
      meta: {
        title: string;
        icon: string;
        roles: string[];
        noCache?: boolean;
        affix?: boolean;
        alwaysShow?: boolean;
      };
      children?: any[];
    }>;
  }>>({
    url: ApiEnum.MenuRouters,
    method: 'GET',
  });
}

/**
 * @description: 获取菜单树形结构
 */
export function getMenuTree(params?: { status?: number }) {
  return http.request<BasicResponseModel<MenuInfo[]>>({
    url: '/menu/tree',
    method: 'GET',
    params,
  });
}

/**
 * @description: 获取用户可访问的菜单
 */
export function getUserMenus() {
  return http.request<BasicResponseModel<MenuInfo[]>>({
    url: '/menu/user-menus',
    method: 'GET',
  });
}
