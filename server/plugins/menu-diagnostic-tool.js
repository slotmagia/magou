/**
 * 菜单渲染诊断工具
 * 用于诊断前端菜单不渲染的问题
 */

(function() {
    'use strict';

    // 诊断结果存储
    const diagnostics = {
        timestamp: new Date().toISOString(),
        checks: [],
        recommendations: []
    };

    // 执行诊断
    function runDiagnostics() {
        console.log('🔍 开始菜单渲染诊断...');

        // 检查1：API连通性
        checkApiConnectivity();

        // 检查2：数据结构验证
        checkDataStructure();

        // 检查3：前端代码检查
        checkFrontendCode();

        // 检查4：DOM检查
        checkDOMElements();

        // 检查5：网络请求检查
        checkNetworkRequests();

        // 生成报告
        generateReport();
    }

    // 检查API连通性
    function checkApiConnectivity() {
        console.log('📡 检查API连通性...');

        fetch('http://localhost:8001/api/menu/routers', {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': localStorage.getItem('token') || ''
            }
        })
        .then(response => response.json())
        .then(data => {
            diagnostics.checks.push({
                name: 'API连通性检查',
                status: 'success',
                message: 'API响应正常',
                details: {
                    status: data.code,
                    message: data.message,
                    dataLength: data.data ? data.data.list.length : 0,
                    hasData: !!data.data && !!data.data.list
                }
            });
        })
        .catch(error => {
            diagnostics.checks.push({
                name: 'API连通性检查',
                status: 'error',
                message: 'API请求失败: ' + error.message,
                details: { error: error.message }
            });
        });
    }

    // 检查数据结构
    function checkDataStructure() {
        console.log('📊 检查数据结构...');

        // 模拟用户提供的数据进行验证
        const sampleData = {
            "code": 0,
            "message": "操作成功",
            "data": {
                "list": [
                    {
                        "id": 1,
                        "name": "",
                        "path": "/dashboard",
                        "component": "dashboard/index",
                        "alwaysShow": false,
                        "hidden": false,
                        "meta": {
                            "title": "",
                            "icon": "dashboard",
                            "noCache": false,
                            "breadcrumb": false,
                            "permissions": "dashboard"
                        }
                    }
                ]
            },
            "timestamp": 1756258242,
            "traceID": "ec80ead3376e8059ad3f5a6b316dd7d2"
        };

        const issues = [];

        // 检查基本结构
        if (sampleData.code !== 0) {
            issues.push('响应码不是0，可能表示错误');
        }

        if (!sampleData.data) {
            issues.push('缺少data字段');
        }

        if (!sampleData.data.list) {
            issues.push('缺少data.list字段');
        }

        // 检查菜单项结构
        if (sampleData.data.list.length > 0) {
            const menuItem = sampleData.data.list[0];
            const requiredFields = ['id', 'path', 'component', 'meta'];
            const missingFields = requiredFields.filter(field => !(field in menuItem));

            if (missingFields.length > 0) {
                issues.push(`菜单项缺少必要字段: ${missingFields.join(', ')}`);
            }

            // 检查meta字段
            if (menuItem.meta) {
                const requiredMetaFields = ['title', 'icon', 'permissions'];
                const missingMetaFields = requiredMetaFields.filter(field => !(field in menuItem.meta));
                if (missingMetaFields.length > 0) {
                    issues.push(`菜单项meta字段缺少必要字段: ${missingMetaFields.join(', ')}`);
                }
            } else {
                issues.push('菜单项缺少meta字段');
            }
        }

        diagnostics.checks.push({
            name: '数据结构检查',
            status: issues.length === 0 ? 'success' : 'warning',
            message: issues.length === 0 ? '数据结构完整' : `发现${issues.length}个问题`,
            details: { issues }
        });
    }

    // 检查前端代码
    function checkFrontendCode() {
        console.log('💻 检查前端代码...');

        const issues = [];

        // 检查是否加载了必要的库
        if (typeof Vue === 'undefined' && typeof React === 'undefined') {
            issues.push('未检测到Vue.js或React，可能前端框架未加载');
        }

        // 检查是否有菜单相关的组件
        const menuElements = document.querySelectorAll('[data-menu], .menu, .sidebar, .nav');
        if (menuElements.length === 0) {
            issues.push('未找到菜单相关的DOM元素');
        }

        // 检查JavaScript错误
        if (window.console && window.console.errors) {
            issues.push('检测到JavaScript错误，可能影响菜单渲染');
        }

        diagnostics.checks.push({
            name: '前端代码检查',
            status: issues.length === 0 ? 'success' : 'warning',
            message: issues.length === 0 ? '前端代码正常' : `发现${issues.length}个问题`,
            details: { issues, menuElements: menuElements.length }
        });
    }

    // 检查DOM元素
    function checkDOMElements() {
        console.log('🌐 检查DOM元素...');

        const menuSelectors = [
            '.sidebar',
            '.menu',
            '.navigation',
            '[role="navigation"]',
            '.nav-menu',
            '.el-menu',
            '.ant-menu'
        ];

        const foundElements = [];
        menuSelectors.forEach(selector => {
            const elements = document.querySelectorAll(selector);
            if (elements.length > 0) {
                foundElements.push({
                    selector,
                    count: elements.length,
                    visible: Array.from(elements).some(el => el.offsetWidth > 0 && el.offsetHeight > 0)
                });
            }
        });

        diagnostics.checks.push({
            name: 'DOM元素检查',
            status: foundElements.length > 0 ? 'success' : 'warning',
            message: foundElements.length > 0 ? `找到${foundElements.length}个菜单相关元素` : '未找到菜单相关DOM元素',
            details: { foundElements }
        });
    }

    // 检查网络请求
    function checkNetworkRequests() {
        console.log('🔗 检查网络请求...');

        // 使用Performance API检查网络请求
        if ('performance' in window && performance.getEntriesByType) {
            const entries = performance.getEntriesByType('resource');
            const apiRequests = entries.filter(entry =>
                entry.name.includes('/api/menu/') ||
                entry.name.includes('menu') && entry.name.includes('api')
            );

            diagnostics.checks.push({
                name: '网络请求检查',
                status: apiRequests.length > 0 ? 'success' : 'warning',
                message: apiRequests.length > 0 ? `找到${apiRequests.length}个菜单相关API请求` : '未找到菜单相关API请求',
                details: {
                    apiRequests: apiRequests.map(req => ({
                        url: req.name,
                        duration: req.duration,
                        status: req.responseEnd - req.requestStart > 0 ? 'success' : 'error'
                    }))
                }
            });
        }
    }

    // 生成报告
    function generateReport() {
        console.log('📋 生成诊断报告...');

        setTimeout(() => {
            // 生成建议
            generateRecommendations();

            // 输出完整报告
            console.group('🔍 菜单渲染诊断报告');
            console.log('诊断时间:', diagnostics.timestamp);
            console.log('检查项目:', diagnostics.checks.length);

            diagnostics.checks.forEach((check, index) => {
                console.group(`${index + 1}. ${check.name}`);
                console.log('状态:', check.status === 'success' ? '✅' : check.status === 'warning' ? '⚠️' : '❌');
                console.log('消息:', check.message);
                if (check.details) {
                    console.log('详情:', check.details);
                }
                console.groupEnd();
            });

            if (diagnostics.recommendations.length > 0) {
                console.group('💡 修复建议');
                diagnostics.recommendations.forEach((rec, index) => {
                    console.log(`${index + 1}. ${rec}`);
                });
                console.groupEnd();
            }

            console.groupEnd();

            // 创建可视化报告
            createVisualReport();
        }, 2000);
    }

    // 生成建议
    function generateRecommendations() {
        diagnostics.checks.forEach(check => {
            switch (check.name) {
                case 'API连通性检查':
                    if (check.status === 'error') {
                        diagnostics.recommendations.push('检查API服务器是否运行正常');
                        diagnostics.recommendations.push('验证API端点地址是否正确');
                        diagnostics.recommendations.push('检查网络连接和防火墙设置');
                    }
                    break;

                case '数据结构检查':
                    if (check.status === 'warning' && check.details.issues) {
                        check.details.issues.forEach(issue => {
                            diagnostics.recommendations.push(`修复数据结构问题: ${issue}`);
                        });
                    }
                    break;

                case '前端代码检查':
                    if (check.status === 'warning' && check.details.issues) {
                        check.details.issues.forEach(issue => {
                            diagnostics.recommendations.push(`修复前端问题: ${issue}`);
                        });
                    }
                    break;

                case 'DOM元素检查':
                    if (check.status === 'warning') {
                        diagnostics.recommendations.push('检查菜单组件是否正确导入和使用');
                        diagnostics.recommendations.push('验证CSS样式是否影响菜单显示');
                    }
                    break;

                case '网络请求检查':
                    if (check.status === 'warning') {
                        diagnostics.recommendations.push('确认前端代码中调用了菜单API');
                        diagnostics.recommendations.push('检查API调用时机是否正确');
                    }
                    break;
            }
        });

        // 通用建议
        diagnostics.recommendations.push('打开浏览器开发者工具查看Console错误信息');
        diagnostics.recommendations.push('检查Network标签页确认API请求是否发送成功');
        diagnostics.recommendations.push('验证用户权限设置是否正确');
    }

    // 创建可视化报告
    function createVisualReport() {
        const reportDiv = document.createElement('div');
        reportDiv.id = 'menu-diagnostic-report';
        reportDiv.innerHTML = `
            <div style="
                position: fixed;
                top: 20px;
                right: 20px;
                width: 350px;
                max-height: 500px;
                background: #fff;
                border: 2px solid #0066cc;
                border-radius: 8px;
                padding: 15px;
                z-index: 9999;
                overflow-y: auto;
                box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
                font-family: Arial, sans-serif;
                font-size: 12px;
            ">
                <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 15px;">
                    <h3 style="margin: 0; color: #0066cc;">🔍 菜单诊断报告</h3>
                    <button onclick="this.parentNode.parentNode.remove()" style="
                        background: #ef4444;
                        color: white;
                        border: none;
                        border-radius: 4px;
                        padding: 4px 8px;
                        cursor: pointer;
                        font-size: 10px;
                    ">✕</button>
                </div>

                <div id="diagnostic-results">
                    <div style="color: #666; margin-bottom: 10px;">正在生成报告...</div>
                </div>

                <div style="margin-top: 15px;">
                    <button onclick="runDiagnostics()" style="
                        background: #0066cc;
                        color: white;
                        border: none;
                        border-radius: 4px;
                        padding: 8px 16px;
                        cursor: pointer;
                        font-size: 11px;
                        margin-right: 8px;
                    ">🔄 重新诊断</button>
                    <button onclick="exportDiagnostics()" style="
                        background: #10b981;
                        color: white;
                        border: none;
                        border-radius: 4px;
                        padding: 8px 16px;
                        cursor: pointer;
                        font-size: 11px;
                    ">📤 导出报告</button>
                </div>
            </div>
        `;

        // 填充报告内容
        setTimeout(() => {
            const resultsDiv = reportDiv.querySelector('#diagnostic-results');
            resultsDiv.innerHTML = diagnostics.checks.map(check => `
                <div style="margin-bottom: 8px; padding: 8px; border-radius: 4px; ${
                    check.status === 'success' ? 'background: #f0fdf4; border-left: 4px solid #22c55e;' :
                    check.status === 'warning' ? 'background: #fffbeb; border-left: 4px solid #f59e0b;' :
                    'background: #fef2f2; border-left: 4px solid #ef4444;'
                }">
                    <div style="font-weight: bold; margin-bottom: 4px;">
                        ${check.status === 'success' ? '✅' : check.status === 'warning' ? '⚠️' : '❌'}
                        ${check.name}
                    </div>
                    <div style="color: #666; font-size: 11px;">${check.message}</div>
                </div>
            `).join('');

            if (diagnostics.recommendations.length > 0) {
                resultsDiv.innerHTML += `
                    <div style="margin-top: 15px; padding-top: 10px; border-top: 1px solid #e2e8f0;">
                        <div style="font-weight: bold; margin-bottom: 8px; color: #0066cc;">💡 修复建议:</div>
                        <ul style="margin: 0; padding-left: 20px;">
                            ${diagnostics.recommendations.map(rec => `<li style="margin-bottom: 4px; font-size: 11px;">${rec}</li>`).join('')}
                        </ul>
                    </div>
                `;
            }
        }, 2500);

        document.body.appendChild(reportDiv);
    }

    // 导出诊断数据
    window.exportDiagnostics = function() {
        const dataStr = JSON.stringify(diagnostics, null, 2);
        const dataBlob = new Blob([dataStr], { type: 'application/json' });
        const url = URL.createObjectURL(dataBlob);

        const a = document.createElement('a');
        a.href = url;
        a.download = `menu-diagnostic-${new Date().toISOString().slice(0, 19)}.json`;
        a.click();

        URL.revokeObjectURL(url);
    };

    // 暴露全局函数
    window.runDiagnostics = runDiagnostics;

    // 自动运行诊断（延迟执行，给页面加载时间）
    setTimeout(() => {
        runDiagnostics();
    }, 1000);

    console.log('🔧 菜单诊断工具已加载');
    console.log('💡 使用方法:');
    console.log('   1. 运行 runDiagnostics() 开始诊断');
    console.log('   2. 查看控制台输出结果');
    console.log('   3. 查看页面右上角的可视化报告');

})();


