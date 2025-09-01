import { getMenuRouters } from '@/api/system/menu';
import { constantRouterIcon } from './router-icons';
import { RouteRecordRaw } from 'vue-router';
import { Layout, ParentLayout } from '@/router/constant';
import type { AppRouteRecordRaw } from '@/router/types';

const Iframe = () => import('@/views/iframe/index.vue');
const LayoutMap = new Map<string, () => Promise<typeof import('*.vue')>>();

LayoutMap.set('LAYOUT', Layout);
LayoutMap.set('IFRAME', Iframe);

/**
 * 格式化 后端 结构信息并递归生成层级路由表
 * @param routerMap
 * @param parent
 * @returns {*}
 */
export const routerGenerator = (routerMap, parent?): any[] => {
  console.log('🔧 routerGenerator 开始处理:', routerMap, '父路由:', parent?.path);
  
  return routerMap.map((item) => {
    console.log('📝 处理路由项:', item);
    // 处理路径，确保以/开头但避免重复/
    let routePath = item.path;
    if (parent && parent.path) {
      // 如果有父路由，拼接路径
      routePath = `${parent.path}/${item.path}`.replace('//', '/');
    }

    // 处理组件路径
    let componentPath = item.component;
    if (componentPath === 'Layout') {
      componentPath = 'LAYOUT';
    } else if (componentPath) {
      // 映射后端返回的组件路径到实际的Vue组件路径
      componentPath = mapComponentPath(componentPath);
    }

    // 生成路由名称
    const routeName = item.name || generateRouteName(item.path, item.id);

    // 处理菜单标题 - 如果title为空，使用路径生成默认标题
    const title = item.meta?.title || generateDefaultTitle(item.path);

    const currentRouter: any = {
      // 路由地址
      path: routePath,
      // 路由名称，建议唯一
      name: routeName,
      // 该路由对应页面的 组件
      component: componentPath,
      // meta: 页面标题, 菜单图标, 页面权限(供指令权限用，可去掉)
      meta: {
        ...item.meta,
        title: title,
        label: title,
        icon: constantRouterIcon[item.meta?.icon] || null,
        permissions: item.meta?.permissions || null,
        hidden: item.hidden || false,
        alwaysShow: item.alwaysShow || false,
        noCache: item.meta?.noCache || false,
        breadcrumb: item.meta?.breadcrumb !== false,
        type: item.meta?.type || (item.children ? 1 : 2), // 1=目录 2=菜单
      },
    };

    console.log(`✨ 生成路由: ${routePath} -> ${componentPath} (${title})`);

    // 处理重定向
    if (item.redirect) {
      currentRouter.redirect = item.redirect;
    }

    // 是否有子菜单，并递归处理
    if (item.children && item.children.length > 0) {
      // 如果是目录类型且未定义redirect，默认重定向到第一个子路由
      if (!item.redirect && (item.meta?.type === 1 || item.component === 'Layout')) {
        const firstChild = item.children[0];
        if (firstChild) {
          currentRouter.redirect = `${routePath}/${firstChild.path}`.replace('//', '/');
        }
      }
      // 递归处理子路由
      currentRouter.children = routerGenerator(item.children, currentRouter);
    }

    return currentRouter;
  });
};

/**
 * 生成路由名称
 */
function generateRouteName(path: string, id?: number): string {
  if (id) {
    return `Route${id}`;
  }
  // 将路径转换为驼峰命名
  return path
    .split('/')
    .filter(Boolean)
    .map((segment, index) => {
      if (index === 0) {
        return segment.toLowerCase();
      }
      return segment.charAt(0).toUpperCase() + segment.slice(1).toLowerCase();
    })
    .join('');
}

/**
 * 生成默认标题
 */
function generateDefaultTitle(path: string): string {
  const pathSegments = path.split('/').filter(Boolean);
  const lastSegment = pathSegments[pathSegments.length - 1];
  
  // 简单的路径到中文映射
  const titleMap: Record<string, string> = {
    'dashboard': '仪表板',
    'system': '系统管理',
    'user': '用户管理',
    'role': '角色管理',
    'menu': '菜单管理',
    'dept': '部门管理',
    'tenant': '租户管理',
    'list': '列表',
    'config': '配置'
  };
  
  return titleMap[lastSegment] || lastSegment || '未命名';
}

/**
 * 映射组件路径到实际的Vue组件文件路径
 */
