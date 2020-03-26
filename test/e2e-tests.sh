# Dispatch tests to the different plugin directories


run() {
  local basedir=$(basedir)

  # Environment variable which can be used my plugins
  export TEST_INFRA_SCRIPTS="$basedir/test-infra/scripts"

  source $TEST_INFRA_SCRIPTS/e2e-tests.sh

  # Iterate over all plugin directories and check whether they have testing
  # enabled
  for plugin in $(ls $basedir/plugins); do
    local test_script=$basedir/plugins/"${plugin}"/test/e2e-tests.sh
    if [ -x $basedir/plugins/$plugin/test/e2e-tests.sh ]; then
      echo "== $plugin ==============================="
      eval "$test_script" || fail_test
      echo "=========================================="
    fi
  done
}

basedir() {
  # Default is current directory
  local script=${BASH_SOURCE[0]}
  local dir=$(dirname "$script")
  local full_dir=$(cd "${dir}/.." && pwd)
  echo ${full_dir}
}

run $@
