/**
 * Ozon API 请求工具 - 内容脚本
 * 功能：在页面上下文中执行请求，提取 cookies，与页面交互
 */

// 全局变量
let isToolEnabled = false;
let debugMode = false;

// 初始化内容脚本
(function() {
    'use strict';
    
    console.log('🛒 Ozon API 工具内容脚本已加载');
    
    // 监听来自 popup 的消息
    chrome.runtime.onMessage.addListener((request, sender, sendResponse) => {
        handleMessage(request, sender, sendResponse);
        return true; // 保持消息通道开放以支持异步响应
    });

    // 页面加载完成后的初始化
    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', initialize);
    } else {
        initialize();
    }
})();

/**
 * 处理来自插件的消息
 */
async function handleMessage(request, sender, sendResponse) {
    try {
        switch (request.action) {
            case 'makeRequest':
                const result = await makeApiRequest(request.url, request.config, request.cookies);
                sendResponse(result);
                break;
                
            case 'extractCookies':
                const cookies = await extractPageCookies();
                sendResponse(cookies);
                break;
                
            case 'getPageInfo':
                const pageInfo = getPageInformation();
                sendResponse(pageInfo);
                break;
                
            case 'injectTool':
                await injectApiTool();
                sendResponse({ success: true });
                break;
                
            case 'toggleDebug':
                debugMode = !debugMode;
                sendResponse({ success: true, debugMode });
                break;
                
            default:
                sendResponse({ success: false, error: 'Unknown action' });
        }
    } catch (error) {
        console.error('处理消息时出错:', error);
        sendResponse({ success: false, error: error.message });
    }
}

/**
 * 发送 API 请求
 */
async function makeApiRequest(url, config, storedCookies) {
    try {
        logDebug('开始发送 API 请求', { url, config });

        // 如果有存储的 cookies，尝试设置它们
        if (storedCookies && storedCookies.length > 0) {
            await applyStoredCookies(storedCookies);
        }

        // 添加请求拦截和日志
        const startTime = Date.now();
        
        const response = await fetch(url, {
            ...config,
            // 确保包含凭据
            credentials: 'include'
        });
        
        const endTime = Date.now();
        const duration = endTime - startTime;

        logDebug('API 请求完成', { 
            status: response.status, 
            duration: `${duration}ms`,
            headers: Object.fromEntries(response.headers.entries())
        });

        // 处理响应数据
        let data;
        const contentType = response.headers.get('content-type');
        
        if (contentType && contentType.includes('application/json')) {
            data = await response.json();
        } else {
            data = await response.text();
        }

        // 记录请求信息
        await logApiRequest({
            url,
            method: config.method,
            status: response.status,
            duration,
            timestamp: Date.now()
        });

        return {
            success: true,
            status: response.status,
            statusText: response.statusText,
            data: data,
            headers: Object.fromEntries(response.headers.entries()),
            duration: duration
        };

    } catch (error) {
        logDebug('API 请求失败', { error: error.message });
        
        // 特殊错误处理
        if (error.name === 'TypeError' && error.message.includes('fetch')) {
            return {
                success: false,
                error: '网络错误：请检查网络连接或 URL 是否正确'
            };
        }
        
        if (error.message.includes('CORS')) {
            return {
                success: false,
                error: 'CORS 错误：请确保在正确的域名下使用插件'
            };
        }

        return {
            success: false,
            error: `请求失败: ${error.message}`
        };
    }
}

/**
 * 提取页面 cookies
 */
async function extractPageCookies() {
    try {
        logDebug('开始提取页面 cookies');

        // 方法1：从 document.cookie 提取
        const documentCookies = extractDocumentCookies();
        
        // 方法2：尝试从存储中获取（如果有权限）
        let storageCookies = [];
        try {
            // 这里可以添加从浏览器存储中读取 cookies 的逻辑
            // 需要相应的权限
        } catch (e) {
            logDebug('无法从存储中读取 cookies:', e.message);
        }

        // 合并 cookies
        const allCookies = [...documentCookies, ...storageCookies];
        
        // 去重
        const uniqueCookies = Array.from(
            new Map(allCookies.map(cookie => [cookie.name, cookie])).values()
        );

        logDebug(`成功提取 ${uniqueCookies.length} 个 cookies`);

        return {
            success: true,
            cookies: uniqueCookies,
            source: 'document.cookie',
            timestamp: Date.now()
        };

    } catch (error) {
        logDebug('提取 cookies 失败:', error.message);
        return {
            success: false,
            error: `提取 cookies 失败: ${error.message}`
        };
    }
}

