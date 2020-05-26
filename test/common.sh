#!/bin/bash

# Loop over all plugins and call a test script ($1) with the given arguments ($2)
function loop_over_plugins() {
  local script=${1:-}
  local opts=${2:-}

  local basedir=$(basedir)

  # Environment variable which can be used my plugins
  export TEST_INFRA_SCRIPTS="$basedir/scripts/test-infra"


  plugins=$(list_plugins_changed_in_pr)
  echo "--- Plugins changed in PR: ----------------"
  echo "$plugins"
  echo "-------------------------------------------"
  for plugin in ${plugins}; do
    local plugin_dir="${basedir}/plugins/$plugin"
    local test_script="${plugin_dir}/test/$script"
    if [ -x "$test_script" ]; then
      echo "## $plugin ###############################"
      bash -c "REPO_ROOT_DIR=$plugin_dir $test_script $opts"
      local err=$?
      if [ $err -gt 0 ]; then
        fail_sub_test "Plugin $plugin failed with $err"
      fi
      echo "##########################################"
    fi
  done
}

function list_plugins_changed_in_pr() {
  list_changed_files | grep "^plugins/" | sed -e 's|plugins/\([^/]*\).*|\1|' | uniq | sort
}

function fail_sub_test() {
  [[ -n $1 ]] && echo "ERROR: $1"
  exit 1
}

# Calculate the base directory
function basedir() {
  # Default is current directory
  local script=${BASH_SOURCE[0]}
  local dir=$(dirname "$script")
  local full_dir=$(cd "${dir}/.." && pwd)
  echo ${full_dir}
}

function knative_setup() {
  local serving_version=${KNATIVE_SERVING_VERSION:-latest}
  header "Installing Knative Serving (${serving_version})"

  if [ "${serving_version}" = "latest" ]; then
    start_latest_knative_serving
  else
    start_release_knative_serving "${serving_version}"
  fi

  local eventing_version=${KNATIVE_EVENTING_VERSION:-latest}
  header "Installing Knative Eventing (${eventing_version})"

  if [ "${eventing_version}" = "latest" ]; then
    start_latest_knative_eventing
  else
    start_release_knative_eventing "${eventing_version}"
  fi
}

function cluster_setup() {
  header "Installing client"
  go get knative.dev/client
  go install knative.dev/client/cmd/kn
}
