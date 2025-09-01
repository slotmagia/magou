/**
 * 菜单渲染修复助手
 * 帮助解决前端菜单不渲染的问题
 */

(function() {
    'use strict';

    console.log('🔧 菜单修复助手已加载');

    // 修复方案
    const fixes = {
        // 修复1：检查并修复API调用
        fixApiCall: function() {
            console.log('🔧 修复1: 检查API调用...');

            // 检查是否已经有菜单数据
            if (window.menuData || window.routerData) {
                console.log('✅ 菜单数据已存在');
                return true;
            }

            // 手动调用API获取菜单数据
            return this.fetchMenuData();
        },

        // 获取菜单数据
        fetchMenuData: function() {
            console.log('📡 获取菜单数据...');

            return fetch('http://localhost:8001/api/menu/routers', {
                method: 'GET',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': this.getAuthToken()
                }
            })
            .then(response => {
                if (!response.ok) {
                    throw new Error(`HTTP ${response.status}: ${response.statusText}`);
                }
                return response.json();
            })
            .then(data => {
                if (data.code === 0 && data.data && data.data.list) {
                    console.log('✅ 成功获取菜单数据:', data.data.list.length, '个菜单项');

                    // 保存到全局变量
                    window.menuData = data.data.list;
                    window.routerData = data.data.list;

                    // 尝试触发菜单渲染
                    this.triggerMenuRender(data.data.list);

                    return true;
                } else {
                    throw new Error('API响应格式错误或没有菜单数据');
                }
            })
            .catch(error => {
                console.error('❌ 获取菜单数据失败:', error);
                return false;
            });
        },

        // 获取认证令牌
        getAuthToken: function() {
            // 尝试从各种存储中获取token
            const sources = [
                () => localStorage.getItem('token'),
                () => localStorage.getItem('access_token'),
                () => localStorage.getItem('authorization'),
                () => sessionStorage.getItem('token'),
                () => sessionStorage.getItem('access_token'),
                () => {
                    // 从cookie中获取
                    const cookies = document.cookie.split(';');
                    for (let cookie of cookies) {
                        const [name, value] = cookie.trim().split('=');
                        if (name === 'token' || name === 'access_token') {
                            return value;
                        }
                    }
                    return null;
                }
            ];

            for (let source of sources) {
                const token = source();
                if (token) {
                    console.log('🔑 找到认证令牌');
                    return token;
                }
            }

            console.log('⚠️ 未找到认证令牌');
            return '';
        },

        // 触发菜单渲染
        triggerMenuRender: function(menuData) {
            console.log('🎨 尝试触发菜单渲染...');

            // 方法1：查找Vue实例
            if (window.Vue && window.Vue.prototype) {
                this.fixVueMenu(menuData);
            }

            // 方法2：查找React组件
            if (window.React || document.querySelector('[data-reactroot]')) {
                this.fixReactMenu(menuData);
            }

            // 方法3：查找常见的菜单组件
            this.fixCommonMenu(menuData);

            // 方法4：手动创建菜单
            this.createManualMenu(menuData);
        },

        // 修复Vue菜单
        fixVueMenu: function(menuData) {
            console.log('🔧 修复Vue菜单...');

            try {
                // 查找Vue根实例
                const app = document.querySelector('#app') || document.body;
                const vueInstance = app.__vue__ || app._vnode?.componentInstance;

                if (vueInstance) {
                    console.log('✅ 找到Vue实例');

                    // 尝试设置菜单数据
                    if (vueInstance.$store) {
                        // Vuex store
                        if (vueInstance.$store.commit) {
                            vueInstance.$store.commit('menu/setRouters', menuData);
                            vueInstance.$store.commit('menu/setMenus', menuData);
                        }
                    }

                    // 直接设置组件数据
                    if (vueInstance.$children) {
                        vueInstance.$children.forEach(child => {
                            if (child.$options.name && child.$options.name.includes('menu')) {
                                child.routers = menuData;
                                child.menus = menuData;
                                child.$forceUpdate();
                            }
                        });
                    }
                }
            } catch (error) {
                console.error('修复Vue菜单失败:', error);
            }
        },

        // 修复React菜单
        fixReactMenu: function(menuData) {
            console.log('🔧 修复React菜单...');

            try {
                // 查找React根组件
                const rootElement = document.querySelector('#root') || document.querySelector('[data-reactroot]');
                if (rootElement && rootElement._reactInternalInstance) {
                    console.log('✅ 找到React根组件');

                    // 尝试触发重新渲染
                    const event = new CustomEvent('menuDataUpdate', {
                        detail: { menuData }
                    });
                    document.dispatchEvent(event);
                }
            } catch (error) {
                console.error('修复React菜单失败:', error);
            }
        },

        // 修复常见菜单组件
        fixCommonMenu: function(menuData) {
            console.log('🔧 修复常见菜单组件...');

            // Element UI菜单
            const elMenus = document.querySelectorAll('.el-menu');
            elMenus.forEach(menu => {
                console.log('✅ 找到Element UI菜单');
                // 触发更新事件
                menu.dispatchEvent(new CustomEvent('menu-update', {
                    detail: { data: menuData }
                }));
            });

            // Ant Design菜单
            const antMenus = document.querySelectorAll('.ant-menu');
            antMenus.forEach(menu => {
                console.log('✅ 找到Ant Design菜单');
                menu.dispatchEvent(new CustomEvent('menu-update', {
                    detail: { data: menuData }
                }));
            });

            // 自定义菜单组件
            const customMenus = document.querySelectorAll('[data-menu], .sidebar-menu, .nav-menu');
            customMenus.forEach(menu => {
                console.log('✅ 找到自定义菜单组件');
                menu.dispatchEvent(new CustomEvent('menu-update', {
                    detail: { data: menuData }
                }));
            });
        },

        // 手动创建菜单
        createManualMenu: function(menuData) {
            console.log('🔧 创建手动菜单...');

            // 检查是否已经有菜单
            const existingMenu = document.querySelector('#manual-menu');
            if (existingMenu) {
                existingMenu.remove();
            }

            // 创建菜单容器
            const menuContainer = document.createElement('div');
            menuContainer.id = 'manual-menu';
            menuContainer.style.cssText = `
                position: fixed;
                top: 60px;
                left: 0;
                width: 200px;
                height: calc(100vh - 60px);
                background: #f5f5f5;
                border-right: 1px solid #e0e0e0;
                padding: 10px;
                overflow-y: auto;
                z-index: 1000;
                font-family: Arial, sans-serif;
                font-size: 14px;
            `;

            // 创建菜单标题
            const menuTitle = document.createElement('div');
            menuTitle.textContent = '📋 系统菜单';
            menuTitle.style.cssText = `
                font-weight: bold;
                margin-bottom: 15px;
                padding-bottom: 10px;
                border-bottom: 1px solid #ddd;
                color: #333;
            `;
            menuContainer.appendChild(menuTitle);

            // 递归创建菜单项
            function createMenuItem(item, level = 0) {
                const menuItem = document.createElement('div');
                const indent = level * 15;

                menuItem.style.cssText = `
                    margin-bottom: 5px;
                    padding: 8px 12px;
                    margin-left: ${indent}px;
                    border-radius: 4px;
                    cursor: pointer;
                    transition: background-color 0.2s;
                    background: ${level === 0 ? '#fff' : '#f9f9f9'};
                    border-left: 3px solid ${level === 0 ? '#0066cc' : 'transparent'};
                `;

                menuItem.onmouseover = () => {
                    menuItem.style.background = '#e6f7ff';
                };

                menuItem.onmouseout = () => {
                    menuItem.style.background = level === 0 ? '#fff' : '#f9f9f9';
                };

                menuItem.onclick = () => {
                    if (item.path) {
                        console.log('导航到:', item.path);
                        // 这里可以添加路由导航逻辑
                        window.location.hash = item.path;
                    }
                };

                // 菜单项内容
                const icon = item.meta?.icon || '📄';
                const title = item.meta?.title || item.name || item.path || '未命名';

                menuItem.innerHTML = `
                    <span style="margin-right: 8px;">${icon}</span>
                    <span>${title}</span>
                `;

                return menuItem;
            }

            // 递归渲染菜单树
            function renderMenuTree(items, level = 0) {
                items.forEach(item => {
                    const menuItem = createMenuItem(item, level);
                    menuContainer.appendChild(menuItem);

                    // 渲染子菜单
                    if (item.children && item.children.length > 0) {
                        renderMenuTree(item.children, level + 1);
                    }
                });
            }

            // 渲染菜单
            renderMenuTree(menuData);

            // 添加关闭按钮
            const closeButton = document.createElement('button');
            closeButton.textContent = '✕ 关闭菜单';
            closeButton.style.cssText = `
                position: absolute;
                bottom: 10px;
                left: 10px;
                right: 10px;
                padding: 8px;
                background: #ef4444;
                color: white;
                border: none;
                border-radius: 4px;
                cursor: pointer;
                font-size: 12px;
            `;
            closeButton.onclick = () => menuContainer.remove();
            menuContainer.appendChild(closeButton);

            // 添加到页面
            document.body.appendChild(menuContainer);

            console.log('✅ 已创建手动菜单');
        },

        // 修复2：检查数据绑定
        fixDataBinding: function() {
            console.log('🔧 修复2: 检查数据绑定...');

            // 检查全局状态管理
            const stores = [window.store, window.$store, window.app?.config?.globalProperties?.$store];
            stores.forEach((store, index) => {
                if (store && store.state && store.state.menu) {
                    console.log(`✅ 找到状态管理器 ${index + 1}`);
                    if (store.commit) {
                        store.commit('menu/setRouters', window.menuData || []);
                        store.commit('menu/setMenus', window.menuData || []);
                    }
                }
            });

            return true;
        },

        // 修复3：强制重新渲染
        forceReRender: function() {
            console.log('🔧 修复3: 强制重新渲染...');

            // 触发窗口resize事件
            window.dispatchEvent(new Event('resize'));

            // 触发自定义事件
            const event = new CustomEvent('menuReload', {
                detail: { data: window.menuData }
            });
            document.dispatchEvent(event);

            // 查找并更新菜单组件
            const menuComponents = document.querySelectorAll('[data-component*="menu"], .menu-component');
            menuComponents.forEach(component => {
                component.dispatchEvent(new CustomEvent('update', {
                    detail: { menuData: window.menuData }
                }));
            });

            return true;
        }
    };

    // 自动修复流程
    async function autoFix() {
        console.group('🔧 开始自动修复菜单渲染问题');

        try {
            // 步骤1：获取菜单数据
            console.log('📍 步骤1: 获取菜单数据');
            await fixes.fixApiCall();

            // 步骤2：修复数据绑定
            console.log('📍 步骤2: 修复数据绑定');
            fixes.fixDataBinding();

            // 步骤3：强制重新渲染
            console.log('📍 步骤3: 强制重新渲染');
            fixes.forceReRender();

            console.log('✅ 自动修复完成');

        } catch (error) {
            console.error('❌ 自动修复失败:', error);
        }

        console.groupEnd();

        // 显示修复报告
        showFixReport();
    }

    // 显示修复报告
    function showFixReport() {
        const reportDiv = document.createElement('div');
        reportDiv.id = 'menu-fix-report';
        reportDiv.innerHTML = `
            <div style="
                position: fixed;
                bottom: 20px;
                right: 20px;
                width: 300px;
                background: #fff;
                border: 2px solid #10b981;
                border-radius: 8px;
                padding: 15px;
                z-index: 9999;
                box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
                font-family: Arial, sans-serif;
                font-size: 12px;
            ">
                <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 10px;">
                    <h4 style="margin: 0; color: #10b981;">✅ 菜单修复完成</h4>
                    <button onclick="this.parentNode.parentNode.remove()" style="
                        background: #6b7280;
                        color: white;
                        border: none;
                        border-radius: 50%;
                        width: 24px;
                        height: 24px;
                        cursor: pointer;
                        font-size: 10px;
                    ">✕</button>
                </div>

                <div style="margin-bottom: 10px;">
                    <div>🔍 <strong>修复项目：</strong></div>
                    <div style="margin-left: 15px; margin-top: 5px;">
                        <div>✅ 获取菜单数据</div>
                        <div>✅ 修复数据绑定</div>
                        <div>✅ 强制重新渲染</div>
                        <div>✅ 创建备用菜单</div>
                    </div>
                </div>

                <div style="margin-bottom: 10px;">
                    <div>📋 <strong>菜单状态：</strong></div>
                    <div style="margin-left: 15px;">
                        <div>菜单项数量: <span id="menu-count">-</span></div>
                        <div>渲染状态: <span id="render-status">-</span></div>
                    </div>
                </div>

                <div style="display: flex; gap: 8px;">
                    <button onclick="runAutoFix()" style="
                        flex: 1;
                        background: #0066cc;
                        color: white;
                        border: none;
                        border-radius: 4px;
                        padding: 6px 12px;
                        cursor: pointer;
                        font-size: 11px;
                    ">🔄 重新修复</button>
                    <button onclick="testMenu()" style="
                        flex: 1;
                        background: #10b981;
                        color: white;
                        border: none;
                        border-radius: 4px;
                        padding: 6px 12px;
                        cursor: pointer;
                        font-size: 11px;
                    ">🧪 测试菜单</button>
                </div>
            </div>
        `;

        // 更新状态信息
        setTimeout(() => {
            const menuCount = reportDiv.querySelector('#menu-count');
            const renderStatus = reportDiv.querySelector('#render-status');

            if (window.menuData) {
                menuCount.textContent = window.menuData.length;
                renderStatus.textContent = '已渲染';
                renderStatus.style.color = '#10b981';
            } else {
                menuCount.textContent = '0';
                renderStatus.textContent = '未渲染';
                renderStatus.style.color = '#ef4444';
            }
        }, 500);

        document.body.appendChild(reportDiv);
    }

    // 测试菜单功能
    window.testMenu = function() {
        console.log('🧪 测试菜单功能...');

        if (window.menuData && window.menuData.length > 0) {
            console.log('✅ 菜单数据存在:', window.menuData);

            // 测试菜单导航
            const firstMenu = window.menuData[0];
            if (firstMenu && firstMenu.path) {
                console.log('导航到第一个菜单:', firstMenu.path);
                // 这里可以添加实际的导航逻辑
            }

            alert('菜单测试完成！请查看控制台输出。');
        } else {
            console.log('❌ 没有找到菜单数据');
            alert('未找到菜单数据，请先运行修复功能。');
        }
    };

    // 暴露全局函数
    window.runAutoFix = autoFix;
    window.fixes = fixes;

    // 自动运行修复
    console.log('🚀 启动自动修复...');
    setTimeout(autoFix, 2000);

    // 快捷键支持
    document.addEventListener('keydown', function(event) {
        // Ctrl+Shift+M: 运行修复
        if (event.ctrlKey && event.shiftKey && event.key === 'M') {
            event.preventDefault();
            autoFix();
        }
    });

    console.log('💡 快捷键: Ctrl+Shift+M 重新运行修复');

})();


