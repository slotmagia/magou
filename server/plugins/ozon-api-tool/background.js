/**
 * Ozon API 请求工具 - 后台脚本 (Service Worker)
 * 功能：处理跨域请求、管理存储、处理通知、维护状态
 */

// 全局变量
const API_LOGS_KEY = 'ozon_api_logs';
const MAX_LOGS = 1000;
let isDebugMode = false;

// 插件安装/更新时的初始化
chrome.runtime.onInstalled.addListener((details) => {
    console.log('🛒 Ozon API 工具后台脚本已安装');
    
    if (details.reason === 'install') {
        // 首次安装
        handleFirstInstall();
    } else if (details.reason === 'update') {
        // 更新
        handleUpdate(details.previousVersion);
    }
});

// 插件启动时的初始化
chrome.runtime.onStartup.addListener(() => {
    console.log('🛒 Ozon API 工具后台脚本已启动');
    initializeBackground();
});

// 处理首次安装
async function handleFirstInstall() {
    try {
        // 设置默认配置
        await chrome.storage.local.set({
            'ozon_tool_config': {
                version: '1.0.0',
                firstInstall: Date.now(),
                debugMode: false,
                autoExtractCookies: true,
                maxHistoryItems: 50
            }
        });

        // 显示欢迎通知
        chrome.notifications.create({
            type: 'basic',
            iconUrl: 'icons/icon48.png',
            title: 'Ozon API 工具',
            message: '安装成功！点击插件图标开始使用。'
        });

        logDebug('首次安装完成');
    } catch (error) {
        console.error('处理首次安装失败:', error);
    }
}

// 处理更新
async function handleUpdate(previousVersion) {
    try {
        logDebug(`从版本 ${previousVersion} 更新到 1.0.0`);
        
        // 这里可以添加数据迁移逻辑
        // migrateData(previousVersion);

    } catch (error) {
        console.error('处理更新失败:', error);
    }
}

// 初始化后台脚本
async function initializeBackground() {
    try {
        // 加载配置
        const config = await getStorageData('ozon_tool_config');
        if (config) {
            isDebugMode = config.debugMode || false;
        }

        // 清理过期数据
        await cleanupOldData();

        logDebug('后台脚本初始化完成');
    } catch (error) {
        console.error('初始化后台脚本失败:', error);
    }
}

// 监听来自其他脚本的消息
chrome.runtime.onMessage.addListener((request, sender, sendResponse) => {
    handleBackgroundMessage(request, sender, sendResponse);
    return true; // 保持消息通道开放
});

// 处理后台消息
async function handleBackgroundMessage(request, sender, sendResponse) {
    try {
        switch (request.action) {
            case 'makeApiRequest':
                const result = await makeBackgroundApiRequest(request.url, request.config);
                sendResponse(result);
                break;

            case 'saveCookies':
                await saveCookiesToStorage(request.cookies);
                sendResponse({ success: true });
                break;

            case 'saveApiLog':
                await saveApiLog(request.data);
                sendResponse({ success: true });
                break;

            case 'getConfig':
                const config = await getToolConfig();
                sendResponse({ success: true, config });
                break;

            case 'updateConfig':
                await updateToolConfig(request.config);
                sendResponse({ success: true });
                break;

            case 'exportData':
                const exportData = await exportAllData();
                sendResponse({ success: true, data: exportData });
                break;

            case 'importData':
                await importAllData(request.data);
                sendResponse({ success: true });
                break;

            case 'clearAllData':
                await clearAllData();
                sendResponse({ success: true });
                break;

            case 'toggleDebug':
                isDebugMode = !isDebugMode;
                await updateToolConfig({ debugMode: isDebugMode });
                sendResponse({ success: true, debugMode: isDebugMode });
                break;

            default:
                sendResponse({ success: false, error: 'Unknown action' });
        }
    } catch (error) {
        console.error('处理后台消息失败:', error);
        sendResponse({ success: false, error: error.message });
    }
}

