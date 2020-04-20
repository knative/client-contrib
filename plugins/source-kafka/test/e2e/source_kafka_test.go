// Copyright 2020 The Knative Authors

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// +build e2e
// +build !eventing

package e2e

import (
	"os"
	"path/filepath"
	"testing"

	testcommon "github.com/maximilien/kn-source-pkg/test/e2e"
	"gotest.tools/assert"
	"knative.dev/client/lib/test"
	"knative.dev/client/pkg/util"
)

const (
	kafkaBootstrapUrl     = "my-cluster-kafka-bootstrap.kafka.svc:9092"
	kafkaClusterName      = "my-cluster"
	kafkaClusterNamespace = "kafka"
	kafkaTopic            = "test-topic"
)

func TestSourceKafka(t *testing.T) {
	t.Parallel()

	currentDir, err := os.Getwd()
	assert.NilError(t, err)

	it, err := testcommon.NewE2ETest("kn-source_kafka", filepath.Join(currentDir, "../.."), false)
	assert.NilError(t, err)
	defer func() {
		assert.NilError(t, it.KnTest().Teardown())
	}()

	r := test.NewKnRunResultCollector(t, it.KnTest())
	defer r.DumpIfFailed()

	err = it.KnPlugin().Install()
	assert.NilError(t, err)

	serviceCreate(r, "sinksvc")

	t.Log("test kn-source_kafka create source-name")
	knSourceKafkaCreate(it, r, "mykafka1", "sinksvc")

	err = it.KnPlugin().Uninstall()
	assert.NilError(t, err)
}

// Private

func knSourceKafkaCreate(it *testcommon.E2ETest, r *test.KnRunResultCollector, sourceName, sinkName string) {
	out := it.KnPlugin().Run("create", sourceName, "--servers", kafkaBootstrapUrl, "--topics", kafkaTopic, "--consumergroup", "test-consumer-group", "--sink", sinkName)
	r.AssertNoError(out)
	assert.Check(r.T(), util.ContainsAllIgnoreCase(out.Stdout, "create", sourceName))
}

func serviceCreate(r *test.KnRunResultCollector, serviceName string) {
	out := r.KnTest().Kn().Run("service", "create", serviceName, "--image", test.KnDefaultTestImage)
	r.AssertNoError(out)
	assert.Check(r.T(), util.ContainsAllIgnoreCase(out.Stdout, "service", serviceName, "creating", "namespace", r.KnTest().Kn().Namespace(), "ready"))
}
