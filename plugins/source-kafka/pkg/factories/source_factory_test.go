// Copyright © 2020 The Knative Authors
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

package factories

import (
	"testing"

	"gotest.tools/assert"
	"k8s.io/client-go/rest"
)

func TestNewKafkaSourceFactory(t *testing.T) {
	factory := NewFakeKafkaSourceFactory("fake-namespace")

	assert.Assert(t, factory != nil)
}

func TestCreateKafkaSourceParams(t *testing.T) {
	factory := NewFakeKafkaSourceFactory("fake-namespace")

	sourceParams := factory.CreateKafkaSourceParams()
	assert.Assert(t, sourceParams != nil)
	assert.Equal(t, factory.KafkaSourceParams(), sourceParams)
}

func TestCreateKafkaSourceClient(t *testing.T) {
	factory := NewFakeKafkaSourceFactory("fake-namespace")
	client, _ := factory.CreateKafkaSourceClient(&rest.Config{}, "fake-namespace")

	assert.Assert(t, client != nil)
	assert.Equal(t, factory.KafkaSourceClient(), client)
	assert.Equal(t, factory.CreateKnSourceClient(&rest.Config{}, "fake-namespace"), client)
	assert.Equal(t, client.Namespace(), "fake-namespace")
}
