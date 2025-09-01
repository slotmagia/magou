/**
 * Ozon API 请求工具 - 弹窗脚本
 * 功能：处理用户界面交互，发送 API 请求，管理 cookies
 */

document.addEventListener('DOMContentLoaded', function() {
    // 获取 DOM 元素
    const elements = {
        // 标签页
        tabs: document.querySelectorAll('.tab'),
        tabContents: document.querySelectorAll('.tab-content'),
        
        // 请求配置
        apiTypeSelect: document.getElementById('apiType'),
        urlInput: document.getElementById('url'),
        urlPreview: document.getElementById('url-preview'),
        methodSelect: document.getElementById('method'),
        companyIdInput: document.getElementById('companyId'),
        customHeadersTextarea: document.getElementById('customHeaders'),
        requestBodyTextarea: document.getElementById('requestBody'),
        bodyGroup: document.getElementById('body-group'),
        
        // 按钮
        sendButton: document.getElementById('sendRequest'),
        clearFormButton: document.getElementById('clearForm'),
        extractCookiesButton: document.getElementById('extractCookies'),
        clearCookiesButton: document.getElementById('clearCookies'),
        parseCookiesButton: document.getElementById('parseCookies'),
        clearHistoryButton: document.getElementById('clearHistory'),
        exportHistoryButton: document.getElementById('exportHistory'),
        
        // 状态显示
        connectionStatus: document.getElementById('connection-status'),
        statusText: document.getElementById('status-text'),
        cookieCount: document.getElementById('cookie-count'),
        resultDiv: document.getElementById('result'),
        
        // Cookie 管理
        cookieList: document.getElementById('cookieList'),
        cookieString: document.getElementById('cookieString'),
        
        // 历史记录
        historyList: document.getElementById('historyList'),
        
        // 加载状态
        sendText: document.getElementById('send-text'),
        loadingSpinner: document.getElementById('loading-spinner'),
        
        // JSON 验证
        headersValidation: document.getElementById('headers-validation'),
        bodyValidation: document.getElementById('body-validation')
    };

    // API 类型预设配置
    const apiPresets = {
        'finance-info': {
            url: 'https://seller.ozon.ru/api/company/finance-info',
            method: 'POST',
            body: '{"marketplaceSellerId":"2361110"}',
            description: '获取公司财务信息'
        },
        'get-language': {
            url: 'https://seller.ozon.ru/api/user/get-language',
            method: 'POST',
            body: '{"company_id":2361110}',
            description: '获取用户语言设置'
        },
        'company-info': {
            url: 'https://seller.ozon.ru/api/company/info',
            method: 'GET',
            body: '',
            description: '获取公司基本信息'
        },
        'user-profile': {
            url: 'https://seller.ozon.ru/api/user/profile',
            method: 'GET',
            body: '',
            description: '获取用户配置信息'
        }
    };

    // 初始化应用
    async function init() {
        setupEventListeners();
        await loadStoredData();
        updateConnectionStatus();
        setupJSONValidation();
        updateMethodVisibility();
        updateUrlPreview();
    }

    // 设置事件监听器
    function setupEventListeners() {
        // 标签页切换
        elements.tabs.forEach(tab => {
            tab.addEventListener('click', () => switchTab(tab.dataset.tab));
        });

        // API 类型选择
        elements.apiTypeSelect.addEventListener('change', handleApiTypeChange);

        // 请求方法变化
        elements.methodSelect.addEventListener('change', updateMethodVisibility);

        // URL 预览更新
        elements.urlInput.addEventListener('input', updateUrlPreview);

        // 按钮事件
        elements.sendButton.addEventListener('click', handleSendRequest);
        elements.clearFormButton.addEventListener('click', handleClearForm);
        elements.extractCookiesButton.addEventListener('click', handleExtractCookies);
        elements.clearCookiesButton.addEventListener('click', handleClearCookies);
        elements.parseCookiesButton.addEventListener('click', handleParseCookies);
        elements.clearHistoryButton.addEventListener('click', handleClearHistory);
        elements.exportHistoryButton.addEventListener('click', handleExportHistory);

        // 自动保存输入
        elements.companyIdInput.addEventListener('input', () => {
            saveToStorage('companyId', elements.companyIdInput.value);
            updateApiPresets();
        });

        elements.customHeadersTextarea.addEventListener('input', () => {
            saveToStorage('customHeaders', elements.customHeadersTextarea.value);
        });

        elements.urlInput.addEventListener('input', () => {
            saveToStorage('lastUrl', elements.urlInput.value);
        });

        elements.requestBodyTextarea.addEventListener('input', () => {
            saveToStorage('lastBody', elements.requestBodyTextarea.value);
        });
    }

    // 设置 JSON 验证
    function setupJSONValidation() {
        elements.customHeadersTextarea.addEventListener('input', () => {
            validateJSON(elements.customHeadersTextarea, elements.headersValidation);
        });

        elements.requestBodyTextarea.addEventListener('input', () => {
            validateJSON(elements.requestBodyTextarea, elements.bodyValidation);
        });
    }

    // JSON 验证函数
    function validateJSON(textarea, validationElement) {
        const value = textarea.value.trim();
        if (!value) {
            textarea.classList.remove('json-valid', 'json-invalid');
            validationElement.style.display = 'none';
            return true;
        }

        try {
            JSON.parse(value);
            textarea.classList.remove('json-invalid');
            textarea.classList.add('json-valid');
            validationElement.textContent = '✓';
            validationElement.className = 'validation-message validation-valid';
            validationElement.style.display = 'block';
            return true;
        } catch (e) {
            textarea.classList.remove('json-valid');
            textarea.classList.add('json-invalid');
            validationElement.textContent = '✗';
            validationElement.className = 'validation-message validation-invalid';
            validationElement.style.display = 'block';
            return false;
        }
    }

    // 切换标签页
    let currentTab = 'request';
    function switchTab(tabName) {
        // 防止重复切换到同一标签页
        if (currentTab === tabName) {
            return;
        }
        
        currentTab = tabName;
        
        elements.tabs.forEach(tab => {
            tab.classList.toggle('active', tab.dataset.tab === tabName);
        });

        elements.tabContents.forEach(content => {
            content.classList.toggle('active', content.id === `tab-${tabName}`);
        });

        // 使用 requestAnimationFrame 延迟异步操作，避免阻塞DOM更新
        requestAnimationFrame(() => {
            if (tabName === 'cookies') {
                loadCookieList();
            } else if (tabName === 'history') {
                loadHistoryList();
            }
        });
    }

    // 处理 API 类型变化
    function handleApiTypeChange() {
        const selectedType = elements.apiTypeSelect.value;
        if (apiPresets[selectedType]) {
            const preset = apiPresets[selectedType];
            elements.urlInput.value = preset.url;
            elements.methodSelect.value = preset.method;
            elements.requestBodyTextarea.value = preset.body;
            updateMethodVisibility();
            updateUrlPreview();
            showResult(`已加载预设：${preset.description}`, 'info');
        }
    }

    // 更新方法可见性（GET 请求隐藏请求体）
    function updateMethodVisibility() {
        const isGet = elements.methodSelect.value === 'GET';
        elements.bodyGroup.style.display = isGet ? 'none' : 'block';
    }

    // 更新 URL 预览
    let lastUrlPreview = null;
    function updateUrlPreview() {
        const url = elements.urlInput.value;
        let newPreview = '';
        let shouldShow = false;
        
        if (url) {
            try {
                const urlObj = new URL(url);
                newPreview = `${urlObj.hostname}${urlObj.pathname}`;
                shouldShow = true;
            } catch {
                shouldShow = false;
            }
        } else {
            shouldShow = false;
        }
        
        // 防止重复更新
        if (lastUrlPreview !== newPreview || elements.urlPreview.style.display !== (shouldShow ? 'block' : 'none')) {
            if (shouldShow) {
                elements.urlPreview.textContent = newPreview;
                elements.urlPreview.style.display = 'block';
            } else {
                elements.urlPreview.style.display = 'none';
            }
            lastUrlPreview = newPreview;
        }
    }

    // 更新 API 预设中的公司 ID
    function updateApiPresets() {
        const companyId = elements.companyIdInput.value;
        if (companyId) {
            apiPresets['finance-info'].body = `{"marketplaceSellerId":"${companyId}"}`;
            apiPresets['get-language'].body = `{"company_id":${companyId}}`;
            
            // 如果当前选择的是预设，自动更新请求体
            const currentType = elements.apiTypeSelect.value;
            if (apiPresets[currentType] && elements.methodSelect.value === 'POST') {
                elements.requestBodyTextarea.value = apiPresets[currentType].body;
            }
        }
    }

    // 发送请求
    async function handleSendRequest() {
        const url = elements.urlInput.value.trim();
        const method = elements.methodSelect.value;
        const companyId = elements.companyIdInput.value.trim();
        const customHeaders = elements.customHeadersTextarea.value.trim();
        const requestBody = elements.requestBodyTextarea.value.trim();

        // 验证输入
        if (!url) {
            showResult('❌ 请输入 URL', 'error');
            return;
        }

        if (!validateJSON(elements.customHeadersTextarea, elements.headersValidation) && customHeaders) {
            showResult('❌ 自定义请求头格式错误', 'error');
            return;
        }

        if (!validateJSON(elements.requestBodyTextarea, elements.bodyValidation) && requestBody && method !== 'GET') {
            showResult('❌ 请求体格式错误', 'error');
            return;
        }

        // 设置加载状态
        setLoadingState(true);
        showResult('🔄 正在发送请求...', 'loading');

        try {
            // 获取存储的 cookies
            const storedData = await getFromStorage(['ozonCookies']);
            
            // 构建请求头
            const headers = {
                'accept': 'application/json, text/plain, */*',
                'accept-language': 'zh-Hans',
                'content-type': 'application/json',
                'origin': 'https://seller.ozon.ru',
                'priority': 'u=1, i',
                'referer': 'https://seller.ozon.ru/app/registration/signin?__rr=1&abt_att=2&origin_referer=seller.ozon.ru',
                'sec-ch-ua': '"Not;A=Brand";v="99", "Google Chrome";v="139", "Chromium";v="139"',
                'sec-ch-ua-mobile': '?0',
                'sec-ch-ua-platform': '"Windows"',
                'sec-fetch-dest': 'empty',
                'sec-fetch-mode': 'cors',
                'sec-fetch-site': 'same-origin',
                'user-agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/139.0.0.0 Safari/537.36',
                'x-o3-app-name': 'seller-ui',
                'x-o3-language': 'zh-Hans',
                'x-o3-page-type': 'registration'
            };

            // 添加公司 ID 头部
            if (companyId) {
                headers['x-o3-company-id'] = companyId;
            }

            // 添加自定义头部
            if (customHeaders) {
                const customHeadersObj = JSON.parse(customHeaders);
                Object.assign(headers, customHeadersObj);
            }

            // 构建请求配置
            const requestConfig = {
                method: method,
                headers: headers,
                credentials: 'include'
            };

            // 添加请求体
            if (method !== 'GET' && requestBody) {
                requestConfig.body = requestBody;
            }

            // 发送请求到 content script
            const [tab] = await chrome.tabs.query({active: true, currentWindow: true});
            
            const response = await chrome.tabs.sendMessage(tab.id, {
                action: 'makeRequest',
                url: url,
                config: requestConfig,
                cookies: storedData.ozonCookies || []
            });

            if (response && response.success) {
                const result = `✅ 请求成功！

📊 状态码: ${response.status} ${response.statusText}

📅 时间: ${new Date().toLocaleString()}

📝 响应数据:
${JSON.stringify(response.data, null, 2)}

📋 响应头:
${JSON.stringify(response.headers, null, 2)}`;

                showResult(result, 'success');
                
                // 保存到历史记录
                await saveToHistory({
                    url,
                    method,
                    headers,
                    body: requestBody,
                    response: response.data,
                    status: response.status,
                    timestamp: Date.now()
                });

            } else {
                showResult(`❌ 请求失败: ${response?.error || '未知错误'}`, 'error');
            }

        } catch (error) {
            showResult(`💥 发送请求时出错: ${error.message}`, 'error');
        } finally {
            setLoadingState(false);
        }
    }

    // 设置加载状态
    let isLoading = false;
    function setLoadingState(loading) {
        if (isLoading === loading) {
            return; // 防止重复设置相同状态
        }
        isLoading = loading;
        elements.sendButton.disabled = loading;
        elements.sendText.style.display = loading ? 'none' : 'inline';
        elements.loadingSpinner.style.display = loading ? 'inline-block' : 'none';
        
        // 只有在加载状态时才改变状态指示器，否则保持连接状态
        if (loading) {
            elements.connectionStatus.className = 'status-indicator status-loading';
            elements.statusText.textContent = '请求中...';
        } else {
            // 请求完成后，恢复到连接状态检查
            updateConnectionStatus();
        }
    }

    // 清空表单
    function handleClearForm() {
        elements.urlInput.value = '';
        elements.customHeadersTextarea.value = '';
        elements.requestBodyTextarea.value = '';
        elements.resultDiv.innerHTML = '';
        elements.apiTypeSelect.value = 'finance-info';
        elements.methodSelect.value = 'POST';
        updateMethodVisibility();
        updateUrlPreview();
        showResult('✨ 表单已清空', 'info');
    }

    // 提取 cookies
    async function handleExtractCookies() {
        try {
            const [tab] = await chrome.tabs.query({active: true, currentWindow: true});
            
            if (!tab.url.includes('ozon.ru')) {
                showResult('⚠️ 请在 Ozon 页面上提取 cookies', 'warning');
                return;
            }

            const response = await chrome.tabs.sendMessage(tab.id, {
                action: 'extractCookies'
            });

            if (response && response.success) {
                await saveToStorage('ozonCookies', response.cookies);
                updateCookieCount(response.cookies.length);
                showResult(`✅ 已提取并保存 ${response.cookies.length} 个 cookies`, 'success');
                loadCookieList();
            } else {
                showResult(`❌ 提取 cookies 失败: ${response?.error || '未知错误'}`, 'error');
            }
        } catch (error) {
            showResult(`💥 提取 cookies 时出错: ${error.message}`, 'error');
        }
    }

    // 清空 cookies
    async function handleClearCookies() {
        await saveToStorage('ozonCookies', []);
        updateCookieCount(0);
        elements.cookieList.value = '';
        showResult('🗑️ Cookies 已清空', 'info');
    }

    // 解析 cookies 字符串
    async function handleParseCookies() {
        const cookieString = elements.cookieString.value.trim();
        if (!cookieString) {
            showResult('❌ 请输入 cookie 字符串', 'error');
            return;
        }

        try {
            const cookies = parseCookieString(cookieString);
            await saveToStorage('ozonCookies', cookies);
            updateCookieCount(cookies.length);
            elements.cookieString.value = '';
            loadCookieList();
            showResult(`✅ 已解析并保存 ${cookies.length} 个 cookies`, 'success');
        } catch (error) {
            showResult(`❌ 解析 cookies 失败: ${error.message}`, 'error');
        }
    }

    // 解析 cookie 字符串
    function parseCookieString(cookieString) {
        return cookieString.split(';').map(cookie => {
            const [name, ...valueParts] = cookie.trim().split('=');
            return {
                name: name.trim(),
                value: valueParts.join('=').trim()
            };
        }).filter(cookie => cookie.name && cookie.value);
    }

    // 更新连接状态
    let lastConnectionStatus = null;
    let isUpdatingStatus = false;
    async function updateConnectionStatus() {
        // 防止重复调用
        if (isUpdatingStatus) {
            return;
        }
        
        isUpdatingStatus = true;
        try {
            const storedData = await getFromStorage(['ozonCookies']);
            const cookieCount = storedData.ozonCookies ? storedData.ozonCookies.length : 0;
            updateCookieCount(cookieCount);
            
            const newStatus = cookieCount > 0 ? 'status-success' : 'status-error';
            const newStatusClass = `status-indicator ${newStatus}`;
            
            // 防止重复设置相同状态
            if (lastConnectionStatus !== newStatusClass) {
                elements.connectionStatus.className = newStatusClass;
                lastConnectionStatus = newStatusClass;
            }
            
            elements.statusText.textContent = cookieCount > 0 ? '已连接' : '未连接';
        } catch (error) {
            const errorStatusClass = 'status-indicator status-error';
            if (lastConnectionStatus !== errorStatusClass) {
                elements.connectionStatus.className = errorStatusClass;
                lastConnectionStatus = errorStatusClass;
            }
            elements.statusText.textContent = '状态未知';
        } finally {
            isUpdatingStatus = false;
        }
    }

    // 更新 cookie 计数
    function updateCookieCount(count) {
        elements.cookieCount.textContent = `${count} cookies`;
        elements.cookieCount.style.background = count > 0 ? '#10b981' : '#ef4444';
    }

    // 加载 cookie 列表
    let cookieListLoaded = false;
    let lastCookieContent = null;
    async function loadCookieList() {
        try {
            const storedData = await getFromStorage(['ozonCookies']);
            const cookies = storedData.ozonCookies || [];
            
            const newContent = cookies.length > 0 
                ? cookies.map(cookie => `${cookie.name}=${cookie.value}`).join('\n')
                : '暂无 cookies 数据';
            
            // 防止重复设置相同内容
            if (lastCookieContent !== newContent) {
                elements.cookieList.value = newContent;
                lastCookieContent = newContent;
            }
            
            cookieListLoaded = true;
        } catch (error) {
            console.error('加载cookie列表失败:', error);
        }
    }

    // 保存到历史记录
    async function saveToHistory(requestData) {
        const storedData = await getFromStorage(['requestHistory']);
        const history = storedData.requestHistory || [];
        
        history.unshift(requestData);
        
        // 只保留最近 50 条记录
        if (history.length > 50) {
            history.splice(50);
        }
        
        await saveToStorage('requestHistory', history);
    }

    // 加载历史记录列表
    let historyListLoaded = false;
    let lastHistoryContent = null;
    async function loadHistoryList() {
        try {
            const storedData = await getFromStorage(['requestHistory']);
            const history = storedData.requestHistory || [];
            
            let newContent;
            if (history.length === 0) {
                newContent = '<div class="info-panel">暂无请求历史</div>';
            } else {
                newContent = history.map((item, index) => `
                    <div class="history-item" style="border: 1px solid #e2e8f0; border-radius: 6px; padding: 12px; margin-bottom: 8px; cursor: pointer;" onclick="loadHistoryItem(${index})">
                        <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 8px;">
                            <span style="font-weight: 500;">${item.method} ${item.url.split('/').pop()}</span>
                            <span style="font-size: 11px; color: #64748b;">${new Date(item.timestamp).toLocaleString()}</span>
                        </div>
                        <div style="font-size: 12px; color: #64748b;">状态: ${item.status} | ${item.url}</div>
                    </div>
                `).join('');
            }
            
            // 防止重复设置相同内容
            if (lastHistoryContent !== newContent) {
                elements.historyList.innerHTML = newContent;
                lastHistoryContent = newContent;
            }
            
            historyListLoaded = true;
        } catch (error) {
            console.error('加载历史列表失败:', error);
        }
    }

    // 清空历史记录
    async function handleClearHistory() {
        await saveToStorage('requestHistory', []);
        loadHistoryList();
        showResult('🗑️ 历史记录已清空', 'info');
    }

    // 导出历史记录
    async function handleExportHistory() {
        const storedData = await getFromStorage(['requestHistory']);
        const history = storedData.requestHistory || [];
        
        if (history.length === 0) {
            showResult('❌ 没有历史记录可导出', 'warning');
            return;
        }

        const dataStr = JSON.stringify(history, null, 2);
        const dataBlob = new Blob([dataStr], { type: 'application/json' });
        const url = URL.createObjectURL(dataBlob);
        
        const a = document.createElement('a');
        a.href = url;
        a.download = `ozon-api-history-${new Date().toISOString().slice(0, 10)}.json`;
        a.click();
        
        URL.revokeObjectURL(url);
        showResult(`📤 已导出 ${history.length} 条历史记录`, 'success');
    }

    // 加载存储的数据
    async function loadStoredData() {
        const stored = await getFromStorage(['companyId', 'customHeaders', 'lastUrl', 'lastBody', 'ozonCookies']);
        
        if (stored.companyId) {
            elements.companyIdInput.value = stored.companyId;
            updateApiPresets();
        }
        
        if (stored.customHeaders) {
            elements.customHeadersTextarea.value = stored.customHeaders;
        }
        
        if (stored.lastUrl) {
            elements.urlInput.value = stored.lastUrl;
            updateUrlPreview();
        }
        
        if (stored.lastBody) {
            elements.requestBodyTextarea.value = stored.lastBody;
        }

        // 验证 JSON 格式
        validateJSON(elements.customHeadersTextarea, elements.headersValidation);
        validateJSON(elements.requestBodyTextarea, elements.bodyValidation);
    }

    // 显示结果
    function showResult(message, type) {
        elements.resultDiv.innerHTML = `<pre>${message}</pre>`;
        elements.resultDiv.className = `result-${type}`;
        elements.resultDiv.scrollTop = 0;
    }

    // 存储工具函数
    async function saveToStorage(key, value) {
        await chrome.storage.local.set({ [key]: value });
    }

    async function getFromStorage(keys) {
        return await chrome.storage.local.get(keys);
    }

    // 全局函数：加载历史记录项
    window.loadHistoryItem = function(index) {
        getFromStorage(['requestHistory']).then(storedData => {
            const history = storedData.requestHistory || [];
            const item = history[index];
            
            if (item) {
                elements.urlInput.value = item.url;
                elements.methodSelect.value = item.method;
                elements.requestBodyTextarea.value = item.body || '';
                
                // 切换到请求配置标签页
                switchTab('request');
                
                showResult(`📋 已加载历史记录: ${item.method} ${item.url}`, 'info');
                updateMethodVisibility();
                updateUrlPreview();
            }
        });
    };

    // 初始化应用
    init();
});


