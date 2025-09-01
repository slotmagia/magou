# 部署指南

本文档详细说明 HotGo 2.0 前端项目的部署流程和配置方法。

## 环境准备

### 系统要求

| 环境 | 最低版本 | 推荐版本 |
|------|----------|----------|
| Node.js | 16.0.0 | 18.17.0+ |
| pnpm | 7.0.0 | 8.6.0+ |
| Git | 2.0.0 | 2.40.0+ |

### 服务器要求

#### 开发环境
- CPU: 2核
- 内存: 4GB
- 磁盘: 50GB
- 网络: 10Mbps

#### 生产环境
- CPU: 4核+
- 内存: 8GB+
- 磁盘: 100GB+
- 网络: 100Mbps+

## 构建配置

### 环境变量配置

创建对应环境的配置文件：

#### 开发环境 (`.env.development`)

```bash
# 环境标识
NODE_ENV=development
VITE_NODE_ENV=development

# 应用配置
VITE_APP_TITLE=HotGo 管理系统
VITE_APP_SHORT_NAME=HotGo
VITE_PUBLIC_PATH=/

# 服务配置
VITE_PORT=3100
VITE_PROXY_TYPE=dev

# API配置
VITE_GLOB_API_URL=http://localhost:8000
VITE_GLOB_API_URL_PREFIX=/api
VITE_GLOB_UPLOAD_URL=/upload

# 功能开关
VITE_USE_MOCK=true
VITE_USE_PWA=false
VITE_USE_CDN=false

# 开发工具
VITE_BUILD_GZIP=false
VITE_BUILD_ANALYZE=false
VITE_DROP_CONSOLE=false
VITE_DROP_DEBUGGER=false
```

#### 测试环境 (`.env.staging`)

```bash
# 环境标识
NODE_ENV=production
VITE_NODE_ENV=staging

# 应用配置
VITE_APP_TITLE=HotGo 管理系统(测试)
VITE_APP_SHORT_NAME=HotGo
VITE_PUBLIC_PATH=/

# API配置
VITE_GLOB_API_URL=https://test-api.example.com
VITE_GLOB_API_URL_PREFIX=/api
VITE_GLOB_UPLOAD_URL=https://test-api.example.com/upload

# 功能开关
VITE_USE_MOCK=false
VITE_USE_PWA=false
VITE_USE_CDN=false

# 构建优化
VITE_BUILD_GZIP=true
VITE_BUILD_ANALYZE=false
VITE_DROP_CONSOLE=true
VITE_DROP_DEBUGGER=true
```

#### 生产环境 (`.env.production`)

```bash
# 环境标识
NODE_ENV=production
VITE_NODE_ENV=production

# 应用配置
VITE_APP_TITLE=HotGo 管理系统
VITE_APP_SHORT_NAME=HotGo
VITE_PUBLIC_PATH=/

# API配置
VITE_GLOB_API_URL=https://api.example.com
VITE_GLOB_API_URL_PREFIX=/api
VITE_GLOB_UPLOAD_URL=https://api.example.com/upload

# CDN配置
VITE_USE_CDN=true
VITE_CDN_URL=https://cdn.example.com

# 功能开关
VITE_USE_MOCK=false
VITE_USE_PWA=true

# 构建优化
VITE_BUILD_GZIP=true
VITE_BUILD_ANALYZE=false
VITE_DROP_CONSOLE=true
VITE_DROP_DEBUGGER=true

# 错误监控
VITE_SENTRY_DSN=https://xxx@sentry.io/xxx
```

### 构建脚本

```json
{
  "scripts": {
    "dev": "vite --mode development",
    "build": "vite build --mode production",
    "build:test": "vite build --mode staging",
    "build:dev": "vite build --mode development",
    "build:analyze": "cross-env VITE_BUILD_ANALYZE=true vite build",
    "preview": "vite preview",
    "preview:build": "pnpm run build && vite preview"
  }
}
```

## 本地部署

### 开发环境启动

```bash
# 1. 克隆项目
git clone https://github.com/bufanyun/hotgo.git
cd hotgo/web

# 2. 安装依赖
pnpm install

# 3. 启动开发服务器
pnpm run dev

# 4. 访问应用
# 浏览器打开 http://localhost:3100
```

### 本地构建预览

```bash
# 构建生产版本
pnpm run build

# 预览构建结果
pnpm run preview

# 或者使用其他静态服务器
npx serve dist -p 3000
```

