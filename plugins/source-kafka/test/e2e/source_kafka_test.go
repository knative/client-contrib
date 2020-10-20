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

type e2eTest struct {
	it *testcommon.E2ETest
}

func newE2ETest(t *testing.T) *e2eTest {
	currentDir, err := os.Getwd()
	if err != nil {
		return nil
	}

	it, err := testcommon.NewE2ETest("kn-source-kafka", filepath.Join(currentDir, "../.."), false)
	if err != nil {
		return nil
	}

	e2eTest := &e2eTest{
		it: it,
	}
	return e2eTest
}

func TestSourceKafka(t *testing.T) {
	t.Parallel()

	e2eTest := newE2ETest(t)
	assert.Assert(t, e2eTest != nil)
	defer func() {
		assert.NilError(t, e2eTest.it.KnTest().Teardown())
	}()

	r := test.NewKnRunResultCollector(t, e2eTest.it.KnTest())
	defer r.DumpIfFailed()

	err := e2eTest.it.KnPlugin().Install()
	assert.NilError(t, err)

	serviceCreate(r, "sinksvc")

	t.Log("test kn-source-kafka create source-name")
	e2eTest.knSourceKafkaCreate(t, r, "mykafka1", "sinksvc")

	t.Log("test kn-source-kafka describe source-name")
	e2eTest.knSourceKafkaDescribe(t, r, "mykafka1", "sinksvc")

	t.Log("test kn-source-kafka delete source-name")
	e2eTest.knSourceKafkaDelete(t, r, "mykafka1")

	err = e2eTest.it.KnPlugin().Uninstall()
	assert.NilError(t, err)
}

// Private

func (et *e2eTest) knSourceKafkaCreate(t *testing.T, r *test.KnRunResultCollector, sourceName, sinkName string) {
	out := et.it.KnPlugin().Run("create", sourceName, "--servers", kafkaBootstrapUrl, "--topics", kafkaTopic, "--consumergroup", "test-consumer-group", "--sink", sinkName)
	r.AssertNoError(out)
	assert.Check(t, util.ContainsAllIgnoreCase(out.Stdout, "create", sourceName))
}

func (et *e2eTest) knSourceKafkaDelete(t *testing.T, r *test.KnRunResultCollector, sourceName string) {
	out := et.it.KnPlugin().Run("delete", sourceName)
	r.AssertNoError(out)
	assert.Check(t, util.ContainsAllIgnoreCase(out.Stdout, "delete", sourceName))
}

func (et *e2eTest) knSourceKafkaDescribe(t *testing.T, r *test.KnRunResultCollector, sourceName, sinkName string) {
	out := et.it.KnPlugin().Run("describe", sourceName)
	r.AssertNoError(out)
	assert.Check(t, util.ContainsAllIgnoreCase(out.Stdout, sourceName, sinkName))
}

func serviceCreate(r *test.KnRunResultCollector, serviceName string) {
	out := r.KnTest().Kn().Run("service", "create", serviceName, "--image", "gcr.io/knative-samples/helloworld-go")
	r.AssertNoError(out)
	assert.Check(r.T(), util.ContainsAllIgnoreCase(out.Stdout, "service", serviceName, "creating", "namespace", r.KnTest().Kn().Namespace(), "ready"))
}
