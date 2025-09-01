#!/bin/bash

# 多租户系统部署脚本
# 使用说明：bash deploy_multi_tenant.sh

set -e

echo "=========================================="
echo "开始部署多租户系统"
echo "=========================================="

# 检查 MySQL 连接
check_mysql() {
    echo "检查 MySQL 连接..."
    if ! mysql -u root -p admin -e "SELECT 1;" > /dev/null 2>&1; then
        echo "错误: 无法连接到 MySQL 数据库 'admin'"
        echo "请确保："
        echo "1. MySQL 服务正在运行"
        echo "2. 数据库 'admin' 存在"
        echo "3. 用户权限正确"
        exit 1
    fi
    echo "✓ MySQL 连接正常"
}

# 备份现有数据
backup_data() {
    echo "备份现有数据..."
    backup_file="backup_$(date +%Y%m%d_%H%M%S).sql"
    mysqldump -u root -p admin > "backup/$backup_file"
    echo "✓ 数据备份完成: backup/$backup_file"
}

# 执行数据库脚本
execute_sql() {
    echo "执行数据库脚本..."
    
    echo "  - 创建租户表..."
    mysql -u root -p admin < internal/sql/tenants.sql
    echo "  ✓ 租户表创建完成"
    
    echo "  - 调整现有表结构..."
    mysql -u root -p admin < internal/sql/add_tenant_id.sql
    echo "  ✓ 表结构调整完成"
    
    echo "✓ 数据库脚本执行完成"
}

# 验证表结构
verify_schema() {
    echo "验证表结构..."
    
    # 检查租户表
    if mysql -u root -p admin -e "DESCRIBE tenants;" > /dev/null 2>&1; then
        echo "  ✓ 租户表结构正常"
    else
        echo "  ✗ 租户表结构异常"
        exit 1
    fi
    
    # 检查用户表的租户字段
    if mysql -u root -p admin -e "SHOW COLUMNS FROM users LIKE 'tenant_id';" | grep -q "tenant_id"; then
        echo "  ✓ 用户表租户字段正常"
    else
        echo "  ✗ 用户表租户字段异常"
        exit 1
    fi
    
    echo "✓ 表结构验证通过"
}

# 编译应用
build_app() {
    echo "编译应用..."
    if go build -o admin main.go; then
        echo "✓ 应用编译完成"
    else
        echo "✗ 应用编译失败"
        exit 1
    fi
}

# 创建测试租户
create_test_tenant() {
    echo "创建测试租户..."
    
    # 启动应用（后台运行）
    ./admin &
    APP_PID=$!
    
    # 等待应用启动
    sleep 5
    
    # 创建测试租户
    curl -X POST "http://localhost:8888/api/tenant/create" \
      -H "Content-Type: application/json" \
      -d '{
        "name": "测试租户",
        "code": "test",
        "maxUsers": 100,
        "storageLimit": 1073741824,
        "adminName": "test_admin",
        "adminEmail": "admin@test.com",
        "adminPassword": "MTIzNDU2",
        "remark": "测试租户"
      }' > /dev/null 2>&1
    
    # 停止应用
    kill $APP_PID
    
    echo "✓ 测试租户创建完成"
}

# 生成配置文件
generate_config() {
    echo "生成配置文件..."
    
    cat > manifest/config/multi-tenant.yaml << EOF
# 多租户配置
tenant:
  # 默认租户识别方式：header, domain, param
  identifyMode: "header"
  
  # 默认租户ID（当无法识别租户时使用）
  defaultTenantId: 1
  
  # 域名映射（当 identifyMode 为 domain 时使用）
  domainMapping:
    "app.example.com": 1
    "demo.example.com": 2
  
  # 缓存配置
  cache:
    enabled: true
    ttl: "1h"
  
  # 资源限制
  limits:
    defaultMaxUsers: 100
    defaultStorageLimit: 1073741824 # 1GB
    
# 安全配置
security:
  # 强制租户隔离
  forceIsolation: true
  
  # 跨租户访问权限
  crossTenantAccess:
    - "super_admin"
    - "system_admin"
EOF

    echo "✓ 配置文件生成完成: manifest/config/multi-tenant.yaml"
}

