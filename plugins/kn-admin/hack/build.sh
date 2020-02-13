#!/bin/bash
set -o pipefail
source_dirs="cmd"
# Dir where this script is located
basedir() {
  # Default is current directory
  local script=${BASH_SOURCE[0]}

  # Resolve symbolic links
  if [ -L "$script" ]; then
    if readlink -f "$script" >/dev/null 2>&1; then
      script=$(readlink -f "$script")
    elif readlink "$script" >/dev/null 2>&1; then
      script=$(readlink "$script")
    elif realpath "$script" >/dev/null 2>&1; then
      script=$(realpath "$script")
    else
      echo "ERROR: Cannot resolve symbolic link $script"
      exit 1
    fi
  fi

  local dir=$(dirname "$script")
  local full_dir=$(cd "${dir}/.." && pwd)
  echo "${full_dir}"
}

go_fmt() {
  echo "🧹 ${S}Format"
  find $(echo "${source_dirs}") -name "*.go" -print0 | xargs -0 gofmt -s -w
}

go_build() {
  echo "🚧 Compile"
  go build -mod=vendor -o build/_output/kn-admin $(basedir)/main.go
}

export GO111MODULE=on
export GOPROXY=direct
go_fmt
go_build
