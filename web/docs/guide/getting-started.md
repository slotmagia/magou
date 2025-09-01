# 快速上手

本指南将帮助您快速搭建开发环境并启动项目。

## 环境准备

### 1. 安装 Node.js

确保您的系统已安装 Node.js 16.0.0 或更高版本。

```bash
# 查看 Node.js 版本
node --version

# 查看 npm 版本
npm --version
```

### 2. 安装 pnpm

推荐使用 pnpm 作为包管理器：

```bash
# 全局安装 pnpm
npm install -g pnpm

# 查看版本
pnpm --version
```

### 3. 安装 Git

确保系统已安装 Git：

```bash
# 查看 Git 版本
git --version
```

## 获取项目

### 克隆仓库

```bash
# 克隆项目
git clone https://github.com/bufanyun/hotgo.git

# 进入项目目录
cd hotgo/web
```

### 安装依赖

```bash
# 安装项目依赖
pnpm install
```

## 开发环境配置

### 1. 环境变量配置

项目支持多环境配置，复制环境变量文件：

```bash
# 复制开发环境配置
cp .env.example .env.development

# 复制生产环境配置
cp .env.example .env.production
```

### 2. 修改配置文件

编辑 `.env.development` 文件：

```bash
# 开发环境配置
VITE_NODE_ENV=development

# 项目基础路径
VITE_PUBLIC_PATH=/

# 开发服务器端口
VITE_PORT=3100

# API 接口地址
VITE_GLOB_API_URL=http://localhost:8000

# 接口前缀
VITE_GLOB_API_URL_PREFIX=/api

# 是否启用 Mock
VITE_USE_MOCK=true

# 是否启用 PWA
VITE_USE_PWA=false

# 是否开启包分析
VITE_USE_ANALYZE=false

# 是否启用 gzip 压缩
VITE_BUILD_GZIP=false

# 是否删除 console
VITE_DROP_CONSOLE=false
```

## 启动项目

### 开发模式

```bash
# 启动开发服务器
pnpm run dev

# 或者使用
pnpm run serve
```

启动成功后，浏览器会自动打开 `http://localhost:3100`

### 构建项目

```bash
# 构建生产版本
pnpm run build

# 预览构建结果
pnpm run preview
```

## 开发工具配置

### VS Code 配置

推荐使用 VS Code 作为开发编辑器，并安装以下插件：

#### 必备插件
- **Vetur** - Vue 语法高亮和智能提示
- **TypeScript Importer** - 自动导入 TypeScript 模块
- **ESLint** - 代码规范检查
- **Prettier** - 代码格式化
- **Auto Rename Tag** - 自动重命名配对标签

#### 推荐插件
- **GitLens** - Git 增强工具
- **Bracket Pair Colorizer** - 括号配对着色
- **Material Icon Theme** - 文件图标主题
- **Path Intellisense** - 路径智能提示

### VS Code 设置

在项目根目录创建 `.vscode/settings.json`：

```json
{
  "editor.formatOnSave": true,
  "editor.defaultFormatter": "esbenp.prettier-vscode",
  "editor.codeActionsOnSave": {
    "source.fixAll.eslint": true
  },
  "typescript.preferences.importModuleSpecifier": "relative",
  "vue.codeActions.enabled": true,
  "vue.complete.casing.tags": "kebab",
  "vue.complete.casing.props": "camel"
}
```

## 项目脚本说明

```json
{
  "scripts": {
    "dev": "vite",                              // 启动开发服务器
    "build": "vite build",                      // 构建生产版本
    "preview": "vite preview",                  // 预览构建结果
    "lint:eslint": "eslint src --fix",          // ESLint 检查和修复
    "lint:prettier": "prettier --write src",    // Prettier 格式化
    "lint:stylelint": "stylelint src/**/*.{vue,css,less} --fix", // 样式检查
    "type-check": "vue-tsc --noEmit",          // TypeScript 类型检查
    "clean:cache": "rimraf node_modules/.cache", // 清理缓存
    "clean:lib": "rimraf node_modules"          // 清理依赖
  }
}
```

## 开发流程

### 1. 功能开发

```bash
# 创建新分支
git checkout -b feature/new-feature

# 开发功能...

# 代码检查
pnpm run lint:eslint
pnpm run lint:prettier

# 类型检查
pnpm run type-check

# 提交代码
git add .
git commit -m "feat: 添加新功能"
```

### 2. 代码规范

项目使用严格的代码规范：

- **ESLint** - JavaScript/TypeScript 代码规范
- **Prettier** - 代码格式化
- **Stylelint** - CSS/Less 样式规范
- **Commitizen** - 提交信息规范

### 3. Git 提交规范

项目使用 [Conventional Commits](https://conventionalcommits.org/) 规范：

```bash
# 功能开发
git commit -m "feat: 添加用户管理功能"

# 问题修复
git commit -m "fix: 修复登录页面样式问题"

# 文档更新
git commit -m "docs: 更新 API 文档"

# 样式调整
git commit -m "style: 调整按钮样式"

# 代码重构
git commit -m "refactor: 重构用户服务"

# 性能优化
git commit -m "perf: 优化表格渲染性能"

# 测试相关
git commit -m "test: 添加用户服务测试"

# 构建相关
git commit -m "build: 更新构建配置"

# CI 相关
git commit -m "ci: 更新 GitHub Actions 配置"

# 其他杂项
git commit -m "chore: 更新依赖版本"
```

## 常见问题

### 1. 依赖安装失败

```bash
# 清理缓存重新安装
pnpm run clean:cache
pnpm install

# 或使用 npm
rm -rf node_modules package-lock.json
npm install
```

### 2. 端口被占用

```bash
# 查看端口占用
netstat -ano | findstr :3100

# 修改端口
# 在 .env.development 中修改 VITE_PORT
```

### 3. TypeScript 类型错误

```bash
# 运行类型检查
pnpm run type-check

# 查看具体错误信息
npx vue-tsc --noEmit
```

### 4. ESLint 错误

```bash
# 自动修复 ESLint 错误
pnpm run lint:eslint

# 查看具体错误
npx eslint src --ext .ts,.vue
```

## 下一步

- 📖 [项目结构](./project-structure.md) - 了解项目目录结构
- 🏗️ [核心架构](../architecture/core.md) - 深入了解项目架构
- 🧩 [组件文档](../components/overview.md) - 学习组件使用

## 技术支持

如果您在环境搭建过程中遇到问题：

1. 查看 [常见问题](../faq.md)
2. 搜索 [GitHub Issues](https://github.com/bufanyun/hotgo/issues)
3. 提交新的 [Issue](https://github.com/bufanyun/hotgo/issues/new)

---

恭喜！您已经成功搭建了开发环境，可以开始愉快的开发之旅了！ 🎉

