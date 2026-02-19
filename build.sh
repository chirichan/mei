#!/usr/bin/env bash

set -euo pipefail

release_dir="release"

# 获取最新的 Git 标签版本，无 tag 时回退到 commit short hash。
VERSION="$(git describe --tags --abbrev=0 2>/dev/null || git rev-parse --short=8 HEAD)"
BUILD_DATE="$(date -u +'%Y-%m-%dT%H:%M:%SZ')"
GIT_COMMIT_SHORT="$(git rev-parse --short=8 HEAD)"

LDFLAGS="-X github.com/chirichan/mei/version.Version=${VERSION} -X github.com/chirichan/mei/version.BuildTime=${BUILD_DATE} -X github.com/chirichan/mei/version.GitCommit=${GIT_COMMIT_SHORT}"

mkdir -p "${release_dir}"
rm -f "${release_dir}"/pwdgen-* "${release_dir}"/bing15-*

build_bin() {
	local goos="$1"
	local goarch="$2"
	local output="$3"
	local pkg="$4"

	echo "building ${output}"
	CGO_ENABLED=0 GOOS="${goos}" GOARCH="${goarch}" go build -ldflags "${LDFLAGS}" -o "${output}" "${pkg}"
}

build_bin darwin arm64 "${release_dir}/pwdgen-darwin-arm64" "./cmd/pwdgen/"
build_bin darwin amd64 "${release_dir}/pwdgen-darwin-x86_64" "./cmd/pwdgen/"
build_bin linux amd64 "${release_dir}/pwdgen-linux-x86_64" "./cmd/pwdgen/"
build_bin windows amd64 "${release_dir}/pwdgen-windows-x86_64.exe" "./cmd/pwdgen/"

build_bin darwin arm64 "${release_dir}/bing15-darwin-arm64" "./cmd/bing15/"
build_bin darwin amd64 "${release_dir}/bing15-darwin-x86_64" "./cmd/bing15/"
build_bin linux amd64 "${release_dir}/bing15-linux-x86_64" "./cmd/bing15/"
build_bin windows amd64 "${release_dir}/bing15-windows-x86_64.exe" "./cmd/bing15/"

echo "build success: version=${VERSION}, git_commit=${GIT_COMMIT_SHORT}, build_time=${BUILD_DATE}"