## 服务器部署

### 1. Nginx 部署

#### 安装 Nginx

```bash
# Ubuntu/Debian
sudo apt update
sudo apt install nginx

# CentOS/RHEL
sudo yum install nginx
# 或
sudo dnf install nginx

# macOS
brew install nginx
```

#### Nginx 配置

创建配置文件 `/etc/nginx/sites-available/hotgo`：

```nginx
server {
    listen 80;
    server_name your-domain.com www.your-domain.com;
    
    # 重定向到HTTPS
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name your-domain.com www.your-domain.com;
    
    # SSL证书配置
    ssl_certificate /path/to/your/certificate.crt;
    ssl_certificate_key /path/to/your/private.key;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers ECDHE-RSA-AES128-GCM-SHA256:ECDHE-RSA-AES256-GCM-SHA384;
    ssl_prefer_server_ciphers off;
    
    # 网站根目录
    root /var/www/hotgo;
    index index.html;
    
    # Gzip压缩
    gzip on;
    gzip_comp_level 6;
    gzip_types
        text/plain
        text/css
        text/xml
        text/javascript
        application/javascript
        application/xml+rss
        application/json;
    
    # 缓存配置
    location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg|woff|woff2|ttf|eot)$ {
        expires 1y;
        add_header Cache-Control "public, immutable";
        access_log off;
    }
    
    # HTML文件不缓存
    location ~* \.html$ {
        expires -1;
        add_header Cache-Control "no-cache, no-store, must-revalidate";
    }
    
    # 前端路由配置 (SPA)
    location / {
        try_files $uri $uri/ /index.html;
        add_header X-Frame-Options "SAMEORIGIN" always;
        add_header X-Content-Type-Options "nosniff" always;
        add_header X-XSS-Protection "1; mode=block" always;
    }
    
    # API代理
    location /api/ {
        proxy_pass http://localhost:8000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # 超时配置
        proxy_connect_timeout 30s;
        proxy_send_timeout 30s;
        proxy_read_timeout 30s;
    }
    
    # 文件上传
    location /upload/ {
        proxy_pass http://localhost:8000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # 上传文件大小限制
        client_max_body_size 100M;
    }
    
    # 安全配置
    location ~ /\. {
        deny all;
        access_log off;
        log_not_found off;
    }
    
    # 日志配置
    access_log /var/log/nginx/hotgo_access.log;
    error_log /var/log/nginx/hotgo_error.log;
}
```

#### 启用配置

```bash
# 创建软链接
sudo ln -s /etc/nginx/sites-available/hotgo /etc/nginx/sites-enabled/

# 测试配置
sudo nginx -t

# 重启 Nginx
sudo systemctl restart nginx

# 设置开机自启
sudo systemctl enable nginx
```

### 2. Apache 部署

#### 安装 Apache

```bash
# Ubuntu/Debian
sudo apt install apache2

# CentOS/RHEL
sudo yum install httpd
# 或
sudo dnf install httpd
```

#### Apache 配置

创建虚拟主机配置文件：

```apache
<VirtualHost *:80>
    ServerName your-domain.com
    ServerAlias www.your-domain.com
    DocumentRoot /var/www/hotgo
    
    # 重定向到HTTPS
    Redirect permanent / https://your-domain.com/
</VirtualHost>

<VirtualHost *:443>
    ServerName your-domain.com
    ServerAlias www.your-domain.com
    DocumentRoot /var/www/hotgo
    
    # SSL配置
    SSLEngine on
    SSLCertificateFile /path/to/your/certificate.crt
    SSLCertificateKeyFile /path/to/your/private.key
    
    # 启用mod_rewrite
    RewriteEngine On
    
    # SPA路由配置
    RewriteCond %{REQUEST_FILENAME} !-f
    RewriteCond %{REQUEST_FILENAME} !-d
    RewriteRule . /index.html [L]
    
    # 缓存配置
    <FilesMatch "\.(css|js|png|jpg|jpeg|gif|ico|svg|woff|woff2|ttf|eot)$">
        ExpiresActive On
        ExpiresDefault "access plus 1 year"
        Header set Cache-Control "public, immutable"
    </FilesMatch>
    
    <FilesMatch "\.html$">
        ExpiresActive On
        ExpiresDefault "access plus 0 seconds"
        Header set Cache-Control "no-cache, no-store, must-revalidate"
    </FilesMatch>
    
    # Gzip压缩
    <IfModule mod_deflate.c>
        AddOutputFilterByType DEFLATE text/plain
        AddOutputFilterByType DEFLATE text/html
        AddOutputFilterByType DEFLATE text/xml
        AddOutputFilterByType DEFLATE text/css
        AddOutputFilterByType DEFLATE application/xml
        AddOutputFilterByType DEFLATE application/xhtml+xml
        AddOutputFilterByType DEFLATE application/rss+xml
        AddOutputFilterByType DEFLATE application/javascript
        AddOutputFilterByType DEFLATE application/x-javascript
    </IfModule>
    
    # 安全配置
    <FilesMatch "^\.">
        Require all denied
    </FilesMatch>
    
    # 日志配置
    ErrorLog /var/log/apache2/hotgo_error.log
    CustomLog /var/log/apache2/hotgo_access.log combined
</VirtualHost>
```