// 在后台发送 API 请求（用于跨域场景）
async function makeBackgroundApiRequest(url, config) {
    try {
        logDebug('后台发送 API 请求', { url, method: config.method });

        const response = await fetch(url, {
            ...config,
            // 后台脚本中不需要 credentials: 'include'
            credentials: 'omit'
        });

        let data;
        const contentType = response.headers.get('content-type');
        
        if (contentType && contentType.includes('application/json')) {
            data = await response.json();
        } else {
            data = await response.text();
        }

        logDebug('后台 API 请求成功', { status: response.status });

        return {
            success: true,
            status: response.status,
            statusText: response.statusText,
            data: data,
            headers: Object.fromEntries(response.headers.entries())
        };

    } catch (error) {
        logDebug('后台 API 请求失败', { error: error.message });
        return {
            success: false,
            error: `后台请求失败: ${error.message}`
        };
    }
}

// 保存 cookies 到存储
async function saveCookiesToStorage(cookies) {
    try {
        await chrome.storage.local.set({
            'ozonCookies': cookies,
            'cookiesUpdated': Date.now()
        });

        logDebug(`保存了 ${cookies.length} 个 cookies`);

        // 发送通知给所有相关标签页
        notifyTabsOfCookieUpdate(cookies.length);

    } catch (error) {
        console.error('保存 cookies 失败:', error);
        throw error;
    }
}

// 通知标签页 cookie 更新
async function notifyTabsOfCookieUpdate(cookieCount) {
    try {
        const tabs = await chrome.tabs.query({ url: 'https://*.ozon.ru/*' });
        
        for (const tab of tabs) {
            try {
                await chrome.tabs.sendMessage(tab.id, {
                    action: 'cookiesUpdated',
                    count: cookieCount
                });
            } catch (e) {
                // 某些标签页可能无法接收消息，忽略错误
            }
        }
    } catch (error) {
        logDebug('通知标签页失败:', error.message);
    }
}

// 保存 API 日志
async function saveApiLog(logData) {
    try {
        const logs = await getStorageData(API_LOGS_KEY) || [];
        
        // 添加新日志
        logs.unshift({
            ...logData,
            id: Date.now() + Math.random().toString(36).substr(2, 9),
            timestamp: Date.now()
        });

        // 限制日志数量
        if (logs.length > MAX_LOGS) {
            logs.splice(MAX_LOGS);
        }

        await chrome.storage.local.set({ [API_LOGS_KEY]: logs });
        logDebug('API 日志已保存');

    } catch (error) {
        console.error('保存 API 日志失败:', error);
    }
}

// 获取工具配置
async function getToolConfig() {
    try {
        const config = await getStorageData('ozon_tool_config') || {};
        return {
            version: '1.0.0',
            debugMode: false,
            autoExtractCookies: true,
            maxHistoryItems: 50,
            ...config
        };
    } catch (error) {
        console.error('获取工具配置失败:', error);
        return {};
    }
}

// 更新工具配置
async function updateToolConfig(newConfig) {
    try {
        const currentConfig = await getToolConfig();
        const updatedConfig = { ...currentConfig, ...newConfig };
        
        await chrome.storage.local.set({ 'ozon_tool_config': updatedConfig });
        
        // 更新全局变量
        if (updatedConfig.debugMode !== undefined) {
            isDebugMode = updatedConfig.debugMode;
        }

        logDebug('工具配置已更新', updatedConfig);
    } catch (error) {
        console.error('更新工具配置失败:', error);
        throw error;
    }
}

// 导出所有数据
async function exportAllData() {
    try {
        const allData = await chrome.storage.local.get(null);
        
        const exportData = {
            version: '1.0.0',
            exportTime: Date.now(),
            data: allData
        };

        logDebug('数据导出完成', { size: Object.keys(allData).length });
        return exportData;

    } catch (error) {
        console.error('导出数据失败:', error);
        throw error;
    }
}

