#!/bin/bash

# Update test infrastructure code in ../test-infra to the latest version and add a VERSION file
# which contains the version updated to. Please commit this directory after the update


run() {
  local infra_dir=$(mktemp -d /tmp/test-infra.XXXXXX)
  local base_dir="$(basedir)"
  pushd $infra_dir >&/dev/null

  git clone --depth 1 https://github.com/knative/test-infra.git
  cd test-infra
  git log --pretty=format:"%as - %H%n" > VERSION
  rm -rf .git
  cd ..

  if [ -d "$base_dir/test-infra" ]; then
    rm -rf "$base_dir/test-infra"
  fi

  mv test-infra "$base_dir/"

  popd >&/dev/null
  rmdir $infra_dir
}

basedir() {
  # Default is current directory
  local script=${BASH_SOURCE[0]}
  local dir=$(dirname "$script")
  local full_dir=$(cd "${dir}/.." && pwd)
  echo ${full_dir}
}

run