/**
 * 从 document.cookie 提取 cookies
 */
function extractDocumentCookies() {
    const cookieString = document.cookie;
    if (!cookieString) return [];

    return cookieString.split(';').map(cookie => {
        const [name, ...valueParts] = cookie.trim().split('=');
        return {
            name: name.trim(),
            value: valueParts.join('=').trim(),
            domain: window.location.hostname,
            path: '/',
            source: 'document'
        };
    }).filter(cookie => cookie.name && cookie.value);
}

/**
 * 应用存储的 cookies（尝试设置到当前域）
 */
async function applyStoredCookies(cookies) {
    try {
        // 注意：由于浏览器安全限制，不能直接设置所有 cookies
        // 但可以尝试设置一些非 httpOnly 的 cookies
        
        const applicableCookies = cookies.filter(cookie => {
            // 过滤出可以设置的 cookies
            return !cookie.name.startsWith('__Secure-') || 
                   window.location.protocol === 'https:';
        });

        for (const cookie of applicableCookies) {
            try {
                // 尝试设置 cookie
                document.cookie = `${cookie.name}=${cookie.value}; path=/; domain=${window.location.hostname}`;
            } catch (e) {
                logDebug(`无法设置 cookie ${cookie.name}:`, e.message);
            }
        }

        logDebug(`尝试应用 ${applicableCookies.length} 个 cookies`);

    } catch (error) {
        logDebug('应用 cookies 失败:', error.message);
    }
}

/**
 * 获取页面信息
 */
function getPageInformation() {
    return {
        success: true,
        pageInfo: {
            url: window.location.href,
            hostname: window.location.hostname,
            pathname: window.location.pathname,
            title: document.title,
            userAgent: navigator.userAgent,
            cookies: extractDocumentCookies(),
            localStorage: getLocalStorageInfo(),
            sessionStorage: getSessionStorageInfo(),
            isOzonSeller: window.location.hostname.includes('seller.ozon.ru'),
            timestamp: Date.now()
        }
    };
}

/**
 * 获取 localStorage 信息
 */
function getLocalStorageInfo() {
    try {
        const items = {};
        for (let i = 0; i < localStorage.length; i++) {
            const key = localStorage.key(i);
            const value = localStorage.getItem(key);
            // 只保存前 100 个字符，避免数据过大
            items[key] = value ? value.substring(0, 100) : value;
        }
        return items;
    } catch (e) {
        return { error: e.message };
    }
}

/**
 * 获取 sessionStorage 信息
 */
function getSessionStorageInfo() {
    try {
        const items = {};
        for (let i = 0; i < sessionStorage.length; i++) {
            const key = sessionStorage.key(i);
            const value = sessionStorage.getItem(key);
            // 只保存前 100 个字符，避免数据过大
            items[key] = value ? value.substring(0, 100) : value;
        }
        return items;
    } catch (e) {
        return { error: e.message };
    }
}

/**
 * 注入 API 工具到页面
 */
async function injectApiTool() {
    if (isToolEnabled) return;

    try {
        // 创建浮动工具栏
        const toolbar = createFloatingToolbar();
        document.body.appendChild(toolbar);
        
        isToolEnabled = true;
        logDebug('API 工具已注入到页面');

    } catch (error) {
        logDebug('注入 API 工具失败:', error.message);
        throw error;
    }
}

/**
 * 创建浮动工具栏
 */
function createFloatingToolbar() {
    const toolbar = document.createElement('div');
    toolbar.id = 'ozon-api-toolbar';
    toolbar.innerHTML = `
        <div style="
            position: fixed;
            top: 20px;
            right: 20px;
            z-index: 10000;
            background: linear-gradient(135deg, #005bb5, #0066cc);
            color: white;
            padding: 10px 15px;
            border-radius: 8px;
            box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
            font-family: Arial, sans-serif;
            font-size: 12px;
            cursor: move;
            user-select: none;
        ">
            <div style="display: flex; align-items: center; gap: 10px;">
                <span>🛒 Ozon API 工具</span>
                <button id="extract-cookies-btn" style="
                    background: #10b981;
                    color: white;
                    border: none;
                    padding: 4px 8px;
                    border-radius: 4px;
                    font-size: 11px;
                    cursor: pointer;
                ">提取 Cookies</button>
                <button id="close-toolbar-btn" style="
                    background: #ef4444;
                    color: white;
                    border: none;
                    padding: 2px 6px;
                    border-radius: 3px;
                    font-size: 10px;
                    cursor: pointer;
                ">×</button>
            </div>
        </div>
    `;

    // 添加事件监听器
    const extractBtn = toolbar.querySelector('#extract-cookies-btn');
    const closeBtn = toolbar.querySelector('#close-toolbar-btn');

    extractBtn.addEventListener('click', async () => {
        const result = await extractPageCookies();
        if (result.success) {
            // 保存到插件存储
            chrome.runtime.sendMessage({
                action: 'saveCookies',
                cookies: result.cookies
            });
            
            showNotification(`✅ 已提取 ${result.cookies.length} 个 cookies`, 'success');
        } else {
            showNotification(`❌ 提取失败: ${result.error}`, 'error');
        }
    });

    closeBtn.addEventListener('click', () => {
        toolbar.remove();
        isToolEnabled = false;
    });

    // 添加拖拽功能
    makeDraggable(toolbar);

    return toolbar;
}

