#!/bin/bash

# 多租户管理系统 数据库初始化脚本 (Linux/macOS)
# 执行前确保MySQL服务已启动

echo "================================"
echo "  多租户管理系统 数据库初始化"
echo "================================"
echo

# 设置颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# MySQL连接参数（请根据实际情况修改）
MYSQL_HOST=${MYSQL_HOST:-localhost}
MYSQL_PORT=${MYSQL_PORT:-3306}
MYSQL_USER=${MYSQL_USER:-root}
MYSQL_DB=${MYSQL_DB:-admin}

# 提示输入密码
echo -e "${BLUE}[信息]${NC} 正在检查MySQL连接..."
read -s -p "请输入MySQL密码: " MYSQL_PASSWORD
echo
echo

# 检查MySQL连接
mysql -h"$MYSQL_HOST" -P"$MYSQL_PORT" -u"$MYSQL_USER" -p"$MYSQL_PASSWORD" -e "SELECT 1" > /dev/null 2>&1

if [ $? -ne 0 ]; then
    echo -e "${RED}[错误]${NC} MySQL连接失败！"
    echo "请检查："
    echo "1. MySQL服务是否启动"
    echo "2. 用户名密码是否正确"
    echo "3. 网络连接是否正常"
    exit 1
fi

echo -e "${GREEN}[成功]${NC} MySQL连接正常！"
echo

# 检查SQL文件是否存在
if [ ! -f "SQL/admin.sql" ]; then
    echo -e "${RED}[错误]${NC} SQL/admin.sql 文件不存在！"
    echo "请确保在项目根目录执行此脚本"
    exit 1
fi

if [ ! -f "init_data.sql" ]; then
    echo -e "${RED}[错误]${NC} init_data.sql 文件不存在！"
    echo "请确保在项目根目录执行此脚本"
    exit 1
fi

echo -e "${BLUE}[信息]${NC} 正在执行数据库结构初始化..."

# 1. 执行表结构创建
mysql -h"$MYSQL_HOST" -P"$MYSQL_PORT" -u"$MYSQL_USER" -p"$MYSQL_PASSWORD" < SQL/admin.sql

if [ $? -ne 0 ]; then
    echo -e "${RED}[错误]${NC} 数据库结构初始化失败！"
    echo "请检查 SQL/admin.sql 文件内容"
    exit 1
fi

echo -e "${GREEN}[成功]${NC} 数据库结构初始化完成！"
echo

echo -e "${BLUE}[信息]${NC} 正在执行初始化数据导入..."

# 2. 执行初始化数据插入
mysql -h"$MYSQL_HOST" -P"$MYSQL_PORT" -u"$MYSQL_USER" -p"$MYSQL_PASSWORD" "$MYSQL_DB" < init_data.sql

if [ $? -ne 0 ]; then
    echo -e "${RED}[错误]${NC} 初始化数据导入失败！"
    echo "请检查 init_data.sql 文件内容"
    exit 1
fi

echo -e "${GREEN}[成功]${NC} 初始化数据导入完成！"
echo

echo "================================"
echo "   初始化完成！"
echo "================================"
echo

echo -e "${GREEN}🎉 恭喜！多租户管理系统数据库初始化成功！${NC}"
echo

echo -e "${YELLOW}📋 默认账号信息：${NC}"
echo "┌─────────────────────────────────┐"
echo "│  系统管理员账号                 │"
echo "│  用户名: admin                  │"
echo "│  密码: 123456                   │"
echo "│  租户: system                   │"
echo "└─────────────────────────────────┘"
echo
echo "┌─────────────────────────────────┐"
echo "│  演示管理员账号                 │"
echo "│  用户名: demo_admin             │"
echo "│  密码: 123456                   │"
echo "│  租户: demo                     │"
echo "└─────────────────────────────────┘"
echo

echo -e "${BLUE}🌐 访问地址:${NC} http://localhost:8000"
echo -e "${BLUE}📚 API文档:${NC} http://localhost:8000/swagger"
echo

echo -e "${YELLOW}📝 注意事项：${NC}"
echo "1. 首次登录建议修改默认密码"
echo "2. 系统管理员拥有所有权限"
echo "3. 演示账号仅用于功能演示"
echo "4. 生产环境请删除演示数据"
echo

echo -e "${GREEN}✅ 数据库初始化完成，可以启动应用程序了！${NC}"