// 导入所有数据
async function importAllData(importData) {
    try {
        if (!importData || !importData.data) {
            throw new Error('导入数据格式错误');
        }

        // 清除现有数据
        await chrome.storage.local.clear();
        
        // 导入新数据
        await chrome.storage.local.set(importData.data);

        logDebug('数据导入完成');

    } catch (error) {
        console.error('导入数据失败:', error);
        throw error;
    }
}

// 清除所有数据
async function clearAllData() {
    try {
        await chrome.storage.local.clear();
        
        // 重新设置默认配置
        await handleFirstInstall();
        
        logDebug('所有数据已清除');

    } catch (error) {
        console.error('清除数据失败:', error);
        throw error;
    }
}

// 清理过期数据
async function cleanupOldData() {
    try {
        const oneWeekAgo = Date.now() - (7 * 24 * 60 * 60 * 1000);
        
        // 清理过期的 API 日志
        const logs = await getStorageData(API_LOGS_KEY) || [];
        const filteredLogs = logs.filter(log => log.timestamp > oneWeekAgo);
        
        if (filteredLogs.length < logs.length) {
            await chrome.storage.local.set({ [API_LOGS_KEY]: filteredLogs });
            logDebug(`清理了 ${logs.length - filteredLogs.length} 条过期日志`);
        }

        // 清理过期的请求历史
        const history = await getStorageData('requestHistory') || [];
        const filteredHistory = history.filter(item => item.timestamp > oneWeekAgo);
        
        if (filteredHistory.length < history.length) {
            await chrome.storage.local.set({ requestHistory: filteredHistory });
            logDebug(`清理了 ${history.length - filteredHistory.length} 条过期历史`);
        }

    } catch (error) {
        console.error('清理过期数据失败:', error);
    }
}

// 监听标签页更新
chrome.tabs.onUpdated.addListener(async (tabId, changeInfo, tab) => {
    // 当标签页加载完成且是 Ozon 页面时
    if (changeInfo.status === 'complete' && tab.url && tab.url.includes('ozon.ru')) {
        try {
            // 获取配置
            const config = await getToolConfig();
            
            // 如果启用了自动提取 cookies
            if (config.autoExtractCookies) {
                setTimeout(async () => {
                    try {
                        await chrome.tabs.sendMessage(tabId, {
                            action: 'autoExtractCookies'
                        });
                    } catch (e) {
                        // 忽略无法发送消息的错误
                    }
                }, 3000);
            }
        } catch (error) {
            logDebug('处理标签页更新失败:', error.message);
        }
    }
});

// 监听通知点击
chrome.notifications.onClicked.addListener((notificationId) => {
    // 打开插件弹窗或相关页面
    chrome.action.openPopup();
});

// 工具函数

// 从存储获取数据
async function getStorageData(key) {
    try {
        const result = await chrome.storage.local.get(key);
        return result[key];
    } catch (error) {
        console.error(`获取存储数据失败 (${key}):`, error);
        return null;
    }
}

// 调试日志
function logDebug(message, data = null) {
    if (isDebugMode) {
        console.log(`[Ozon API Tool Background] ${message}`, data || '');
    }
}

// 格式化数据大小
function formatDataSize(bytes) {
    if (bytes === 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
}

// 生成唯一 ID
function generateId() {
    return Date.now().toString(36) + Math.random().toString(36).substr(2, 9);
}

// 错误处理和报告
chrome.runtime.onSuspend.addListener(() => {
    console.log('🛒 Ozon API 工具后台脚本即将挂起');
});

// 监听存储变化
chrome.storage.onChanged.addListener((changes, namespace) => {
    if (namespace === 'local') {
        for (const key in changes) {
            logDebug(`存储变化: ${key}`, {
                oldValue: changes[key].oldValue ? '已设置' : '未设置',
                newValue: changes[key].newValue ? '已设置' : '未设置'
            });
        }
    }
});

// 初始化
initializeBackground();