### 3. CDN 部署

#### 阿里云 OSS + CDN

```bash
# 1. 安装阿里云CLI工具
npm install -g @alicloud/cli

# 2. 配置认证信息
ossutil config

# 3. 构建项目
pnpm run build

# 4. 上传到OSS
ossutil cp -r dist/ oss://your-bucket-name/web/ --include="*"

# 5. 刷新CDN缓存
aliyun cdn RefreshObjectCaches --ObjectPath https://your-cdn-domain.com/
```

#### 腾讯云 COS + CDN

```bash
# 1. 安装工具
npm install -g coscmd

# 2. 配置认证
coscmd config -a your-secret-id -s your-secret-key -b your-bucket -r your-region

# 3. 上传文件
coscmd upload -r dist/ /web/

# 4. 刷新CDN
qcloud cdn purge urls --urls https://your-cdn-domain.com/
```

## 容器化部署

### Docker 部署

#### Dockerfile

```dockerfile
# 多阶段构建
FROM node:18-alpine AS builder

# 设置工作目录
WORKDIR /app

# 复制package.json文件
COPY package.json pnpm-lock.yaml ./

# 安装pnpm
RUN npm install -g pnpm

# 安装依赖
RUN pnpm install --frozen-lockfile

# 复制源代码
COPY . .

# 构建应用
RUN pnpm run build

# 生产阶段
FROM nginx:alpine

# 复制构建文件
COPY --from=builder /app/dist /usr/share/nginx/html

# 复制nginx配置
COPY nginx.conf /etc/nginx/nginx.conf

# 暴露端口
EXPOSE 80

# 启动nginx
CMD ["nginx", "-g", "daemon off;"]
```

#### nginx.conf

```nginx
user nginx;
worker_processes auto;

error_log /var/log/nginx/error.log notice;
pid /var/run/nginx.pid;

events {
    worker_connections 1024;
}

http {
    include /etc/nginx/mime.types;
    default_type application/octet-stream;
    
    log_format main '$remote_addr - $remote_user [$time_local] "$request" '
                    '$status $body_bytes_sent "$http_referer" '
                    '"$http_user_agent" "$http_x_forwarded_for"';
    
    access_log /var/log/nginx/access.log main;
    
    sendfile on;
    tcp_nopush on;
    keepalive_timeout 65;
    
    gzip on;
    gzip_comp_level 6;
    gzip_types
        text/plain
        text/css
        text/xml
        text/javascript
        application/javascript
        application/xml+rss
        application/json;
    
    server {
        listen 80;
        server_name localhost;
        root /usr/share/nginx/html;
        index index.html;
        
        # SPA路由配置
        location / {
            try_files $uri $uri/ /index.html;
        }
        
        # 静态资源缓存
        location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg|woff|woff2|ttf|eot)$ {
            expires 1y;
            add_header Cache-Control "public, immutable";
        }
        
        # HTML文件不缓存
        location ~* \.html$ {
            expires -1;
            add_header Cache-Control "no-cache, no-store, must-revalidate";
        }
    }
}
```

#### 构建和运行

```bash
# 构建镜像
docker build -t hotgo-web:latest .

# 运行容器
docker run -d \
  --name hotgo-web \
  -p 80:80 \
  hotgo-web:latest

# 查看日志
docker logs hotgo-web

# 停止容器
docker stop hotgo-web

# 删除容器
docker rm hotgo-web
```

