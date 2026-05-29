@echo off
chcp 65001 >nul
echo ========================================
echo   网约车项目 - 停止所有服务
echo ========================================
echo.

echo 正在结束所有相关进程...

:: 结束所有 Go 服务进程
taskkill /F /IM admin-srv.exe /T >nul 2>&1
taskkill /F /IM driver-srv.exe /T >nul 2>&1
taskkill /F /IM user-srv.exe /T >nul 2>&1
taskkill /F /IM user-api.exe /T >nul 2>&1

:: 结束 Python HTTP 服务器 (如果前端用 Python 启动)
taskkill /F /IM python.exe /T >nul 2>&1

echo.
echo 所有服务已停止！
pause
