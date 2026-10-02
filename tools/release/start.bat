@echo off
rem mini-im 启动脚本：先切到脚本所在目录，再拉起同目录下的可执行文件，
rem 因此可以在任意位置调用，也可以直接双击。
rem
rem   start.bat                使用同目录下的 config.yaml
rem   start.bat -c my.yaml     附加参数会原样透传给 mini-im
rem
rem 配置也可以用环境变量覆盖，例如在 PowerShell 中：
rem   $env:IM_RUN_MODE='release'; $env:IM_JWT_SECRET='xxx'; .\start.bat

cd /d "%~dp0"
mini-im.exe %*
