# Copyright 2020 The Knative Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# ===============================================
# Add you integration tests here

source $TEST_INFRA_SCRIPTS/e2e-tests.sh

export PATH=$PWD:$PATH

dir=$(dirname "${BASH_SOURCE[0]}")
base=$(cd "$dir/.." && pwd)

echo "TEST_INFRA_SCRIPTS: $TEST_INFRA_SCRIPTS"
echo "Testing kn-admin plugin"
cd ${REPO_ROOT_DIR}

function plugin_test_setup() {
  header "Setting up plugin kn-admin"
  # TODO: add setup steps
}

function run() {

  header "Running plugin kn-admin e2e tests for Knative Serving $KNATIVE_SERVING_VERSION and Eventing $KNATIVE_EVENTING_VERSION"

  # Will create and delete this namespace (used for all tests, modify if you want a different one used)
  export KN_E2E_NAMESPACE=kne2etests

  echo "🧪  Setup"
  plugin_test_setup
  echo "🧪  Build"
  ./hack/build.sh -f
  echo "🧪  Testing"
  go_test_e2e -timeout=45m ./test/e2e || fail_test
  echo "🧪  Teardown"
  plugin_test_teardown
  success
}

# Fire up
run $@
