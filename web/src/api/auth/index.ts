import { http } from '@/utils/http/axios';

export interface BasicResponseModel<T = any> {
  code: number;
  message: string;
  data: T;
  timestamp: number;
  traceId: string;
}

// 用户信息接口
export interface UserInfo {
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

// 登录参数
export interface LoginParams {
  tenantCode: string;
  username: string;
  password: string; // Base64编码
  captchaId: string;
  captcha: string;
  rememberMe?: boolean;
}

// 登录响应数据
export interface LoginResult {
  accessToken: string;
  refreshToken: string;
  tokenType: string;
  expiresIn: number;
  userInfo: UserInfo;
  permissions: string[];
  menuIds: number[];
}

// 验证码响应数据
export interface CaptchaResult {
  captchaId: string;
  captchaImage: string; // Base64图片数据
}

/**
 * @description: 用户登录（多租户版本）
 */
export function login(params: LoginParams) {
  return http.request<BasicResponseModel<LoginResult>>(
    {
      url: '/login',
      method: 'POST',
      data: params,
    },
    {
      isTransformResponse: false,
    }
  );
}

/**
 * @description: 生成验证码
 */
export function getCaptcha() {
  return http.request<BasicResponseModel<CaptchaResult>>({
    url: '/captcha',
    method: 'GET',
  });
}

/**
 * @description: 用户注销
 */
export function logout() {
  return http.request<BasicResponseModel<null>>(
    {
      url: '/logout',
      method: 'POST',
    },
    {
      isTransformResponse: false,
    }
  );
}

/**
 * @description: 刷新Token
 */
export function refreshToken(refreshToken: string) {
  return http.request<BasicResponseModel<{
    accessToken: string;
    refreshToken: string;
    expiresIn: number;
  }>>({
    url: '/refresh-token',
    method: 'POST',
    data: { refreshToken },
  });
}

/**
 * @description: 获取当前用户信息
 */
export function getCurrentUserInfo() {
  return http.request<BasicResponseModel<UserInfo>>({
    url: '/profile',
    method: 'GET',
  });
}

/**
 * @description: 修改密码
 */
export function changePassword(params: {
  oldPassword: string;
  newPassword: string;
  confirmPassword: string;
}) {
  return http.request<BasicResponseModel<null>>({
    url: '/change-password',
    method: 'POST',
    data: params,
  });
}
