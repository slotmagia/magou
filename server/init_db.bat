@echo off
echo 正在初始化数据库...
echo.

:: 执行数据库初始化脚本
mysql -u root -p -e "source init_database.sql"

if %errorlevel% == 0 (
    echo.
    echo 数据库初始化成功！
    echo.
    echo 默认管理员账号：
    echo 用户名: admin
    echo 密码: 123456
    echo.
) else (
    echo.
    echo 数据库初始化失败！
    echo 请检查MySQL服务是否启动以及用户名密码是否正确
    echo.
)

pause 