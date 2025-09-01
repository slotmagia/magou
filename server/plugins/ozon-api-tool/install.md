# 📦 Ozon API 工具 - 安装指南

## 🚀 快速安装

### 步骤 1：准备文件
确保您有完整的插件文件：
```
ozon-api-tool/
├── manifest.json          # 插件配置文件
├── popup.html             # 弹窗界面
├── popup.js               # 弹窗逻辑
├── content.js             # 内容脚本
├── background.js          # 后台脚本
├── styles.css             # 样式文件
├── icons/                 # 图标目录
│   ├── icon16.png
│   ├── icon32.png
│   ├── icon48.png
│   └── icon128.png
├── README.md              # 说明文档
└── install.md             # 本安装指南
```

### 步骤 2：安装插件

#### Chrome/Edge 浏览器

1. **打开扩展管理页面**
   - Chrome：地址栏输入 `chrome://extensions/`
   - Edge：地址栏输入 `edge://extensions/`

2. **启用开发者模式**
   - 在页面右上角开启"开发者模式"开关

3. **加载插件**
   - 点击"加载已解压的扩展程序"
   - 选择 `ozon-api-tool` 文件夹
   - 点击"选择文件夹"

4. **确认安装**
   - 插件会出现在扩展列表中
   - 确保插件已启用（开关为蓝色）

### 步骤 3：验证安装

1. **检查插件图标**
   - 浏览器工具栏应显示插件图标
   - 点击图标应打开插件弹窗

2. **测试基本功能**
   - 访问 `https://seller.ozon.ru`
   - 点击插件图标
   - 尝试提取 cookies

## 🔧 图标文件

如果缺少图标文件，您可以：

### 方法 1：使用默认图标
创建以下尺寸的图标文件（可以是简单的彩色方块）：
- `icons/icon16.png` (16x16 像素)
- `icons/icon32.png` (32x32 像素)  
- `icons/icon48.png` (48x48 像素)
- `icons/icon128.png` (128x128 像素)

### 方法 2：临时移除图标引用
编辑 `manifest.json`，移除或注释掉图标相关配置：

```json
{
    "action": {
        "default_popup": "popup.html",
        "default_title": "Ozon API 工具"
        // 暂时移除图标配置
        // "default_icon": { ... }
    }
    // 暂时移除图标配置  
    // "icons": { ... }
}
```

## 🐛 常见安装问题

### 问题 1：无法加载插件
**错误信息**：`Manifest file is missing or unreadable`

**解决方案**：
- 确保 `manifest.json` 文件存在且格式正确
- 检查文件编码是否为 UTF-8
- 验证 JSON 语法是否正确

### 问题 2：权限错误
**错误信息**：`Permission denied`

**解决方案**：
- 确保已启用开发者模式
- 重新启动浏览器
- 尝试以管理员身份运行浏览器

### 问题 3：图标加载失败
**错误信息**：`Could not load icon`

**解决方案**：
- 添加图标文件到 `icons/` 目录
- 或临时移除 manifest.json 中的图标配置

### 问题 4：内容脚本注入失败
**错误信息**：`Content script injection failed`

**解决方案**：
- 确保在 Ozon 网站上使用插件
- 检查网站是否阻止脚本注入
- 重新加载页面和插件

## 🔄 更新插件

当插件代码有更新时：

1. **保存数据**（可选）
   - 在插件面板中导出配置和历史

2. **重新加载插件**
   - 在扩展管理页面点击"重新加载"按钮
   - 或删除插件后重新安装

3. **恢复数据**（可选）
   - 导入之前保存的配置

## 🎯 首次使用

安装完成后的首次设置：

1. **访问 Ozon 卖家平台**
   ```
   https://seller.ozon.ru
   ```

2. **登录您的账户**
   - 确保已正常登录

3. **提取认证信息**
   - 点击插件图标
   - 切换到"Cookie 管理"标签
   - 点击"提取 Cookies"

4. **测试 API 请求**
   - 选择"财务信息"预设
   - 输入您的公司 ID
   - 点击"发送请求"

## 💡 提示和技巧

### 固定插件到工具栏
- 点击浏览器工具栏的插件图标（拼图形状）
- 找到"Ozon API 工具"
- 点击固定图标（📌）

### 设置快捷键
1. 在扩展管理页面点击左侧菜单"键盘快捷键"
2. 为插件设置快捷键组合
3. 使用快捷键快速打开插件

### 开启调试模式
在浏览器控制台执行：
```javascript
chrome.runtime.sendMessage({action: 'toggleDebug'});
```

---

如果您在安装过程中遇到任何问题，请参考 README.md 中的故障排除部分。


