# 常见问题

本文档收集了 HotGo 2.0 项目开发和使用过程中的常见问题和解决方案。

## 环境问题

### Q: Node.js 版本要求是什么？

**A:** HotGo 2.0 要求 Node.js >= 16.0.0，推荐使用 18.17.0 或更高版本。

```bash
# 查看当前版本
node --version

# 如果版本过低，请升级 Node.js
# 推荐使用 nvm 管理 Node.js 版本
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.0/install.sh | bash
nvm install 18
nvm use 18
```

### Q: 为什么推荐使用 pnpm 而不是 npm？

**A:** pnpm 具有以下优势：

- 更快的安装速度
- 更少的磁盘空间占用
- 更严格的依赖管理
- 更好的 monorepo 支持

```bash
# 安装 pnpm
npm install -g pnpm

# 使用 pnpm 安装依赖
pnpm install
```

### Q: 安装依赖时出现网络错误怎么办？

**A:** 可以尝试以下解决方案：

```bash
# 1. 使用国内镜像源
pnpm config set registry https://registry.npmmirror.com

# 2. 清理缓存重新安装
pnpm store prune
rm -rf node_modules pnpm-lock.yaml
pnpm install

# 3. 使用代理
pnpm config set proxy http://your-proxy:port
pnpm config set https-proxy http://your-proxy:port
```

## 开发问题

### Q: 启动开发服务器时端口被占用怎么办？

**A:** 可以通过以下方式解决：

```bash
# 1. 查看端口占用情况
netstat -ano | findstr :3100  # Windows
lsof -i :3100                 # macOS/Linux

# 2. 杀掉占用端口的进程
kill -9 <PID>

# 3. 或者修改端口
# 在 .env.development 文件中修改 VITE_PORT
VITE_PORT=3200
```

### Q: TypeScript 类型检查报错怎么解决？

**A:** 常见的类型错误解决方案：

```typescript
// 1. 缺少类型定义
// 安装对应的 @types 包
pnpm install -D @types/node @types/lodash

// 2. any 类型错误
// 使用具体的类型定义
interface UserInfo {
  id: number;
  name: string;
}
const user: UserInfo = { id: 1, name: 'admin' };

// 3. 模块导入错误
// 确保路径和导出正确
import { BasicTable } from '@/components/Table';

// 4. 第三方库类型问题
// 在 types/global.d.ts 中声明
declare module 'your-module' {
  export function someFunction(): void;
}
```

### Q: ESLint 检查失败怎么办？

**A:** 根据不同的错误类型：

```bash
# 1. 自动修复大部分问题
pnpm run lint:eslint --fix

# 2. 查看具体错误
pnpm run lint:eslint

# 3. 忽略特定规则（谨慎使用）
// eslint-disable-next-line @typescript-eslint/no-explicit-any
const data: any = {};

# 4. 配置文件问题
# 检查 .eslintrc.js 配置是否正确
```

### Q: 组件样式不生效怎么解决？

**A:** 可能的原因和解决方案：

```vue
<!-- 1. 确保使用了 scoped 样式 -->
<style scoped lang="less">
.my-component {
  color: red;
}
</style>

<!-- 2. 深度选择器 -->
<style scoped lang="less">
.my-component {
  :deep(.child-element) {
    color: blue;
  }
}
</style>

<!-- 3. 全局样式 -->
<style lang="less">
.global-style {
  font-size: 14px;
}
</style>

<!-- 4. CSS 变量 -->
<style scoped lang="less">
.my-component {
  color: var(--primary-color);
}
</style>
```

## 组件问题

### Q: BasicTable 数据不显示怎么办？

**A:** 检查以下几个方面：

```typescript
// 1. 确保 dataSource 函数正确
const loadData = async (params) => {
  // 必须返回正确格式的数据
  return {
    items: [], // 数据数组
    total: 0,  // 总数量
  };
};

// 2. 检查列配置
const columns = [
  {
    title: '标题',
    key: 'field_name', // 确保 key 与数据字段匹配
    width: 100,
  },
];

// 3. 检查分页配置
const pagination = reactive({
  page: 1,
  pageSize: 10,
  showSizePicker: true,
});

// 4. 检查 API 响应格式
// 后端应返回符合约定的数据格式
```

