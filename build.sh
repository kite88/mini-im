#!/usr/bin/env bash
#
# 交叉编译 mini-im 并打包为可分发的归档，产物落在 dist/。
#
# 每个归档是一个顶层目录 mini-im-<version>-<os>-<arch>/，解压后直接可用：
#
#   mini-im-<version>-<os>-<arch>/
#   ├── mini-im(.exe)        可执行文件（非 Windows 平台带 0755 权限）
#   ├── start.sh             启动脚本（Windows 为 start.bat，带 0755 权限）
#   ├── version.txt          版本号，启动脚本启动时打印它
#   ├── config.example.yaml  配置模板（首次启动进入初始化向导时按它写回 config.yaml）
#   ├── README.md / README.en.md / LICENSE
#   └── web/                 前端静态页（后端直接托管，必须与可执行文件放在一起）
#
# 归档里刻意不带 config.yaml：首次运行会进入初始化配置向导（浏览器打开
# /setup 填写数据库 / Redis 后写回），免得把开发用的默认配置（含默认 JWT 密钥）
# 发到用户机器上。
#
# 另有 checksums.txt 记录全部归档的 SHA256（LF 换行，Linux 下可直接 sha256sum -c）。
#
# 用法: ./build.sh [-v Version] [-o OutDir]

set -euo pipefail

usage() {
    cat <<'EOF'
交叉编译 mini-im 并打包为可分发的归档。

用法: ./build.sh [-v Version] [-o OutDir]

选项:
  -v Version  版本号，形如 v1.1.0（可带 -rc1 后缀）；省略时取最近的 git 标签
  -o OutDir   输出目录（默认: dist，相对本脚本所在目录）
  -h          显示本帮助

平台矩阵: windows/amd64 windows/arm64 linux/amd64 linux/arm64 darwin/amd64 darwin/arm64
每个平台一个归档（Windows 为 .zip，其余为 .tar.gz），外加 checksums.txt。
归档名与包内顶层目录都带版本号，包内另有 version.txt。
EOF
}

out_dir=dist
version=''
while getopts ':ho:v:' opt; do
    case "$opt" in
        h) usage; exit 0 ;;
        o) out_dir=$OPTARG ;;
        v) version=$OPTARG ;;
        '?') usage >&2; exit 2 ;;
    esac
done

cd "$(dirname "$0")"

# 版本号决定归档名与包内 version.txt，必须一眼看得出是正式版本：
# 允许 vX.Y.Z，也允许 vX.Y.Z-rc1 / vX.Y.Z+build 这类后缀，
# 但 git describe 在标签之后有新提交时会给出 v1.0.0-3-gabc123，这种既不像版本
# 也不像标签，一旦打进归档名就会造成困惑，直接拦下要求显式指定。
if [ -z "$version" ]; then
    version=$(git describe --tags --always 2>/dev/null || true)
fi
if [ -z "$version" ]; then
    echo "未能确定版本号：请显式指定 -v vX.Y.Z（例如 -v v1.1.0）" >&2
    exit 2
fi
if ! printf '%s' "$version" | grep -Eq '^v[0-9]+\.[0-9]+\.[0-9]+(-[0-9A-Za-z.-]+)?(\+[0-9A-Za-z.-]+)?$'; then
    echo "版本号 '$version' 不是 vX.Y.Z 形式（可带 -rc1 这类后缀）：请用 -v 显式指定" >&2
    exit 2
fi
# git describe 在标签之后还有提交时给出 v1.0.0-3-gabc123：它符合上面的后缀规则，
# 但意思是「离标签还有 3 个提交」，不是版本号，打进归档名会让人以为是正式版本
if printf '%s' "$version" | grep -Eq -- '-g[0-9a-f]{6,}$'; then
    echo "版本号 '$version' 是 git describe 的距离格式（标签之后还有提交）：请显式指定 -v vX.Y.Z" >&2
    exit 2
fi

if [ -z "$out_dir" ]; then
    echo "OutDir 不能为空" >&2
    exit 2
