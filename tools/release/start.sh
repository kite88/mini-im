#!/usr/bin/env bash
#
# mini-im 启动脚本：先切到脚本所在目录，再拉起同目录下的可执行文件，
# 因此可以在任意位置调用，不需要先 cd 进解压目录。
#
#   ./start.sh                 # 使用同目录下的 config.yaml
#   ./start.sh -c my.yaml      # 附加参数会原样透传给 mini-im
#
# 配置也可以用环境变量覆盖，例如：
#   IM_RUN_MODE=release IM_JWT_SECRET=xxx ./start.sh

set -euo pipefail
cd "$(dirname "$0")"
exec ./mini-im "$@"