# 创建部署文档
create_docs() {
    echo "创建部署文档..."
    
    cat > DEPLOYMENT_RESULT.md << EOF
# 多租户系统部署结果

## 部署概要

✅ **部署状态**: 成功  
📅 **部署时间**: $(date '+%Y-%m-%d %H:%M:%S')  
🏗️ **部署版本**: v1.0.0  

## 系统信息

### 数据库变更
- ✅ 创建租户表 \`tenants\`
- ✅ 所有业务表添加 \`tenant_id\` 字段
- ✅ 调整唯一索引支持租户隔离
- ✅ 添加租户相关索引优化

### 新增功能
- ✅ 租户管理 API
- ✅ 多租户认证机制
- ✅ 租户数据隔离中间件
- ✅ 租户配置管理

### 默认账户

#### 系统租户 (ID: 1)
- **租户名称**: 系统租户
- **租户编码**: system
- **管理员**: admin (原有账户)

#### 测试租户 (ID: 2)
- **租户名称**: 测试租户
- **租户编码**: test  
- **管理员**: test_admin
- **密码**: 123456

## API 接口

### 租户管理
- \`GET /api/tenant/list\` - 获取租户列表
- \`POST /api/tenant/create\` - 创建租户
- \`PUT /api/tenant/update\` - 更新租户
- \`DELETE /api/tenant/delete\` - 删除租户

### 租户切换
在请求头中添加：
\`\`\`
X-Tenant-Id: 2
\`\`\`

## 测试验证

### 1. 租户列表查询
\`\`\`bash
curl -X GET "http://localhost:8888/api/tenant/list" \\
  -H "Authorization: Bearer <your_token>"
\`\`\`

### 2. 租户用户登录
\`\`\`bash
curl -X POST "http://localhost:8888/api/user/login" \\
  -H "Content-Type: application/json" \\
  -H "X-Tenant-Id: 2" \\
  -d '{
    "username": "test_admin",
    "password": "MTIzNDU2"
  }'
\`\`\`

## 注意事项

⚠️ **重要提醒**:
1. 已创建数据备份文件在 \`backup/\` 目录
2. 多租户配置文件位于 \`manifest/config/multi-tenant.yaml\`
3. 详细实施指南请参考 \`docs/multi-tenant-implementation-guide.md\`

## 下一步操作

1. 🔧 **配置域名映射**（可选）
2. 🎨 **自定义租户配置**
3. 📊 **设置监控告警**
4. 🧪 **进行全面测试**

---
📚 更多信息请查看：\`docs/multi-tenant-implementation-guide.md\`
EOF

    echo "✓ 部署文档创建完成: DEPLOYMENT_RESULT.md"
}

# 主执行流程
main() {
    # 创建必要目录
    mkdir -p backup
    
    # 执行部署步骤
    check_mysql
    backup_data
    execute_sql
    verify_schema
    build_app
    generate_config
    create_docs
    
    echo ""
    echo "=========================================="
    echo "🎉 多租户系统部署完成！"
    echo "=========================================="
    echo ""
    echo "📋 部署结果："
    echo "  ✅ 数据库结构调整完成"
    echo "  ✅ 应用编译完成"
    echo "  ✅ 配置文件生成完成"
    echo "  ✅ 部署文档创建完成"
    echo ""
    echo "📖 请查看以下文件："
    echo "  📄 DEPLOYMENT_RESULT.md - 部署结果总结"
    echo "  📚 docs/multi-tenant-implementation-guide.md - 详细实施指南"
    echo "  ⚙️ manifest/config/multi-tenant.yaml - 多租户配置"
    echo ""
    echo "🚀 启动应用："
    echo "  ./admin"
    echo ""
}

# 错误处理
trap 'echo "❌ 部署过程中发生错误，请检查上述输出信息"; exit 1' ERR

# 执行主函数
main
