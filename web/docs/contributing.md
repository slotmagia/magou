# 贡献指南

感谢您对 HotGo 2.0 项目的关注！我们欢迎任何形式的贡献，包括但不限于代码贡献、文档改进、问题反馈和功能建议。

## 贡献方式

### 🐛 问题反馈

如果您在使用过程中遇到问题，请：

1. 在 [GitHub Issues](https://github.com/bufanyun/hotgo/issues) 中搜索是否已有相似问题
2. 如果没有找到，请创建新的 Issue
3. 提供详细的问题描述和复现步骤
4. 包含您的环境信息（操作系统、Node.js 版本、浏览器版本等）

### 💡 功能建议

我们欢迎您的功能建议：

1. 在 [GitHub Issues](https://github.com/bufanyun/hotgo/issues) 中创建功能请求
2. 使用 "Feature Request" 标签
3. 详细描述您期望的功能和使用场景
4. 说明该功能的价值和必要性

### 📝 文档贡献

文档改进也是重要的贡献：

1. Fork 项目到您的 GitHub 账户
2. 在 `docs/` 目录下修改或新增文档
3. 确保文档格式正确且内容准确
4. 提交 Pull Request

### 💻 代码贡献

我们欢迎代码贡献，请遵循以下流程：

## 开发流程

### 1. 环境准备

```bash
# 克隆仓库
git clone https://github.com/bufanyun/hotgo.git
cd hotgo/web

# 安装依赖
pnpm install

# 启动开发服务器
pnpm run dev
```

### 2. 创建分支

```bash
# 从 main 分支创建功能分支
git checkout -b feature/your-feature-name

# 或者从 main 分支创建修复分支
git checkout -b fix/your-fix-name
```

### 3. 开发和测试

在开发过程中，请：

- 遵循项目的[代码规范](development/standards.md)
- 编写必要的单元测试
- 确保代码通过所有检查

```bash
# 运行代码检查
pnpm run lint:eslint
pnpm run lint:prettier

# 运行类型检查
pnpm run type-check

# 运行测试
pnpm run test

# 构建检查
pnpm run build
```

### 4. 提交代码

使用规范的提交信息：

```bash
# 暂存文件
git add .

# 提交（使用规范的提交信息）
git commit -m "feat: 添加用户管理功能"

# 推送到远程仓库
git push origin feature/your-feature-name
```

### 5. 创建 Pull Request

1. 在 GitHub 上创建 Pull Request
2. 填写详细的 PR 描述
3. 关联相关的 Issue
4. 等待代码审查

## 代码规范

### 提交信息规范

我们使用 [Conventional Commits](https://conventionalcommits.org/) 规范：

```
<type>(<scope>): <description>

[optional body]

[optional footer(s)]
```

#### 类型说明

- `feat`: 新功能
- `fix`: 修复 bug
- `docs`: 文档变更
- `style`: 代码格式调整
- `refactor`: 重构
- `perf`: 性能优化
- `test`: 测试相关
- `chore`: 构建过程或辅助工具的变动

#### 示例

```bash
feat(user): 添加用户管理功能
fix(table): 修复表格分页显示异常
docs: 更新 API 文档
style(form): 调整表单样式
refactor(auth): 重构权限验证逻辑
```

### 代码风格

- 使用 TypeScript 编写代码
- 遵循 ESLint 和 Prettier 配置
- 组件命名使用 PascalCase
- 文件命名使用 kebab-case
- 函数和变量命名使用 camelCase

### 测试要求

- 新功能必须包含单元测试
- 修复 bug 时应添加回归测试
- 测试覆盖率不低于 80%
- 确保所有测试通过

## 代码审查

### 审查标准

我们的代码审查关注以下方面：

1. **功能正确性**: 代码是否实现了预期功能
2. **代码质量**: 是否遵循最佳实践和项目规范
3. **性能影响**: 是否会对性能产生负面影响
4. **安全性**: 是否存在安全隐患
5. **可维护性**: 代码是否易于理解和维护
6. **测试完整性**: 是否有足够的测试覆盖

### 审查流程

1. 自动化检查（CI/CD）
2. 至少一名核心开发者审查
3. 解决审查意见
4. 合并到主分支

## 发布流程

### 版本规范

我们使用 [Semantic Versioning](https://semver.org/) 版本规范：

- `MAJOR`: 不兼容的 API 变更
- `MINOR`: 向后兼容的功能性新增
- `PATCH`: 向后兼容的问题修正

### 发布步骤

1. 更新版本号
2. 更新 CHANGELOG.md
3. 创建 Release Tag
4. 自动构建和部署

## 社区行为准则

### 我们的承诺

为了营造一个开放且友好的环境，我们承诺：

- 使用友好和包容的语言
- 尊重不同的观点和经历
- 优雅地接受建设性批评
- 关注对社区最有利的事情
- 与其他社区成员友善相处

### 不当行为

以下行为被视为不当行为：

- 使用性化的言语或图像
- 人身攻击或政治攻击
- 公开或私下的骚扰
- 未经明确许可发布他人的私人信息
- 其他在专业环境中被认为不当的行为

### 执行

项目维护者有权移除、编辑或拒绝与本行为准则不符的评论、提交、代码、wiki 编辑、问题和其他贡献。

## 开发者资源

### 文档资源

- [开发规范](development/standards.md)
- [组件开发指南](development/component-development.md)
- [API 设计规范](api/design.md)
- [部署指南](deployment/deployment.md)

### 工具和环境

- **IDE**: 推荐使用 VS Code
- **Node.js**: >= 16.0.0
- **包管理器**: pnpm
- **代码格式化**: Prettier
- **代码检查**: ESLint
- **类型检查**: TypeScript

### 获取帮助

如果您在贡献过程中遇到问题：

1. 查看[常见问题](faq.md)
2. 在 [GitHub Issues](https://github.com/bufanyun/hotgo/issues) 中搜索
3. 创建新的 Issue 寻求帮助
4. 在 [GitHub Discussions](https://github.com/bufanyun/hotgo/discussions) 中讨论

## 贡献者名单

感谢所有为 HotGo 2.0 做出贡献的开发者！

<!-- 
这里将自动显示贡献者列表
可以使用 all-contributors 工具自动维护
-->

## 许可证

通过向本项目贡献代码，您同意您的贡献将在 [MIT 许可证](../LICENSE) 下获得许可。

---

再次感谢您的贡献！🎉






