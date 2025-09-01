@echo off
chcp 65001 > nul
echo ================================
echo   多租户管理系统 数据库初始化
echo ================================
echo.

echo [信息] 正在检查MySQL连接...

:: 设置MySQL连接参数（请根据实际情况修改）
set MYSQL_HOST=localhost
set MYSQL_PORT=3306
set MYSQL_USER=root
set /p MYSQL_PASSWORD=请输入MySQL密码: 

echo.
echo [信息] 正在执行数据库结构初始化...

:: 1. 执行表结构创建
mysql -h%MYSQL_HOST% -P%MYSQL_PORT% -u%MYSQL_USER% -p%MYSQL_PASSWORD% < SQL/admin.sql

if %errorlevel% neq 0 (
    echo [错误] 数据库结构初始化失败！
    echo 请检查：
    echo 1. MySQL服务是否启动
    echo 2. 用户名密码是否正确
    echo 3. SQL/admin.sql文件是否存在
    echo.
    pause
    exit /b 1
)

echo [成功] 数据库结构初始化完成！
echo.
echo [信息] 正在执行初始化数据导入...

:: 2. 执行初始化数据插入
mysql -h%MYSQL_HOST% -P%MYSQL_PORT% -u%MYSQL_USER% -p%MYSQL_PASSWORD% admin < init_data.sql

if %errorlevel% neq 0 (
    echo [错误] 初始化数据导入失败！
    echo 请检查 init_data.sql 文件内容
    echo.
    pause
    exit /b 1
)

echo [成功] 初始化数据导入完成！
echo.
echo ================================
echo   初始化完成！
echo ================================
echo.
echo 🎉 恭喜！多租户管理系统数据库初始化成功！
echo.
echo 📋 默认账号信息：
echo ┌─────────────────────────────────┐
echo │  系统管理员账号                 │
echo │  用户名: admin                  │
echo │  密码: 123456                   │
echo │  租户: system                   │
echo └─────────────────────────────────┘
echo.
echo ┌─────────────────────────────────┐
echo │  演示管理员账号                 │
echo │  用户名: demo_admin             │
echo │  密码: 123456                   │
echo │  租户: demo                     │
echo └─────────────────────────────────┘
echo.
echo 🌐 访问地址: http://localhost:8000
echo 📚 API文档: http://localhost:8000/swagger
echo.
echo 📝 注意事项：
echo 1. 首次登录建议修改默认密码
echo 2. 系统管理员拥有所有权限
echo 3. 演示账号仅用于功能演示
echo 4. 生产环境请删除演示数据
echo.

pause

