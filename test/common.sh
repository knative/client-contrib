#!/bin/bash

# Loop over all plugins and call a test script ($1) with the given arguments ($2)
function loop_over_plugins() {
  local script=${1:-}
  local opts=${2:-}

  local basedir=$(basedir)

  # Environment variable which can be used my plugins
  export TEST_INFRA_SCRIPTS="$basedir/test-infra/scripts"

  for plugin in "${basedir}"/plugins/*; do
    local test_script="${basedir}/plugins/${plugin}/test/$script"
    if [ -x "$test_script" ]; then
      echo "## $plugin ###############################"
      eval "$test_script $opts" || fail_test
      echo "##########################################"
    fi
  done
}

# Calculate the base directory
function basedir() {
  # Default is current directory
  local script=${BASH_SOURCE[0]}
  local dir=$(dirname "$script")
  local full_dir=$(cd "${dir}/.." && pwd)
  echo ${full_dir}
}

