#!/bin/bash

echo "=========================================="
echo "多租户系统编译测试"
echo "=========================================="

# 检查 Go 版本
echo "Go 版本信息："
go version

echo ""
echo "开始编译测试..."

# 编译项目
if go build -o admin_test main.go; then
    echo "✅ 编译成功！"
    
    # 检查生成的可执行文件
    if [ -f "admin_test" ]; then
        echo "✅ 可执行文件生成成功：admin_test"
        file_size=$(du -h admin_test | cut -f1)
        echo "📦 文件大小：$file_size"
        
        # 清理测试文件
        rm -f admin_test
        echo "🧹 清理测试文件完成"
    else
        echo "❌ 可执行文件未生成"
        exit 1
    fi
else
    echo "❌ 编译失败！请检查错误信息："
    echo ""
    exit 1
fi

echo ""
echo "=========================================="
echo "🎉 编译测试通过！"
echo "=========================================="
echo ""
echo "✅ 多租户相关文件编译正常"
echo "✅ 依赖关系正确"
echo "✅ 代码语法无误"
echo ""
echo "📋 下一步操作："
echo "1. 执行数据库脚本：bash deploy_multi_tenant.sh"
echo "2. 启动应用测试：./admin"
echo "3. 测试多租户功能"
echo ""