### Q: BasicForm 表单验证不生效？

**A:** 检查表单配置：

```typescript
const schemas = [
  {
    field: 'username',
    label: '用户名',
    component: 'NInput',
    // 确保验证规则正确
    rules: [
      { required: true, message: '请输入用户名' },
      { min: 3, max: 20, message: '长度在 3 到 20 个字符' },
    ],
  },
];

// 使用表单验证
const [register, { validate }] = useForm({
  schemas,
});

const handleSubmit = async () => {
  try {
    const values = await validate();
    // 验证通过，提交数据
  } catch (error) {
    // 验证失败
    console.log('验证失败:', error);
  }
};
```

### Q: useModal 模态框显示异常？

**A:** 常见问题和解决方案：

```vue
<template>
  <!-- 确保模态框组件正确注册 -->
  <BasicModal
    @register="register"
    title="编辑用户"
    :width="600"
  >
    <div>模态框内容</div>
  </BasicModal>
</template>

<script setup lang="ts">
import { BasicModal, useModal } from '@/components/Modal';

// 确保正确使用 hook
const [register, { openModal, closeModal }] = useModal();

// 检查是否正确调用
const handleEdit = () => {
  openModal(); // 而不是 openModal(true)
};
</script>
```

## API 问题

### Q: API 请求失败怎么排查？

**A:** 按以下步骤排查：

```typescript
// 1. 检查网络连接
fetch('/api/health')
  .then(response => console.log('API 可达'))
  .catch(error => console.log('网络错误:', error));

// 2. 检查 API 地址配置
console.log('API_URL:', import.meta.env.VITE_GLOB_API_URL);

// 3. 查看浏览器开发者工具
// - Network 面板查看请求详情
// - Console 面板查看错误信息

// 4. 检查请求拦截器
// 确保 token 正确添加
const token = localStorage.getItem('token');
console.log('Token:', token);

// 5. 模拟请求测试
curl -H "Authorization: Bearer your-token" \
     http://localhost:8000/api/users
```

### Q: 跨域问题怎么解决？

**A:** 开发环境和生产环境的解决方案：

```typescript
// 开发环境：配置代理（vite.config.ts）
export default {
  server: {
    proxy: {
      '/api': {
        target: 'http://localhost:8000',
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api/, '/api'),
      },
    },
  },
};

// 生产环境：Nginx 配置
// location /api/ {
//   proxy_pass http://localhost:8000;
//   proxy_set_header Host $host;
//   proxy_set_header X-Real-IP $remote_addr;
// }

// 后端配置 CORS（如果需要）
app.use(cors({
  origin: ['http://localhost:3100', 'https://your-domain.com'],
  credentials: true,
}));
```

## 构建问题

### Q: 构建失败怎么解决？

**A:** 根据错误类型排查：

```bash
# 1. 清理缓存重新构建
pnpm run clean:cache
pnpm run build

# 2. 检查内存使用
# 增加 Node.js 内存限制
export NODE_OPTIONS="--max-old-space-size=4096"
pnpm run build

# 3. 分析构建错误
pnpm run build 2>&1 | tee build.log

# 4. 检查依赖版本冲突
pnpm list
pnpm outdated
```

### Q: 构建后白屏怎么解决？

**A:** 可能的原因和解决方案：

```bash
# 1. 检查构建产物
ls -la dist/
cat dist/index.html

# 2. 检查路径配置
# .env.production 中的 VITE_PUBLIC_PATH
VITE_PUBLIC_PATH=/

# 3. 检查服务器配置
# Nginx 需要配置 SPA 路由
location / {
  try_files $uri $uri/ /index.html;
}

# 4. 检查浏览器控制台
# 查看是否有 JavaScript 错误或资源加载失败
```

## 部署问题

### Q: Nginx 配置后访问 404？

**A:** 检查 Nginx 配置：

```nginx
server {
    listen 80;
    server_name your-domain.com;
    root /var/www/hotgo;
    index index.html;
    
    # 重要：SPA 路由配置
    location / {
        try_files $uri $uri/ /index.html;
    }
    
    # 检查文件权限
    # sudo chown -R www-data:www-data /var/www/hotgo
    # sudo chmod -R 755 /var/www/hotgo
}
```