### Docker Compose 部署

#### docker-compose.yml

```yaml
version: '3.8'

services:
  # 前端服务
  web:
    build:
      context: .
      dockerfile: Dockerfile
    container_name: hotgo-web
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./ssl:/etc/nginx/ssl:ro
      - ./logs:/var/log/nginx
    environment:
      - NODE_ENV=production
    restart: unless-stopped
    depends_on:
      - api
    networks:
      - hotgo-network

  # 后端API服务
  api:
    image: hotgo-api:latest
    container_name: hotgo-api
    ports:
      - "8000:8000"
    environment:
      - GO_ENV=production
      - DB_HOST=database
      - REDIS_HOST=redis
    restart: unless-stopped
    depends_on:
      - database
      - redis
    networks:
      - hotgo-network

  # 数据库服务
  database:
    image: mysql:8.0
    container_name: hotgo-mysql
    ports:
      - "3306:3306"
    environment:
      - MYSQL_ROOT_PASSWORD=your-root-password
      - MYSQL_DATABASE=hotgo
      - MYSQL_USER=hotgo
      - MYSQL_PASSWORD=your-password
    volumes:
      - mysql-data:/var/lib/mysql
      - ./init.sql:/docker-entrypoint-initdb.d/init.sql
    restart: unless-stopped
    networks:
      - hotgo-network

  # Redis服务
  redis:
    image: redis:7-alpine
    container_name: hotgo-redis
    ports:
      - "6379:6379"
    volumes:
      - redis-data:/data
    restart: unless-stopped
    networks:
      - hotgo-network

volumes:
  mysql-data:
  redis-data:

networks:
  hotgo-network:
    driver: bridge
```

#### 部署命令

```bash
# 启动所有服务
docker-compose up -d

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs web

# 停止服务
docker-compose down

# 重建并启动
docker-compose up -d --build
```

## CI/CD 自动化部署

### GitHub Actions

创建 `.github/workflows/deploy.yml`：

```yaml
name: Deploy to Production

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

jobs:
  test:
    runs-on: ubuntu-latest
    
    steps:
    - uses: actions/checkout@v3
    
    - name: Setup Node.js
      uses: actions/setup-node@v3
      with:
        node-version: '18'
        cache: 'pnpm'
    
    - name: Install pnpm
      uses: pnpm/action-setup@v2
      with:
        version: 8
    
    - name: Install dependencies
      run: pnpm install --frozen-lockfile
    
    - name: Run linting
      run: pnpm run lint:eslint
    
    - name: Run type check
      run: pnpm run type-check
    
    - name: Run tests
      run: pnpm run test

  build:
    needs: test
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/main'
    
    steps:
    - uses: actions/checkout@v3
    
    - name: Setup Node.js
      uses: actions/setup-node@v3
      with:
        node-version: '18'
        cache: 'pnpm'
    
    - name: Install pnpm
      uses: pnpm/action-setup@v2
      with:
        version: 8
    
    - name: Install dependencies
      run: pnpm install --frozen-lockfile
    
    - name: Build application
      run: pnpm run build
      env:
        VITE_GLOB_API_URL: ${{ secrets.API_URL }}
    
    - name: Upload build artifacts
      uses: actions/upload-artifact@v3
      with:
        name: dist
        path: dist/

  deploy:
    needs: build
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/main'
    
    steps:
    - name: Download build artifacts
      uses: actions/download-artifact@v3
      with:
        name: dist
        path: dist/
    
    - name: Deploy to server
      uses: easingthemes/ssh-deploy@v2.1.5
      env:
        SSH_PRIVATE_KEY: ${{ secrets.SSH_PRIVATE_KEY }}
        ARGS: "-avzr --delete"
        SOURCE: "dist/"
        REMOTE_HOST: ${{ secrets.REMOTE_HOST }}
        REMOTE_USER: ${{ secrets.REMOTE_USER }}
        TARGET: ${{ secrets.REMOTE_TARGET }}
    
    - name: Restart Nginx
      uses: appleboy/ssh-action@v0.1.5
      with:
        host: ${{ secrets.REMOTE_HOST }}
        username: ${{ secrets.REMOTE_USER }}
        key: ${{ secrets.SSH_PRIVATE_KEY }}
        script: |
          sudo systemctl reload nginx
```

### GitLab CI/CD

