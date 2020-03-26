# Dispatch tests to the different plugin directories

run() {
  # Iterate over all plugin directories and check whether they have testing
  # enabled

  local $basedir=$(basedir)
  for plugin in $(ls $basedir/plugins); do
    local test_script=$basedir/plugins/$plugin/test/e2e-tests.sh
    if [ -x $basedir/plugins/$plugin/test/e2e-tests.sh ]; then
      echo "== $plugin ==============================="
      eval "source $test_script" || fail_test
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
