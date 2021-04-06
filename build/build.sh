#!/bin/bash

set -eu

if [[ $# -ne 1 ]]; then
  echo "Usage: $0 <pkg_version>"
  exit
fi

pkg_version="$1"

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")"&& pwd)"

BUILD_DIR="${SCRIPT_DIR}"
BUILD_TMP_DIR="${BUILD_DIR}/tmp"
pkg_name="alertmanager-webhook-adapter"


function build() {
  local _os="$1"
  local _arch="$2"

  pkg_dir="${pkg_name}-${pkg_version}-${_os}-${_arch}"

  mkdir -p "${BUILD_TMP_DIR}/${pkg_dir}"
  rm -rf ${BUILD_TMP_DIR}/${pkg_dir}.tar.gz

  cd "${SCRIPT_DIR}/../cmd/alertmanager-webhook-adapter"

  env GOOS="${_os}" GOARCH="${_arch}" go build -v -o "${BUILD_TMP_DIR}/${pkg_dir}/${pkg_name}"

  cd "${BUILD_TMP_DIR}"
  tar -czvf ${pkg_dir}.tar.gz ${pkg_dir}
}

build "darwin" "arm64"
build "darwin" "amd64"
build "linux" "arm64"
build "linux" "amd64"
