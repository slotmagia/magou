import { http } from '@/utils/http/axios';
import { ApiEnum } from '@/enums/apiEnum';

export interface BasicResponseModel<T = any> {
  code: number;
  message: string;
  data: T;
}

export interface BasicPageParams {
  pageNumber: number;
  pageSize: number;
  total: number;
}

// 多租户用户信息接口
export interface TenantUserInfo {
  id: number;
  tenantId: number;
  tenantCode: string;
  username: string;
  email: string;
  realName: string;
  nickname: string;
  avatar: string;
  gender: number;
  birthday: string;
  deptId: number;
  position: string;
  status: number;
  loginAt: string;
  loginCount: number;
  twoFactorEnabled: boolean;
  emailVerifiedAt?: string;
  phoneVerifiedAt?: string;
  remark?: string;
  createdAt: string;
  updatedAt: string;
}

// 多租户用户列表查询参数
export interface TenantUserListParams extends BasicPageParams {
  username?: string;
  email?: string;
  realName?: string;
  deptId?: number;
  status?: number;
  tenantId?: number;
}

// 创建多租户用户参数
export interface CreateTenantUserParams {
  tenantId?: number; // 如果不传则使用当前租户
  username: string;
  email: string;
  realName: string;
  nickname?: string;
  password: string; // Base64编码
  gender?: number;
  birthday?: string;
  deptId?: number;
  position?: string;
  status?: number;
  remark?: string;
}

// 用户个人资料接口
export interface UserProfile {
  id: number;
  tenantId: number;
  tenantCode: string;
  username: string;
  email: string;
  realName: string;
  nickname: string;
  avatar: string;
  gender: number;
  birthday: string;
  deptId: number;
  position: string;
  status: number;
  loginAt: string;
  loginCount: number;
  twoFactorEnabled: boolean;
  emailVerifiedAt?: string;
  phoneVerifiedAt?: string;
  remark?: string;
  createdAt: string;
  updatedAt: string;
}

// 修改密码参数
export interface ChangePasswordParams {
  oldPassword: string; // Base64编码
  newPassword: string; // Base64编码
  confirmPassword: string; // Base64编码
}

// 刷新令牌参数
export interface RefreshTokenParams {
  refreshToken: string;
}

export function getConfig() {
  return http.request({
    url: ApiEnum.SiteConfig,
    method: 'get',
    headers: { hostname: location.hostname },
  });
}

/**
 * @description: 获取用户信息
 */
export function getUserInfo() {
  return http.request({
    url: ApiEnum.MemberInfo,
    method: 'get',
  });
}

export function updateMemberProfile(params) {
  return http.request({
    url: '/member/updateProfile',
    method: 'post',
    params,
  });
}

export function updateMemberPwd(params) {
  return http.request({
    url: '/member/updatePwd',
    method: 'post',
    params,
  });
}

export function updateMemberMobile(params) {
  return http.request({
    url: '/member/updateMobile',
    method: 'post',
    params,
  });
}

export function updateMemberEmail(params) {
  return http.request({
    url: '/member/updateEmail',
    method: 'post',
    params,
  });
}

export function SendBindEmail() {
  return http.request({
    url: '/ems/sendBind',
    method: 'post',
  });
}

export function SendBindSms() {
  return http.request({
    url: '/sms/sendBind',
    method: 'post',
  });
}

export function SendSms(params) {
  return http.request({
    url: '/sms/send',
    method: 'post',
    params,
  });
}

export function updateMemberCash(params) {
  return http.request({
    url: '/member/updateCash',
    method: 'post',
    params,
  });
}

/**
 * @description: 用户登录配置
 */
export function getLoginConfig() {
  return http.request<BasicResponseModel>({
    url: ApiEnum.SiteLoginConfig,
    method: 'get',
  });
}

/**
 * @description: 用户注册
 */
export function register(params) {
  return http.request<BasicResponseModel>(
    {
      url: ApiEnum.SiteRegister,
      method: 'POST',
      params,
    },
    {
      isTransformResponse: false,
    }
  );
}

