import { defineStore } from 'pinia';
import { createStorage, storage } from '@/utils/Storage';
import { store } from '@/store';
import {
  ACCESS_TOKEN,
  CURRENT_CONFIG,
  CURRENT_DICT,
  CURRENT_LOGIN_CONFIG,
  CURRENT_USER,
  IS_LOCKSCREEN,
} from '@/store/mutation-types';
import { ResultEnum } from '@/enums/httpEnum';
import {
  getConfig,
  getLoginConfig,
  getUserInfo,
  login,
  logout,
  mobileLogin,
} from '@/api/system/user';
import { login as tenantLoginApi, getCaptcha } from '@/api/auth';
import { isWechatBrowser } from '@/utils/is';
import { DeptTypeEnum } from '@/enums/deptEnum';
const Storage = createStorage({ storage: localStorage });

export interface UserInfoState {
  id: number;
  deptName: string;
  deptType: string;
  roleName: string;
  cityLabel: string;
  permissions: string[];
  username: string;
  realName: string;
  avatar: string;
  balance: number;
  integral: number;
  sex: number;
  qq: string;
  email: string;
  mobile: string;
  birthday: string;
  cityId: number;
  address: string;
  cash: {
    name: string;
    account: string;
    payeeCode: string;
  };
  createdAt: string;
  loginCount: number;
  lastLoginAt: string;
  lastLoginIp: string;
  openId: string;
  inviteCode: string;
}

export interface ConfigState {
  domain: string;
  version: string;
  wsAddr: string;
  mode: string;
}

export interface LoginConfigState {
  loginRegisterSwitch: number;
  loginCaptchaSwitch: number;
  loginAutoOpenId: number;
  loginProtocol: string;
  loginPolicy: string;
}

export interface IUserState {
  token: string;
  username: string;
  realName: string;
  avatar: string;
  permissions: any[];
  info: UserInfoState | null;
  config: ConfigState | null;
  loginConfig: LoginConfigState | null;
}