### Q: Docker 部署失败怎么办？

**A:** 常见问题和解决方案：

```dockerfile
# 1. 检查 Dockerfile
FROM node:18-alpine AS builder
WORKDIR /app
COPY package.json pnpm-lock.yaml ./
RUN npm install -g pnpm
RUN pnpm install --frozen-lockfile
COPY . .
RUN pnpm run build

FROM nginx:alpine
COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/nginx.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]

# 2. 检查构建日志
docker build -t hotgo-web . --no-cache

# 3. 检查容器运行状态
docker logs container-name

# 4. 进入容器调试
docker exec -it container-name sh
```

## 性能问题

### Q: 页面加载慢怎么优化？

**A:** 性能优化建议：

```typescript
// 1. 路由懒加载
const routes = [
  {
    path: '/user',
    component: () => import('@/views/user/index.vue'),
  },
];

// 2. 组件懒加载
const HeavyComponent = defineAsyncComponent({
  loader: () => import('./HeavyComponent.vue'),
  loadingComponent: LoadingComponent,
});

// 3. 使用 v-show 替代频繁切换的 v-if
<template>
  <div v-show="visible">内容</div>
</template>

// 4. 大列表虚拟滚动
<VirtualList
  :items="largeList"
  :item-height="50"
  :visible-count="20"
/>

// 5. 图片懒加载
<img v-lazy="imageSrc" alt="description" />
```

### Q: 打包体积过大怎么优化？

**A:** 减少打包体积的方法：

```typescript
// 1. 按需导入
// ❌ 全量导入
import * as icons from '@vicons/antd';

// ✅ 按需导入
import { UserOutlined } from '@vicons/antd';

// 2. 配置外部依赖（CDN）
// vite.config.ts
export default {
  build: {
    rollupOptions: {
      external: ['vue', 'vue-router'],
      output: {
        globals: {
          vue: 'Vue',
          'vue-router': 'VueRouter',
        },
      },
    },
  },
};

// 3. 代码分割
const UserManagement = () => import('@/views/user/index.vue');

// 4. 移除未使用的代码
pnpm run build:analyze
```

## 其他问题

### Q: 如何升级到新版本？

**A:** 升级步骤：

```bash
# 1. 备份当前项目
git stash
git checkout -b backup-current

# 2. 查看更新日志
# 访问 GitHub Releases 或 CHANGELOG.md

# 3. 更新依赖
pnpm update

# 4. 解决冲突
# 根据迁移指南处理不兼容的更改

# 5. 测试功能
pnpm run dev
pnpm run build
```

### Q: 如何自定义主题？

**A:** 主题定制方法：

```less
// styles/var.less
@primary-color: #1890ff;
@success-color: #52c41a;
@warning-color: #faad14;
@error-color: #f5222d;

// 覆盖 Naive UI 主题
@n-color-primary: @primary-color;
@n-border-radius: 6px;
```

```typescript
// main.ts
import { createApp } from 'vue';
import { darkTheme, lightTheme } from 'naive-ui';

const app = createApp(App);

// 主题配置
const themeOverrides = {
  common: {
    primaryColor: '#1890ff',
    primaryColorHover: '#40a9ff',
    borderRadius: '6px',
  },
};
```

### Q: 如何贡献代码？

**A:** 请参考[贡献指南](contributing.md)，主要步骤：

1. Fork 项目
2. 创建功能分支
3. 开发和测试
4. 提交 Pull Request
5. 等待代码审查

### Q: 遇到未解决的问题怎么办？

**A:** 获取帮助的途径：

1. 查看[项目文档](README.md)
2. 搜索 [GitHub Issues](https://github.com/bufanyun/hotgo/issues)
3. 创建新的 Issue 描述问题
4. 参与 [GitHub Discussions](https://github.com/bufanyun/hotgo/discussions)
5. 查看[官方网站](https://hotgo.facms.cn)

---

如果本文档没有涵盖您遇到的问题，请随时在 GitHub Issues 中提出，我们会及时补充和完善。






