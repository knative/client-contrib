#!/usr/bin/env bash

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

source "$(dirname $0)"/../vendor/knative.dev/test-infra/scripts/e2e-tests.sh

export PATH=$PWD:$PATH

dir=$(dirname "${BASH_SOURCE[0]}")
base=$(cd "$dir/.." && pwd)

# Strimzi installation config template used for starting up Kafka clusters.
readonly STRIMZI_INSTALLATION_CONFIG_TEMPLATE="test/config/100-strimzi-cluster-operator-0.17.0.yaml"
# Strimzi installation config.
readonly STRIMZI_INSTALLATION_CONFIG="$(mktemp)"
# Kafka cluster CR config file.
readonly KAFKA_INSTALLATION_CONFIG="test/config/100-kafka-ephemeral-triple-2.4.0.yaml"
readonly KAFKA_TOPIC_INSTALLATION_CONFIG="test/config/100-kafka-topic.yaml"
# Kafka cluster URL for our installation
readonly KAFKA_CLUSTER_URL="my-cluster-kafka-bootstrap.kafka:9092"
# Kafka channel CRD config template directory.
readonly KAFKA_CRD_CONFIG_TEMPLATE_DIR="kafka/channel/config"
# Kafka channel CRD config template file. It needs to be modified to be the real config file.
readonly KAFKA_CRD_CONFIG_TEMPLATE="400-kafka-config.yaml"
# Real Kafka channel CRD config , generated from the template directory and modified template file.
readonly KAFKA_CRD_CONFIG_DIR="$(mktemp -d)"
# Kafka channel CRD config template directory.
readonly KAFKA_SOURCE_CRD_YAML="https://github.com/knative/eventing-contrib/releases/download/v0.14.0/kafka-source.yaml"

function kafka_setup() {
  echo "Installing Kafka cluster"
  kubectl create namespace kafka || return 1
  sed 's/namespace: .*/namespace: kafka/' ${STRIMZI_INSTALLATION_CONFIG_TEMPLATE} > ${STRIMZI_INSTALLATION_CONFIG}
  kubectl apply -f "${STRIMZI_INSTALLATION_CONFIG}" -n kafka
  kubectl apply -f ${KAFKA_INSTALLATION_CONFIG} -n kafka
  kubectl apply -f ${KAFKA_TOPIC_INSTALLATION_CONFIG} -n kafka
  wait_until_pods_running kafka || fail_test "Failed to start up a Kafka cluster"
}

function kafka_teardown() {
  echo "Uninstalling Kafka cluster"
  kubectl delete -f ${KAFKA_TOPIC_INSTALLATION_CONFIG} -n kafka
  kubectl delete -f ${KAFKA_INSTALLATION_CONFIG} -n kafka
  kubectl delete -f "${STRIMZI_INSTALLATION_CONFIG}" -n kafka
  kubectl delete namespace kafka
}

function plugin_test_setup() {
  kafka_setup || return 1
  install_sources_crds || return 1
}

function plugin_test_teardown() {
  kafka_teardown
  uninstall_sources_crds
}

function install_sources_crds() {
  echo "Installing Kafka Source CRD"
  kubectl apply -f ${KAFKA_SOURCE_CRD_YAML}

  # wait_until_pods_running knative-eventing || fail_test "Failed to install the Kafka Source CRD"
  wait_until_pods_running knative-sources || fail_test "Failed to install the Kafka Source CRD"
}

function uninstall_sources_crds() {
  echo "Uninstalling Kafka Source CRD"
  kubectl delete -f ${KAFKA_SOURCE_CRD_YAML}
}

# Will create and delete this namespace (used for all tests, modify if you want a different one used)
export KN_E2E_NAMESPACE=kne2etests

echo "🧪  Setup"
plugin_test_setup
echo "🧪  Testing"
#go test ${base}/test/e2e/ -timeout=45m -test.v -tags "e2e ${E2E_TAGS}" "$@"
echo "🧪  Teardown"
#plugin_test_teardown