export const useUserStore = defineStore({
  id: 'app-member',
  state: (): IUserState => ({
    token: Storage.get(ACCESS_TOKEN, ''),
    username: '',
    realName: '',
    avatar: '',
    permissions: [],
    info: Storage.get(CURRENT_USER, null),
    config: Storage.get(CURRENT_CONFIG, null),
    loginConfig: Storage.get(CURRENT_LOGIN_CONFIG, null),
  }),
  getters: {
    getToken(): string {
      return this.token;
    },
    getAvatar(): string {
      return this.avatar;
    },
    getUsername(): string {
      return this.username;
    },
    getRealName(): string {
      return this.realName;
    },
    getPermissions(): [any][] {
      return this.permissions;
    },
    getUserInfo(): UserInfoState | null {
      return this.info;
    },
    getConfig(): ConfigState | null {
      return this.config;
    },
    getLoginConfig(): LoginConfigState | null {
      return this.loginConfig;
    },
    isCompanyDept(): boolean {
      return this.info?.deptType == DeptTypeEnum.Company;
    },
    isTenantDept(): boolean {
      return this.info?.deptType == DeptTypeEnum.Tenant;
    },
    isMerchantDept(): boolean {
      return this.info?.deptType == DeptTypeEnum.Merchant;
    },
    isUserDept(): boolean {
      return this.info?.deptType == DeptTypeEnum.User;
    },
  },
  actions: {
    setToken(token: string) {
      this.token = token;
    },
    setAvatar(avatar: string) {
      this.avatar = avatar;
    },
    setUsername(username: string) {
      this.username = username;
    },
    setRealName(realName: string) {
      this.realName = realName;
    },
    setPermissions(permissions: string[]) {
      this.permissions = permissions;
    },
    setUserInfo(info: UserInfoState | null) {
      this.info = info;
    },
    setConfig(config: ConfigState | null) {
      this.config = config;
    },
    setLoginConfig(config: LoginConfigState | null) {
      this.loginConfig = config;
    },
    // 账号登录
    async login(userInfo) {
      return await this.handleLogin(login(userInfo));
    },
    // 手机号登录
    async mobileLogin(userInfo) {
      return await this.handleLogin(mobileLogin(userInfo));
    },
    // 多租户登录
    async tenantLogin(userInfo) {
      return await this.handleTenantLogin(tenantLoginApi(userInfo));
    },
    async handleLogin(request: Promise<any>) {
      try {
        const response = await request;
        const { data, code } = response;
        if (code === ResultEnum.SUCCESS) {
          const ex = 30 * 24 * 60 * 60 * 1000;
          storage.set(ACCESS_TOKEN, data.token, ex);
          storage.set(CURRENT_USER, data, ex);
          storage.set(IS_LOCKSCREEN, false);
          this.setToken(data.token);
          this.setUserInfo(data);
        }
        return Promise.resolve(response);
      } catch (e) {
        return Promise.reject(e);
      }
    },
    async handleTenantLogin(request: Promise<any>) {
      try {
        const response = await request;
        const { data, code } = response;
        if (code === ResultEnum.SUCCESS) {
          const ex = 30 * 24 * 60 * 60 * 1000;
          // 存储多租户相关信息
          storage.set(ACCESS_TOKEN, data.accessToken, ex);
          storage.set(CURRENT_USER, data.userInfo, ex);
          storage.set('TENANT_INFO', { 
            tenantId: data.userInfo.tenantId, 
            tenantCode: data.userInfo.tenantCode 
          }, ex);
          storage.set('USER_PERMISSIONS', data.permissions, ex);
          storage.set('USER_MENU_IDS', data.menuIds, ex);
          storage.set(IS_LOCKSCREEN, false);
          
          this.setToken(data.accessToken);
          this.setUserInfo(data.userInfo);
          this.setPermissions(data.permissions);
        }
        return Promise.resolve(response);
      } catch (e) {
        return Promise.reject(e);
      }
    },
    // 获取用户信息 - 完全使用本地存储，不再调用后端API
    GetInfo() {
      const that: any = this;
      return new Promise((resolve, reject) => {
        // 完全使用本地存储的用户信息（登录时已保存）
        const storedUserInfo = storage.get(CURRENT_USER);
        const storedPermissions = storage.get('USER_PERMISSIONS') || [];
        const storedMenuIds = storage.get('USER_MENU_IDS') || [];
        
        if (storedUserInfo) {
          console.log('使用本地存储的用户信息');
          // 构造与API响应相同格式的结果
          const result = {
            ...storedUserInfo,
            permissions: storedPermissions,
            menuIds: storedMenuIds
          } as UserInfoState;
          
          that.setPermissions(storedPermissions);
          that.setUserInfo(storedUserInfo);
          that.setAvatar(storedUserInfo.avatar || '');
          that.setUsername(storedUserInfo.username || '');
          that.setRealName(storedUserInfo.realName || '');
          
          resolve(result);
          return;
        }
        
        // 如果本地没有数据，说明用户未登录
        console.warn('本地无用户信息，用户可能未登录');
        reject(new Error('用户未登录或本地用户信息丢失'));
      });
    },
    // 获取基础配置 - 由于后端不提供此接口，使用默认配置
    GetConfig() {
      const that = this;
      return new Promise((resolve, reject) => {
        // 检查是否有缓存的配置
        const cachedConfig = storage.get(CURRENT_CONFIG);
        if (cachedConfig) {
          console.log('使用缓存的系统配置');
          that.setConfig(cachedConfig);
          resolve(cachedConfig);
          return;
        }

        // 使用默认配置，因为后端不提供 /api/site/config 接口
        const defaultConfig: ConfigState = {
          domain: window.location.origin,
          version: '2.0.0',
          wsAddr: '',
          mode: 'development'
        };
        
        console.log('使用默认系统配置，后端不提供 /api/site/config 接口');
        that.setConfig(defaultConfig);
        storage.set(CURRENT_CONFIG, defaultConfig);
        resolve(defaultConfig);
      });
    },
    // 获取登录配置
    LoadLoginConfig: function () {
      const that = this;
      return new Promise((resolve, reject) => {
        getLoginConfig()
          .then((res) => {
            const result = res as unknown as LoginConfigState;
            that.setLoginConfig(result);
            storage.set(CURRENT_LOGIN_CONFIG, result);
            resolve(res);
          })
          .catch((error) => {
            reject(error);
          });
      });
    },
    // 是否允许获取微信openid
    allowWxOpenId(): boolean {
      if (!isWechatBrowser()) {
        return false;
      }

      // 如果没有登录配置，默认不允许微信openid
      if (this.loginConfig === null) {
        return false;
      }

      if (this.loginConfig.loginAutoOpenId !== 1) {
        return false;
      }
      return this.info === null || this.info.openId === '';
    },
    // 登出
    async logout() {
      try {
        const response = await logout();
        const { code } = response;
        if (code === ResultEnum.SUCCESS) {
          this.setPermissions([]);
          this.setUserInfo(null);
          storage.remove(ACCESS_TOKEN);
          storage.remove(CURRENT_USER);
          storage.remove(CURRENT_DICT);
        }
        return Promise.resolve(response);
      } catch (e) {
        return Promise.reject(e);
      }
    },
  },
});

// Need to be used outside the setup
export function useUserStoreWidthOut() {
  return useUserStore(store);
}
