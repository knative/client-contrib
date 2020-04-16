# Dispatch tests to the different plugin directories

source "$(dirname $0)"/common.sh
source "$(dirname $0)"/../test-infra/scripts/e2e-tests.sh
source "$(dirname $0)"/../test-infra/scripts/presubmit-tests.sh

run() {
  local basedir=$(basedir)

  # Environment variable which can be used my plugins
  export TEST_INFRA_SCRIPTS="$basedir/test-infra/scripts"

  # Iterate over all plugin directories and check whether they have testing
  # enabled

  echo "=== Running Plugins E2E tests ======================================="
  loop_over_plugins "e2e-tests.sh" ""
}

run $@