fi
case "$out_dir" in
    /*) ;;
    *) out_dir=$PWD/$out_dir ;;
esac

# 发布包的硬前提：页面与配置模板缺一个，包就是残的（接口能起、页面打不开，
# 或进不去初始化向导），这里提前拦下而不是产出坏包。
for required in web/index.html config.example.yaml LICENSE; do
    if [ ! -f "$required" ]; then
        echo "缺少 $required，无法打包" >&2
        exit 1
    fi
done

# 目标矩阵，'os/arch' 形式；增删一行即可调整发布平台。
targets='windows/amd64 windows/arm64 linux/amd64 linux/arm64 darwin/amd64 darwin/arm64'

# 交叉编译不需要 cgo（纯 Go 依赖，含 PostgreSQL / Redis 驱动）
export CGO_ENABLED=0

mkdir -p "$out_dir"
# 只删除输出目录顶层的文件，不递归，避免误伤目录外的东西
find "$out_dir" -maxdepth 1 -type f -delete

staging=$(mktemp -d)
trap 'rm -rf "$staging"' EXIT

# Windows 归档：优先用 zip；Git Bash 与精简环境里经常没有 zip，退回 Python
# 标准库（产出同样是标准 zip）。两者都没有就报错退出，绝不产出半个包。
archive_win() {
    local pkg=$1 archive=$2
    if command -v zip >/dev/null 2>&1; then
        (cd "$staging" && zip -qr "$archive" "$pkg")
        return
    fi
    local py
    # python 排在 python3 前面：Windows 上 python3 常常是 Microsoft Store 的存根，
    # command -v 能查到却执行不了（exit 49 且没有任何输出），所以先真跑一次再采用
    for py in python python3 py; do
        if command -v "$py" >/dev/null 2>&1 && "$py" -c 'import zipfile' >/dev/null 2>&1; then
            (cd "$staging" && "$py" -c '
import os, sys, zipfile
src, dst = sys.argv[1], sys.argv[2]
with zipfile.ZipFile(dst, "w", zipfile.ZIP_DEFLATED) as z:
    for root, _dirs, files in os.walk(src):
        for name in sorted(files):
            path = os.path.join(root, name)
            arc = os.path.join(src, os.path.relpath(path, src)).replace(os.sep, "/")
            z.write(path, arc)
' "$pkg" "$archive")
            return
        fi
    done
    echo "缺少 zip 与 python，无法生成 Windows 归档 $archive" >&2
    exit 1
}

count=0
for target in $targets; do
    goos=${target%%/*}
    goarch=${target#*/}

    # 可执行文件在包内保持短名（用户手动运行时好敲），目录名与归档名带版本号
    bin_name=mini-im
    if [ "$goos" = windows ]; then
        bin_name=mini-im.exe
    fi
    pkg_name="mini-im-$version-$goos-$goarch"
    if [ "$goos" = windows ]; then
        archive_name="$pkg_name.zip"
    else
        archive_name="$pkg_name.tar.gz"
    fi

    pkg_dir="$staging/$pkg_name"
    rm -rf "$pkg_dir"
    mkdir -p "$pkg_dir/web"

    GOOS="$goos" GOARCH="$goarch" go build -trimpath -ldflags '-s -w' -o "$pkg_dir/$bin_name" .
    cp LICENSE config.example.yaml README.md README.en.md "$pkg_dir/"
    cp -r web/. "$pkg_dir/web/"
    # 一行版本号，供启动脚本打印、也便于用户确认手上这份是哪个版本
    printf '%s\n' "$version" > "$pkg_dir/version.txt"

    # 启动脚本：脚本自身会切到所在目录，因此可在任意位置调用
    if [ "$goos" = windows ]; then
        cp tools/release/start.bat "$pkg_dir/start.bat"
    else
        cp tools/release/start.sh "$pkg_dir/start.sh"
    fi

    if [ "$goos" = windows ]; then
        # zip 不保存可执行位，Windows 也不需要
        archive_win "$pkg_name" "$out_dir/$archive_name"
    else
        chmod 0755 "$pkg_dir/$bin_name" "$pkg_dir/start.sh"
        # 固定 mtime，保证同一份源码重复构建产物一致
        find "$pkg_dir" -exec touch -t 202601010000.00 {} +
        (cd "$staging" && tar -czf "$out_dir/$archive_name" "$pkg_name")
    fi

    size=$(du -h "$out_dir/$archive_name" | cut -f1)
    printf '  build %-16s -> %s  (%s)\n' "$target" "$archive_name" "$size"
    count=$((count + 1))
done

# 归档的 SHA256；只统计要发布的归档，LF 换行以便 sha256sum -c
hash_of() {
    if command -v sha256sum >/dev/null 2>&1; then
        sha256sum "$1" | cut -d' ' -f1
    else
        shasum -a 256 "$1" | cut -d' ' -f1   # macOS
    fi
}
find "$out_dir" -maxdepth 1 -type f \( -name '*.zip' -o -name '*.tar.gz' \) | LC_ALL=C sort | while IFS= read -r file; do
    printf '%s  %s\n' "$(hash_of "$file")" "$(basename "$file")"
done > "$out_dir/checksums.txt"

echo
echo "$count 个归档 + checksums.txt -> $out_dir"
