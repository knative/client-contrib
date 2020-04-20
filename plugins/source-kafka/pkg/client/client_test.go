// Copyright © 2019 The Knative Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package client

import (
	"testing"

	"gotest.tools/assert"
	client_testing "k8s.io/client-go/testing"
	"knative.dev/client-contrib/plugins/source-kafka/pkg/types"
	v1alpha1 "knative.dev/eventing-contrib/kafka/source/pkg/apis/sources/v1alpha1"
	fake "knative.dev/eventing-contrib/kafka/source/pkg/client/clientset/versioned/typed/sources/v1alpha1/fake"
)

func TestKafkaSourceClient(t *testing.T) {
	fakeE := fake.FakeSourcesV1alpha1{Fake: &client_testing.Fake{}}
	knSourceClient := NewKafkaSourceClient(&types.KafkaSourceParams{}, &fakeE, "fake-namespace")
	assert.Assert(t, knSourceClient != nil)
}

func TestClient_KnSourceParams(t *testing.T) {
	fakeE := fake.FakeSourcesV1alpha1{Fake: &client_testing.Fake{}}
	fakeKafkaParams := &types.KafkaSourceParams{}
	knSourceClient := NewKafkaSourceClient(fakeKafkaParams, &fakeE, "fake-namespace")
	assert.Equal(t, knSourceClient.KnSourceParams(), fakeKafkaParams.KnSourceParams)
}

func TestNamespace(t *testing.T) {
	fakeE := fake.FakeSourcesV1alpha1{Fake: &client_testing.Fake{}}
	knSourceClient := NewKafkaSourceClient(&types.KafkaSourceParams{}, &fakeE, "fake-namespace")
	assert.Equal(t, knSourceClient.Namespace(), "fake-namespace")
}
func TestCreateKafka(t *testing.T) {
	fakeE := fake.FakeSourcesV1alpha1{Fake: &client_testing.Fake{}}
	cli := NewKafkaSourceClient(&types.KafkaSourceParams{}, &fakeE, "fake-namespace")
	objNew := newKafkaSource("samplekafka")
	err := cli.CreateKafkaSource(objNew)
	assert.NilError(t, err)
}

func newKafkaSource(name string) *v1alpha1.KafkaSource {
	return NewKafkaSourceBuilder(name).
		BootstrapServers("test.server.org").
		Topics("topic").
		ConsumerGroup("mygroup").
		Build()
}
