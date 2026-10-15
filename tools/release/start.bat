@echo off
rem mini-im launcher ---------------------------------------------------------
rem Switches to this script's directory first, then starts the executable next
rem to it, so it works from anywhere (double-click friendly).
rem
rem   start.bat                uses config.yaml in the same directory
rem   start.bat -c my.yaml     extra arguments are passed straight to mini-im
rem
rem First run: there is no config.yaml in the archive, so mini-im starts the
rem setup wizard instead of the service -- open http://localhost:8280/setup,
rem fill in PostgreSQL / Redis and save; the service then starts by itself.
rem To redo the setup later, set SETUP_FORCE=1 and run this script again
rem (the old config.yaml is backed up as config.yaml.bak).
rem
rem Configuration can also come from IM_* environment variables, e.g. in
rem PowerShell:  $env:IM_RUN_MODE='release'; $env:IM_JWT_SECRET='xxx'; .\start.bat
rem
rem This file is intentionally ASCII-only: after "chcp 65001" a batch file
rem containing non-ASCII text can make cmd.exe misparse the following lines.
rem The window deliberately stays open after the process exits, so a failed
rem start does not just flash by.
rem ---------------------------------------------------------------------------

rem Switch the console to UTF-8 so the application's Chinese log output is readable
chcp 65001 >nul 2>&1

setlocal
cd /d "%~dp0"

rem 版本号由打包脚本写入 version.txt；源码目录直接运行时没有这个文件
set "VERSION=unknown"
if exist "version.txt" set /p VERSION=<version.txt

echo [mini-im] dir:  %CD%
echo [mini-im] version: %VERSION%
if not exist "config.yaml" (
    echo [mini-im] no config.yaml yet: the setup wizard will start.
    echo [mini-im] open http://localhost:8280/setup to finish the initial setup.
) else (
    echo [mini-im] conf: config.yaml  [PostgreSQL / Redis / port can be overridden by IM_* env vars]
)
echo [mini-im] starting; a reachable PostgreSQL and Redis are required, errors are shown below.
echo.

mini-im.exe %*
set "exit_code=%ERRORLEVEL%"

echo.
if "%exit_code%"=="0" (
    echo [mini-im] exited normally.
) else (
    echo [mini-im] process exited with code %exit_code%.
    echo [mini-im] common causes: PostgreSQL / Redis not running or wrong connection settings;
    echo [mini-im]               database "mini-im" not created yet; HTTP port 8280 already in use.
)
echo.
echo [mini-im] Press any key to close this window...
pause >nul