function mapComponentPath(componentPath: string): string {
  // 组件路径映射表
  const componentMap: Record<string, string> = {
    // 仪表板相关
    'dashboard/index': '/dashboard/console/console',
    
    // 系统管理相关 - 映射到permission目录
    'system/user/index': '/org/user/user',
    'system/role/index': '/permission/role/role',
    'system/menu/index': '/permission/menu/menu',
    'system/dept/index': '/org/dept/dept',
    
    // 租户管理相关
    'tenant/list/index': '/tenant/index',
    'tenant/config/index': '/system/config/system',
  };

  // 如果在映射表中找到，使用映射的路径
  if (componentMap[componentPath]) {
    return componentMap[componentPath];
  }

  // 否则尝试直接转换
  if (!componentPath.startsWith('/')) {
    componentPath = `/${componentPath}`;
  }

  // 移除可能的.vue后缀
  if (componentPath.endsWith('.vue')) {
    componentPath = componentPath.slice(0, -4);
  }

  return componentPath;
}

/**
 * 动态生成菜单
 */
export const generatorDynamicRouter = (): Promise<RouteRecordRaw[]> => {
  return new Promise((resolve, reject) => {
    getMenuRouters()
      .then((result) => {
        console.log('📊 获取到菜单路由原始数据:', result);
        
        // 处理新API的响应格式
        const menuData = result.data?.list || result.list || [];
        console.log('📋 处理的菜单数据:', menuData);
        
        const routeList = routerGenerator(menuData);
        console.log('🛤️ 生成的路由列表:', routeList);
        
        asyncImportRoute(routeList);
        console.log('✅ 路由生成完成');

        resolve(routeList);
      })
      .catch((err) => {
        console.error('❌ 获取菜单路由失败:', err);
        // 如果API调用失败，提供一个默认的dashboard路由
        const defaultRoutes = [{
          path: '/dashboard',
          name: 'Dashboard',
          component: 'LAYOUT',
          redirect: '/dashboard/console',
          meta: {
            title: '仪表板',
            icon: 'DashboardOutlined',
            type: 1,
            sort: 1
          },
          children: [{
            path: 'console',
            name: 'Console',
            component: '/dashboard/console/console',
            meta: {
              title: '控制台',
              icon: 'DashboardOutlined',
              type: 2,
              sort: 1
            }
          }]
        }];
        
        console.log('📝 使用默认路由配置');
        const routeList = routerGenerator(defaultRoutes);
        asyncImportRoute(routeList);
        resolve(routeList);
      });
  });
};

/**
 * 查找views中对应的组件文件
 * */
let viewsModules: Record<string, () => Promise<Recordable>>;
export const asyncImportRoute = (routes: AppRouteRecordRaw[] | undefined): void => {
  viewsModules = viewsModules || import.meta.glob('../views/**/*.{vue,tsx}');
  if (!routes) return;
  routes.forEach((item) => {
    if (!item.component && item.meta?.frameSrc) {
      item.component = 'IFRAME';
    }
    const { component, name } = item;
    const { children } = item;
    if (component) {
      const layoutFound = LayoutMap.get(component as string);
      if (layoutFound) {
        item.component = layoutFound;
      } else {
        item.component = dynamicImport(viewsModules, component as string);
      }
    } else if (name) {
      item.component = ParentLayout;
    }
    children && asyncImportRoute(children);
  });
};

/**
 * 动态导入
 * */
export const dynamicImport = (
  viewsModules: Record<string, () => Promise<Recordable>>,
  component: string
) => {
  const keys = Object.keys(viewsModules);
  const matchKeys = keys.filter((key) => {
    let k = key.replace('../views', '');
    const lastIndex = k.lastIndexOf('.');
    k = k.substring(0, lastIndex);
    return k === component;
  });
  if (matchKeys?.length === 1) {
    const matchKey = matchKeys[0];
    return viewsModules[matchKey];
  }
  if (matchKeys?.length > 1) {
    console.warn(
      'Please do not create `.vue` and `.TSX` files with the same file name in the same hierarchical directory under the views folder. This will cause dynamic introduction failure'
    );
    return;
  }
};

/**
 * 移除隐藏的菜单
 * @param menus
 */
export const removeHiddenMenus = (menus: any[]) => {
  const arr: any[] = [];
  for (let j = 0; j < menus.length; j++) {
    if (menus[j].meta?.type === 3) {
      continue;
    }
    if (menus[j].meta?.hidden === true) {
      continue;
    }

    if (menus[j].children?.length > 0) {
      menus[j].children = removeHiddenMenus(menus[j].children);
      if (menus[j].children?.length === 0) {
        delete menus[j].children;
      }
    }
    arr.push(menus[j]);
  }
  return arr;
};
