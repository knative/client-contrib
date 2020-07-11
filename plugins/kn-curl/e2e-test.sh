#!/bin/bash

# Copyright © 2020 The Knative Authors
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

# Variables and defaults
KN=`which kn`
KUBECTL=`which kubectl`
E2E_TEST_NS=$"kn-curl-e2e-tests"
E2E_TEST_SERVICE_NAME=$"kn-curl-e2e-test-service"

# Check that kn is in PATH
check_kn_in_path() {
  if [ -z "${KN_PATH}" ]; then
    if [ -x "${KN}" ]; then
      echo "✅ Found kn executable: $KN"
    else
      echo "🔥 Could not find kn executable, please add it to your PATH or set KN_PATH"
      exit -1
    fi
  else
    echo "✅ KN_PATH is set to: $KN_PATH"
    export PATH=$KN_PATH:$PATH
  fi
}

# Check that kubectl is in PATH
check_kubectl_in_path() {
  if [ -z "${KUBECTL_PATH}" ]; then
    if [ -x "${KUBECTL}" ]; then
      echo "✅ Found kubectl executable: $KUBECTL"
    else
      echo "🔥 Could not find kubectl executable, please add it to your PATH or set KUBECTL_PATH"
      exit -1
    fi
  else
    echo "✅ KUBECTL_PATH is set to: $KUBECTL_PATH"
    export PATH=$KUBECTL_PATH:$PATH
  fi
}

# Check kn and kubectl
check_kn_in_path
check_kubectl_in_path

# Create e2e-test namespace
kubectl create ns $E2E_TEST_NS

# Check last call
function check() {
  if [ $? != 0 ]; then
    echo "🔥 failed to $1"
    exit 1
  fi
}

# Start
echo "🧪 start kn curl e2e tests"

# Test 1: create public service
kn service create $E2E_TEST_SERVICE_NAME --image knativesamples/helloworld --namespace $E2E_TEST_NS | check "create service"
kn curl $E2E_TEST_SERVICE_NAME | check "curl public service"

# Test 2: update servcie to private (cluster local) service
kn service update $E2E_TEST_SERVICE_NAME --cluster-local --namespace $E2E_TEST_NS | check "update service to cluster_local"
kn curl $E2E_TEST_SERVICE_NAME | check("curl private service")

# Cleanup
echo "🧹 cleanup"

# Delete service
kn service delete $E2E_TEST_SERVICE_NAME --namespace $E2E_TEST_NS | check "delete service"

# Delete namespace
kubectl delete ns $E2E_TEST_NS

wait 2
echo "✅ Success"
exit 0
