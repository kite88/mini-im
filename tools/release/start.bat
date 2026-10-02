@echo off
rem mini-im 启动脚本：先切到脚本所在目录，再拉起同目录下的可执行文件，
rem 因此可以在任意位置调用，也可以直接双击。
rem
rem   start.bat                使用同目录下的 config.yaml
rem   start.bat -c my.yaml     附加参数会原样透传给 mini-im
rem
rem 配置也可以用环境变量覆盖，例如在 PowerShell 中：
rem   $env:IM_RUN_MODE='release'; $env:IM_JWT_SECRET='xxx'; .\start.bat
rem
rem 双击运行时窗口不会在进程退出后立刻关闭：失败时报错会留在屏幕上，
rem 并显示退出码与常见原因，避免只看到一闪而过的黑框。

rem 切到 UTF-8 代码页，否则中文日志在默认 GBK 控制台下会显示成乱码
chcp 65001 >nul 2>&1

setlocal
cd /d "%~dp0"

echo [mini-im] 目录: %CD%
echo [mini-im] 配置: config.yaml（数据库 / Redis / 端口可用 IM_* 环境变量覆盖）
echo [mini-im] 启动需要可用的 PostgreSQL 与 Redis，失败原因见下方日志。
echo.

mini-im.exe %*
set "exit_code=%ERRORLEVEL%"

echo.
if "%exit_code%"=="0" (
    echo [mini-im] 已正常退出。
) else (
    echo [mini-im] 进程退出，退出码 %exit_code%。
    echo [mini-im] 常见原因：PostgreSQL / Redis 未启动或连接配置不对；
    echo [mini-im]          数据库 mini-im 尚未创建；HTTP 端口 2580 被占用。
)
echo.
echo [mini-im] 按任意键关闭窗口...
pause >nul
