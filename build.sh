#!/usr/bin/env bash
#
# 交叉编译 mini-im 并打包为可分发的归档，产物落在 dist/。
#
# 每个归档是一个顶层目录 mini-im-<os>-<arch>/，解压后直接可用：
#
#   mini-im-<os>-<arch>/
#   ├── mini-im(.exe)   可执行文件（非 Windows 平台带 0755 权限）
#   ├── config.yaml     默认配置（数据库 / Redis / 端口，可用 IM_* 环境变量覆盖）
#   ├── LICENSE
#   └── web/            前端静态页（后端直接托管，必须与可执行文件放在一起）
#
# 另有 checksums.txt 记录全部归档的 SHA256（LF 换行，Linux 下可直接 sha256sum -c）。
#
# 用法: ./build.sh [-o OutDir]

set -euo pipefail

usage() {
    cat <<'EOF'
交叉编译 mini-im 并打包为可分发的归档。

用法: ./build.sh [-o OutDir]

选项:
  -o OutDir   输出目录（默认: dist，相对本脚本所在目录）
  -h          显示本帮助

平台矩阵: windows/amd64 windows/arm64 linux/amd64 linux/arm64 darwin/amd64 darwin/arm64
每个平台一个归档（Windows 为 .zip，其余为 .tar.gz），外加 checksums.txt。
EOF
}

out_dir=dist
while getopts ':ho:' opt; do
    case "$opt" in
        h) usage; exit 0 ;;
        o) out_dir=$OPTARG ;;
        '?') usage >&2; exit 2 ;;
    esac
done

cd "$(dirname "$0")"

if [ -z "$out_dir" ]; then
    echo "OutDir 不能为空" >&2
    exit 2
fi
case "$out_dir" in
    /*) ;;
    *) out_dir=$PWD/$out_dir ;;
esac

# 目标矩阵，'os/arch' 形式；增删一行即可调整发布平台。
targets='windows/amd64 windows/arm64 linux/amd64 linux/arm64 darwin/amd64 darwin/arm64'

# 交叉编译不需要 cgo（纯 Go 依赖，含 PostgreSQL / Redis 驱动）
export CGO_ENABLED=0

mkdir -p "$out_dir"
# 只删除输出目录顶层的文件，不递归，避免误伤目录外的东西
find "$out_dir" -maxdepth 1 -type f -delete

staging=$(mktemp -d)
trap 'rm -rf "$staging"' EXIT

count=0
for target in $targets; do
    goos=${target%%/*}
    goarch=${target#*/}

    # 归档内保留短名与顶层目录名，都不带版本号
    bin_name=mini-im
    if [ "$goos" = windows ]; then
        bin_name=mini-im.exe
    fi
    pkg_name="mini-im-$goos-$goarch"
    if [ "$goos" = windows ]; then
        archive_name="$pkg_name.zip"
    else
        archive_name="$pkg_name.tar.gz"
    fi

    pkg_dir="$staging/$pkg_name"
    rm -rf "$pkg_dir"
    mkdir -p "$pkg_dir/web"

    GOOS="$goos" GOARCH="$goarch" go build -trimpath -ldflags '-s -w' -o "$pkg_dir/$bin_name" .
    cp config.yaml LICENSE "$pkg_dir/"
    cp -r web/. "$pkg_dir/web/"

    if [ "$goos" = windows ]; then
        # zip 不保存可执行位，Windows 也不需要
        (cd "$staging" && zip -qr "$out_dir/$archive_name" "$pkg_name")
    else
        chmod 0755 "$pkg_dir/$bin_name"
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
