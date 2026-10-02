@echo off
rem mini-im launcher ---------------------------------------------------------
rem Switches to this script's directory first, then starts the executable next
rem to it, so it works from anywhere (double-click friendly).
rem
rem   start.bat                uses config.yaml in the same directory
rem   start.bat -c my.yaml     extra arguments are passed straight to mini-im
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

echo [mini-im] dir:  %CD%
echo [mini-im] conf: config.yaml  [PostgreSQL / Redis / port can be overridden by IM_* env vars]
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
    echo [mini-im]               database "mini-im" not created yet; HTTP port 2580 already in use.
)
echo.
echo [mini-im] Press any key to close this window...
pause >nul