创建 `.gitlab-ci.yml`：

```yaml
stages:
  - test
  - build
  - deploy

variables:
  NODE_VERSION: "18"
  PNPM_VERSION: "8"

# 缓存配置
cache:
  key: ${CI_COMMIT_REF_SLUG}
  paths:
    - node_modules/
    - .pnpm-store/

before_script:
  - npm install -g pnpm@$PNPM_VERSION
  - pnpm config set store-dir .pnpm-store
  - pnpm install --frozen-lockfile

# 测试阶段
test:
  stage: test
  image: node:$NODE_VERSION-alpine
  script:
    - pnpm run lint:eslint
    - pnpm run type-check
    - pnpm run test
  coverage: '/Statements.*?(\d+(?:\.\d+)?)%/'
  artifacts:
    reports:
      coverage_report:
        coverage_format: cobertura
        path: coverage/cobertura-coverage.xml

# 构建阶段
build:
  stage: build
  image: node:$NODE_VERSION-alpine
  script:
    - pnpm run build
  artifacts:
    paths:
      - dist/
    expire_in: 1 hour
  only:
    - main
    - develop

# 部署到测试环境
deploy:staging:
  stage: deploy
  image: alpine:latest
  before_script:
    - apk add --no-cache rsync openssh
    - eval $(ssh-agent -s)
    - echo "$SSH_PRIVATE_KEY" | ssh-add -
    - mkdir -p ~/.ssh
    - echo "$SSH_KNOWN_HOSTS" >> ~/.ssh/known_hosts
  script:
    - rsync -avz --delete dist/ $STAGING_USER@$STAGING_HOST:$STAGING_PATH
    - ssh $STAGING_USER@$STAGING_HOST "sudo systemctl reload nginx"
  environment:
    name: staging
    url: https://staging.example.com
  only:
    - develop

# 部署到生产环境
deploy:production:
  stage: deploy
  image: alpine:latest
  before_script:
    - apk add --no-cache rsync openssh
    - eval $(ssh-agent -s)
    - echo "$SSH_PRIVATE_KEY" | ssh-add -
    - mkdir -p ~/.ssh
    - echo "$SSH_KNOWN_HOSTS" >> ~/.ssh/known_hosts
  script:
    - rsync -avz --delete dist/ $PRODUCTION_USER@$PRODUCTION_HOST:$PRODUCTION_PATH
    - ssh $PRODUCTION_USER@$PRODUCTION_HOST "sudo systemctl reload nginx"
  environment:
    name: production
    url: https://example.com
  when: manual
  only:
    - main
```

## 监控和运维

### 性能监控

#### 1. 前端性能监控

```typescript
// utils/monitor.ts
class PerformanceMonitor {
  static init() {
    // 页面加载性能监控
    window.addEventListener('load', () => {
      const navigation = performance.getEntriesByType('navigation')[0] as PerformanceNavigationTiming;
      
      const metrics = {
        dns: navigation.domainLookupEnd - navigation.domainLookupStart,
        tcp: navigation.connectEnd - navigation.connectStart,
        ssl: navigation.connectEnd - navigation.secureConnectionStart,
        ttfb: navigation.responseStart - navigation.requestStart,
        download: navigation.responseEnd - navigation.responseStart,
        domParse: navigation.domInteractive - navigation.responseEnd,
        domContentLoaded: navigation.domContentLoadedEventEnd - navigation.domContentLoadedEventStart,
        loadComplete: navigation.loadEventEnd - navigation.loadEventStart,
        total: navigation.loadEventEnd - navigation.navigationStart,
      };
      
      // 发送性能数据到监控平台
      this.sendMetrics('performance', metrics);
    });
    
    // 错误监控
    window.addEventListener('error', (event) => {
      this.sendError({
        message: event.message,
        filename: event.filename,
        lineno: event.lineno,
        colno: event.colno,
        stack: event.error?.stack,
      });
    });
    
    // Promise错误监控
    window.addEventListener('unhandledrejection', (event) => {
      this.sendError({
        message: 'Unhandled Promise Rejection',
        reason: event.reason,
      });
    });
  }
  
  static sendMetrics(type: string, data: any) {
    // 发送到监控平台
    fetch('/api/monitor/metrics', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ type, data, timestamp: Date.now() }),
    }).catch(console.error);
  }
  
  static sendError(error: any) {
    // 发送错误信息
    fetch('/api/monitor/errors', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ ...error, timestamp: Date.now() }),
    }).catch(console.error);
  }
}

// 在main.ts中初始化
PerformanceMonitor.init();
```

