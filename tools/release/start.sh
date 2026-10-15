#!/usr/bin/env bash
#
# mini-im 启动脚本：先切到脚本所在目录，再拉起同目录下的可执行文件，
# 因此可以在任意位置调用，不需要先 cd 进解压目录。
#
#   ./start.sh                 # 使用同目录下的 config.yaml
#   ./start.sh -c my.yaml      # 附加参数会原样透传给 mini-im
#
# 首次运行归档里没有 config.yaml，服务会先进初始化配置向导：
# 浏览器打开 http://localhost:8280/setup 填写 PostgreSQL / Redis 并保存，
# 配置写回后业务服务自动启动。想重新配置：SETUP_FORCE=1 ./start.sh
# （原配置备份为 config.yaml.bak）。
#
# 配置也可以用环境变量覆盖，例如：
#   IM_RUN_MODE=release IM_JWT_SECRET=xxx ./start.sh

set -euo pipefail
cd "$(dirname "$0")"

# 版本号由打包脚本写入 version.txt；源码目录直接运行时没有这个文件
version=unknown
if [ -f version.txt ]; then
    version=$(tr -d '\r' < version.txt)
fi

echo "[mini-im] 版本：$version"
if [ ! -f config.yaml ]; then
    echo "[mini-im] 未检测到 config.yaml，将进入初始化配置向导："
    echo "[mini-im]   浏览器打开 http://localhost:8280/setup"
fi

exec ./mini-im "$@"
