# HotGo 2.0 前端开发文档

> 一个基于 Vue 3 + TypeScript + Naive UI 的现代化前端管理系统完整开发指南

## 📖 文档导航

### 🚀 快速开始
- [项目介绍](guide/introduction.md) - 了解 HotGo 2.0 的技术特色和应用场景
- [快速上手](guide/getting-started.md) - 环境搭建和项目启动指南
- [项目结构](guide/project-structure.md) - 详细的目录结构和文件组织

### 🏗️ 架构设计
- [核心架构](architecture/core.md) - 项目整体架构设计和设计原则
- [状态管理](architecture/state-management.md) - Pinia 状态管理完整方案
- [路由系统](architecture/routing.md) - 动态路由生成和权限控制
- [HTTP 请求](architecture/http.md) - Axios 封装和请求处理架构

### 🧩 组件文档
- [组件概览](components/overview.md) - 组件库架构和设计原则
- [表格组件](components/table.md) - BasicTable 详细使用指南
- [表单组件](components/form.md) - BasicForm 使用指南
- [模态框组件](components/modal.md) - BasicModal 使用指南
- [上传组件](components/upload.md) - 文件上传组件
- [其他组件](components/others.md) - 其他通用组件

### 🎣 Hooks 文档
- [Hooks 概览](hooks/overview.md) - 组合式函数介绍和设计模式
- [通用 Hooks](hooks/common.md) - 常用工具 hooks
- [业务 Hooks](hooks/business.md) - 业务相关 hooks
- [权限 Hooks](hooks/permission.md) - 权限相关 hooks

### 📡 API 文档
- [API 设计规范](api/design.md) - 接口设计原则和规范
- [接口封装](api/encapsulation.md) - 请求封装方案
- [错误处理](api/error-handling.md) - 统一错误处理
- [接口文档](api/endpoints.md) - 具体接口说明

### 🎨 样式系统
- [主题配置](styles/theme.md) - 主题系统说明
- [样式规范](styles/guidelines.md) - CSS 编写规范
- [Tailwind CSS](styles/tailwind.md) - 原子化样式使用

### 🔧 开发指南
- [开发规范](development/standards.md) - 代码规范和最佳实践
- [组件开发](development/component-development.md) - 组件开发指南
- [调试技巧](development/debugging.md) - 开发调试技巧
- [性能优化](development/performance.md) - 性能优化指南

### 🚀 部署运维
- [构建配置](deployment/build.md) - 构建相关配置
- [环境配置](deployment/environment.md) - 不同环境配置
- [部署指南](deployment/deployment.md) - 生产环境部署
- [监控运维](deployment/monitoring.md) - 监控和运维

## 🌟 项目特色

### 💎 技术栈
- **Vue 3.4.38** - 最新的 Vue 3 框架，使用 Composition API
- **TypeScript 5.5.4** - 完整的类型安全支持
- **Naive UI 2.42.0** - 优秀的 Vue 3 组件库
- **Vite 5.4.2** - 极速的开发体验和构建工具
- **Pinia 2.2.2** - Vue 3 官方推荐的状态管理
- **Vue Router 4.4.3** - 强大的路由管理

### ⚡ 核心功能
- 🔐 **完整的权限系统** - RBAC 权限模型，精确到按钮级别
- 📊 **丰富的组件库** - 开箱即用的高质量组件
- 🎨 **主题系统** - 支持多主题切换和暗黑模式
- 📱 **响应式设计** - 完美适配各种设备
- 🚀 **性能优化** - 代码分割、懒加载、虚拟滚动等
- 🔧 **开发工具** - 完整的开发工具链和调试支持

### 🏆 设计亮点
- **模块化架构** - 清晰的分层设计，易于维护和扩展
- **TypeScript 优先** - 全面的类型支持，提升开发体验
- **组件化开发** - 高度复用的组件设计
- **插件化扩展** - 支持功能模块的插件化开发

## 📋 快速开始

### 环境要求
- Node.js >= 16.0.0
- pnpm >= 7.0.0
- Git >= 2.0.0

### 安装和启动

```bash
# 克隆项目
git clone https://github.com/bufanyun/hotgo.git
cd hotgo/web

# 安装依赖
pnpm install

# 启动开发服务器
pnpm run dev

# 访问 http://localhost:3100
```

### 文档本地预览

```bash
# 启动文档服务器
pnpm run docs:serve

# 访问 http://localhost:3001
```

## 🤝 参与贡献

我们欢迎所有形式的贡献！请阅读[贡献指南](contributing.md)了解如何参与项目开发。

### 贡献方式
- 🐛 [报告问题](https://github.com/bufanyun/hotgo/issues)
- 💡 [功能建议](https://github.com/bufanyun/hotgo/issues)
- 📝 改进文档
- 💻 提交代码

## 📞 获取帮助

- 📖 [常见问题](faq.md)
- 💬 [GitHub Discussions](https://github.com/bufanyun/hotgo/discussions)
- 🐛 [Issue 反馈](https://github.com/bufanyun/hotgo/issues)
- 🌐 [官方网站](https://hotgo.facms.cn)

## 📄 许可证

本项目采用 [MIT](https://github.com/bufanyun/hotgo/blob/main/LICENSE) 许可证。

---

**Happy Coding! 🎉**