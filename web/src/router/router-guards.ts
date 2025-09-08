import type { RouteRecordRaw } from 'vue-router';
import { isNavigationFailure, Router } from 'vue-router';
import { useUserStoreWidthOut } from '@/store/modules/user';
import { useAsyncRouteStoreWidthOut } from '@/store/modules/asyncRoute';
import { ACCESS_TOKEN } from '@/store/mutation-types';
import { storage } from '@/utils/Storage';
import { PageEnum } from '@/enums/pageEnum';
import { ErrorPageRoute } from '@/router/base';
import { jump } from '@/utils/http/axios';
import { getNowUrl } from '@/utils/urlUtils';

const LOGIN_PATH = PageEnum.BASE_LOGIN;
const whitePathList = [LOGIN_PATH]; // no redirect whitelist

export function createRouterGuards(router: Router) {
  const userStore = useUserStoreWidthOut();
  const asyncRouteStore = useAsyncRouteStoreWidthOut();
  router.beforeEach(async (to, from, next) => {
    const Loading = window['$loading'] || null;
    Loading && Loading.start();

    console.log('🛡️ 路由守卫:', { to: to.path, from: from.path, name: to.name });

    if (from.path === LOGIN_PATH && to.name === 'errorPage') {
      console.log('🔄 从登录页跳转到错误页，重定向到首页');
      next(PageEnum.BASE_HOME);
      return;
    }

    // Whitelist can be directly entered
    if (whitePathList.includes(to.path as PageEnum)) {
      console.log('✅ 白名单路径，直接通过');
      next();
      return;
    }

    const token = storage.get(ACCESS_TOKEN);

    if (!token) {
      console.log('❌ 未登录，重定向到登录页');
      // You can access without permissions. You need to set the routing meta.ignoreAuth to true
      if (to.meta.ignoreAuth) {
        next();
        return;
      }

      // redirect login page
      const redirectData: { path: string; replace: boolean; query?: Recordable<string> } = {
        path: LOGIN_PATH,
        replace: true,
      };
      if (to.path) {
        redirectData.query = {
          ...redirectData.query,
          redirect: to.path,
        };
      }
      next(redirectData);
      return;
    }

    if (asyncRouteStore.getIsDynamicAddedRoute) {
      console.log('✅ 动态路由已添加，直接通过');
      next();
      return;
    }

    console.log('🔄 开始生成动态路由...');
    const redirectPath = (from.query.redirect || to.path) as string;
    const redirect = decodeURIComponent(redirectPath);
    const nextData = to.path === redirect ? { ...to, replace: true } : { path: redirect };
    
    try {
      const userInfo = await userStore.GetInfo();
      console.log('👤 获取用户信息成功:', userInfo);

      // 是否允许获取微信openid
      if (userStore.allowWxOpenId()) {
        let path = nextData.path;
        if (path === LOGIN_PATH) {
          path = PageEnum.BASE_HOME_REDIRECT;
        }

        const URI = getNowUrl() + '#' + path;
        jump('/wechat/authorize', { type: 'openId', syncRedirect: URI });
        return;
      }

      await userStore.GetConfig();
      const routes = await asyncRouteStore.generateRoutes(userInfo);
      console.log('🛤️ 动态路由生成完成:', routes);

      // 动态添加可访问路由表
      routes.forEach((item) => {
        router.addRoute(item as unknown as RouteRecordRaw);
      });

      //添加404
      const isErrorPage = router.getRoutes().findIndex((item) => item.name === ErrorPageRoute.name);
      if (isErrorPage === -1) {
        router.addRoute(ErrorPageRoute as unknown as RouteRecordRaw);
      }

      asyncRouteStore.setDynamicAddedRoute(true);
      console.log('✅ 路由添加完成，跳转到:', nextData);
      next(nextData);
    } catch (error) {
      console.error('❌ 动态路由生成失败:', error);
      // 如果路由生成失败，尝试跳转到默认首页
      next(PageEnum.BASE_HOME_REDIRECT);
    } finally {
      Loading && Loading.finish();
    }
  });

  router.afterEach((to, _, failure) => {
    document.title = (to?.meta?.title as string) || document.title;
    if (isNavigationFailure(failure)) {
      //console.log('failed navigation', failure)
    }
    const Loading = window['$loading'] || null;
    Loading && Loading.finish();
  });

  router.onError((error) => {
    console.log(error, '路由错误');
  });
}