/**
 * @description: 用户登录
 */
export function login(params) {
  return http.request<BasicResponseModel>(
    {
      url: ApiEnum.SiteAccountLogin,
      method: 'POST',
      params,
    },
    {
      isTransformResponse: false,
    }
  );
}

/**
 * @description: 手机号登录
 */
export function mobileLogin(params) {
  return http.request<BasicResponseModel>(
    {
      url: ApiEnum.SiteMobileLogin,
      method: 'POST',
      params,
    },
    {
      isTransformResponse: false,
    }
  );
}

/**
 * @description: 用户注销
 */
export function logout() {
  return http.request<BasicResponseModel>(
    {
      url: ApiEnum.SiteLogout,
      method: 'POST',
    },
    {
      isTransformResponse: false,
    }
  );
}

/**
 * @description: 用户修改密码
 */
export function changePassword(params, uid) {
  return http.request(
    {
      url: `/user/u${uid}/changepw`,
      method: 'POST',
      params,
    },
    {
      isTransformResponse: false,
    }
  );
}

/**
 * @description: 获取多租户用户列表
 */
export function getTenantUserList(params?: TenantUserListParams) {
  return http.request<BasicResponseModel<{
    list: TenantUserInfo[];
    total: number;
    page: number;
    pageSize: number;
  }>>({
    url: '/tenant-user/list',
    method: 'GET',
    params,
  });
}

/**
 * @description: 创建多租户用户
 */
export function createTenantUser(params: CreateTenantUserParams) {
  return http.request<BasicResponseModel<TenantUserInfo>>({
    url: '/tenant-user/create',
    method: 'POST',
    params,
  });
}

/**
 * @description: 更新多租户用户
 */
export function updateTenantUser(params: CreateTenantUserParams & { id: number }) {
  return http.request<BasicResponseModel<TenantUserInfo>>({
    url: '/tenant-user/update',
    method: 'PUT',
    params,
  });
}

/**
 * @description: 删除多租户用户
 */
export function deleteTenantUser(params: { id: number }) {
  return http.request<BasicResponseModel<null>>({
    url: '/tenant-user/delete',
    method: 'DELETE',
    params,
  });
}

/**
 * @description: 获取多租户用户详情
 */
export function getTenantUserDetail(params: { id: number }) {
  return http.request<BasicResponseModel<TenantUserInfo>>({
    url: '/tenant-user/detail',
    method: 'GET',
    params,
  });
}

/**
 * @description: 更新多租户用户状态
 */
export function updateTenantUserStatus(params: { id: number; status: number }) {
  return http.request<BasicResponseModel<null>>({
    url: '/tenant-user/status',
    method: 'PUT',
    params,
  });
}

/**
 * @description: 重置多租户用户密码
 */
export function resetTenantUserPassword(params: { id: number; password: string }) {
  return http.request<BasicResponseModel<null>>({
    url: '/tenant-user/reset-password',
    method: 'PUT',
    params,
  });
}

/**
 * @description: 用户退出登录（新版本）
 */
export function logoutV2() {
  return http.request<BasicResponseModel<null>>({
    url: '/logout',
    method: 'POST',
  });
}

/**
 * @description: 获取用户个人资料
 */
export function getProfile() {
  return http.request<BasicResponseModel<UserProfile>>({
    url: '/profile',
    method: 'GET',
  });
}

/**
 * @description: 刷新访问令牌
 */
export function refreshToken(params: RefreshTokenParams) {
  return http.request<BasicResponseModel<{
    accessToken: string;
    refreshToken: string;
    expiresIn: number;
  }>>({
    url: '/refresh-token',
    method: 'POST',
    data: params,
  });
}

/**
 * @description: 修改密码（新版本）
 */
export function changePasswordV2(params: ChangePasswordParams) {
  return http.request<BasicResponseModel<null>>({
    url: '/change-password',
    method: 'POST',
    data: params,
  });
}
