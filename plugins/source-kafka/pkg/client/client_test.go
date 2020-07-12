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
	v1alpha1 "knative.dev/eventing-contrib/kafka/source/pkg/apis/sources/v1alpha1"
)

func TestKafkaSourceClient(t *testing.T) {
	knSourceClient := NewFakeKafkaSourceClient("fake-namespace")
	assert.Assert(t, knSourceClient != nil)
}

func TestClient_KnSourceParams(t *testing.T) {
	knSourceClient := NewFakeKafkaSourceClient("fake-namespace")
	fakeKafkaParams := knSourceClient.KafkaSourceParams()
	assert.Equal(t, knSourceClient.KnSourceParams(), fakeKafkaParams.KnSourceParams)
}

func TestNamespace(t *testing.T) {
	knSourceClient := NewFakeKafkaSourceClient("fake-namespace")
	assert.Equal(t, knSourceClient.Namespace(), "fake-namespace")
}
func TestCreateKafka(t *testing.T) {
	cli := NewFakeKafkaSourceClient("fake-namespace")
	objNew := newKafkaSource("samplekafka")
	err := cli.CreateKafkaSource(objNew)
	assert.NilError(t, err)
}

func TestDeleteKafka(t *testing.T) {
	cli := NewFakeKafkaSourceClient("fake-namespace")
	objNew := newKafkaSource("samplekafka")
	err := cli.CreateKafkaSource(objNew)
	assert.NilError(t, err)
	err = cli.DeleteKafkaSource("samplekafka")
	assert.NilError(t, err)
}

func newKafkaSource(name string) *v1alpha1.KafkaSource {
	return NewKafkaSourceBuilder(name).
		BootstrapServers("test.server.org").
		Topics("topic").
		ConsumerGroup("mygroup").
		Build()
}