#### 2. 健康检查

```bash
# 创建健康检查脚本
cat > /usr/local/bin/hotgo-health-check.sh << 'EOF'
#!/bin/bash

# 检查网站可访问性
STATUS=$(curl -o /dev/null -s -w "%{http_code}\n" https://your-domain.com)

if [ $STATUS -eq 200 ]; then
  echo "✅ Website is healthy (HTTP $STATUS)"
  exit 0
else
  echo "❌ Website is down (HTTP $STATUS)"
  # 发送告警通知
  curl -X POST "https://api.your-alert-service.com/alert" \
    -H "Content-Type: application/json" \
    -d '{"message":"HotGo website is down","status":'$STATUS'}'
  exit 1
fi
EOF

chmod +x /usr/local/bin/hotgo-health-check.sh

# 添加到crontab，每分钟检查一次
echo "* * * * * /usr/local/bin/hotgo-health-check.sh" | crontab -
```

### 日志管理

#### Nginx 日志分析

```bash
# 安装GoAccess
sudo apt install goaccess

# 分析访问日志
goaccess /var/log/nginx/hotgo_access.log -o /var/www/html/report.html --log-format=COMBINED

# 实时监控
goaccess /var/log/nginx/hotgo_access.log -o /var/www/html/report.html --log-format=COMBINED --real-time-html
```

#### 日志轮转配置

```bash
# 创建logrotate配置
cat > /etc/logrotate.d/hotgo << 'EOF'
/var/log/nginx/hotgo_*.log {
    daily
    missingok
    rotate 52
    compress
    delaycompress
    notifempty
    create 644 www-data www-data
    postrotate
        systemctl reload nginx
    endscript
}
EOF
```

## 安全配置

### SSL/TLS 配置

```bash
# 使用Let's Encrypt获取免费SSL证书
sudo apt install certbot python3-certbot-nginx

# 获取证书
sudo certbot --nginx -d your-domain.com -d www.your-domain.com

# 自动续期
echo "0 12 * * * /usr/bin/certbot renew --quiet" | sudo crontab -
```

### 安全头配置

```nginx
# 添加到nginx配置中
add_header X-Frame-Options "SAMEORIGIN" always;
add_header X-Content-Type-Options "nosniff" always;
add_header X-XSS-Protection "1; mode=block" always;
add_header Referrer-Policy "strict-origin-when-cross-origin" always;
add_header Content-Security-Policy "default-src 'self'; script-src 'self' 'unsafe-inline' 'unsafe-eval'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; font-src 'self' data:; connect-src 'self' https:; frame-ancestors 'self';" always;
add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;
```

## 故障排查

### 常见问题

#### 1. 白屏问题

```bash
# 检查nginx配置
sudo nginx -t

# 检查文件权限
ls -la /var/www/hotgo/

# 查看nginx错误日志
sudo tail -f /var/log/nginx/hotgo_error.log

# 检查文件是否存在
ls -la /var/www/hotgo/index.html
```

#### 2. API请求失败

```bash
# 检查API服务状态
curl -I http://localhost:8000/api/health

# 查看nginx代理配置
nginx -T | grep -A 10 "location /api/"

# 检查防火墙
sudo ufw status
```

#### 3. 静态资源404

```bash
# 检查资源文件
ls -la /var/www/hotgo/assets/

# 检查nginx配置
grep -A 5 "location.*\.\(js\|css\)" /etc/nginx/sites-available/hotgo
```

### 恢复流程

```bash
# 1. 备份当前版本
sudo cp -r /var/www/hotgo /var/www/hotgo.backup.$(date +%Y%m%d_%H%M%S)

# 2. 恢复到上一个版本
sudo rm -rf /var/www/hotgo
sudo cp -r /var/www/hotgo.backup.20231201_120000 /var/www/hotgo

# 3. 重启服务
sudo systemctl reload nginx

# 4. 验证恢复
curl -I https://your-domain.com
```

---

这份文档涵盖了 HotGo 2.0 前端项目的完整部署流程，从本地开发到生产环境部署，包括了各种部署方式和运维监控。在实际使用中，请根据具体的环境和需求进行调整。






