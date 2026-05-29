@echo off
chcp 65001 >nul
echo ========================================
echo   网约车项目 - 一键启动脚本
echo ========================================
echo.

:: 获取脚本所在目录
cd /d "%~dp0"

:: 第一步：结束所有相关进程
echo [1/3] 正在结束所有相关进程...
taskkill /F /IM admin-srv.exe /T >nul 2>&1
taskkill /F /IM driver-srv.exe /T >nul 2>&1
taskkill /F /IM user-srv.exe /T >nul 2>&1
taskkill /F /IM user-api.exe /T >nul 2>&1
echo       已结束所有相关进程
echo.

:: 第二步：启动后端服务
echo [2/3] 正在启动后端服务...
cd admin-srv
start "admin-srv" cmd /c "admin-srv.exe -f admin-srv.yaml"
cd ..

cd driver-srv
start "driver-srv" cmd /c "driver-srv.exe -f driver-srv.yaml"
cd ..

cd user-srv
start "user-srv" cmd /c "user-srv.exe -f user-srv.yaml"
cd ..

cd user-api
start "user-api" cmd /c "user-api.exe -f user-api.yaml"
cd ..

echo       后端服务启动中...
echo.

:: 第三步：启动前端
echo [3/3] 正在启动前端服务...
cd user-web
:: 尝试使用 Python 启动静态服务器
start "user-web" cmd /c "python -m http.server 8080"
cd ..

echo.
echo ========================================
echo   启动完成！
echo ========================================
echo.
echo   服务地址：
echo   - admin-srv:  http://localhost:8080 (admin管理后台)
echo   - driver-srv: http://localhost:8081
echo   - user-srv:   http://localhost:8082
echo   - user-api:   http://localhost:8083
echo   - user-web:   http://localhost:8084
echo   - 前端页面:   http://localhost:8080 (user-web)
echo ========================================
pause