/**
 * 使元素可拖拽
 */
function makeDraggable(element) {
    let pos1 = 0, pos2 = 0, pos3 = 0, pos4 = 0;
    
    element.onmousedown = dragMouseDown;

    function dragMouseDown(e) {
        e = e || window.event;
        e.preventDefault();
        pos3 = e.clientX;
        pos4 = e.clientY;
        document.onmouseup = closeDragElement;
        document.onmousemove = elementDrag;
    }

    function elementDrag(e) {
        e = e || window.event;
        e.preventDefault();
        pos1 = pos3 - e.clientX;
        pos2 = pos4 - e.clientY;
        pos3 = e.clientX;
        pos4 = e.clientY;
        element.style.top = (element.offsetTop - pos2) + "px";
        element.style.left = (element.offsetLeft - pos1) + "px";
        element.style.right = 'auto';
    }

    function closeDragElement() {
        document.onmouseup = null;
        document.onmousemove = null;
    }
}

/**
 * 显示通知
 */
function showNotification(message, type = 'info') {
    const notification = document.createElement('div');
    notification.innerHTML = `
        <div style="
            position: fixed;
            top: 80px;
            right: 20px;
            z-index: 10001;
            background: ${type === 'success' ? '#10b981' : type === 'error' ? '#ef4444' : '#0066cc'};
            color: white;
            padding: 12px 16px;
            border-radius: 6px;
            box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
            font-family: Arial, sans-serif;
            font-size: 13px;
            max-width: 300px;
            animation: slideIn 0.3s ease-out;
        ">
            ${message}
        </div>
    `;

    // 添加动画样式
    if (!document.getElementById('notification-styles')) {
        const styles = document.createElement('style');
        styles.id = 'notification-styles';
        styles.textContent = `
            @keyframes slideIn {
                from { transform: translateX(100%); opacity: 0; }
                to { transform: translateX(0); opacity: 1; }
            }
        `;
        document.head.appendChild(styles);
    }

    document.body.appendChild(notification);

    // 3秒后自动移除
    setTimeout(() => {
        notification.remove();
    }, 3000);
}

/**
 * 记录 API 请求信息
 */
async function logApiRequest(requestInfo) {
    try {
        // 发送到后台脚本保存
        chrome.runtime.sendMessage({
            action: 'saveApiLog',
            data: requestInfo
        });
    } catch (e) {
        logDebug('记录 API 请求失败:', e.message);
    }
}

/**
 * 调试日志
 */
function logDebug(message, data = null) {
    if (debugMode) {
        console.log(`[Ozon API Tool] ${message}`, data || '');
    }
}

/**
 * 初始化函数
 */
function initialize() {
    logDebug('内容脚本初始化完成');
    
    // 检查是否是 Ozon 页面
    if (window.location.hostname.includes('ozon.ru')) {
        logDebug('检测到 Ozon 页面，工具可用');
        
        // 自动提取一次 cookies（如果用户已同意）
        setTimeout(() => {
            extractPageCookies().then(result => {
                if (result.success && result.cookies.length > 0) {
                    logDebug(`自动提取了 ${result.cookies.length} 个 cookies`);
                }
            });
        }, 2000);
    }

    // 监听页面变化（SPA 应用）
    observePageChanges();
}

/**
 * 监听页面变化
 */
function observePageChanges() {
    const observer = new MutationObserver((mutations) => {
        mutations.forEach((mutation) => {
            if (mutation.type === 'childList' && mutation.addedNodes.length > 0) {
                // 页面内容发生变化，可以在这里添加相应的处理逻辑
                logDebug('页面内容发生变化');
            }
        });
    });

    observer.observe(document.body, {
        childList: true,
        subtree: true
    });
}

// 页面卸载时的清理
window.addEventListener('beforeunload', () => {
    logDebug('页面即将卸载，清理内容脚本');
});


