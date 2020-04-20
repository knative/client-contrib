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
)

func TestNewKafkaSourceFactory(t *testing.T) {
	factory := NewKafkaSourceFactory()

	assert.Assert(t, factory != nil)
}

func TestKafkaCreateKnSourceParams(t *testing.T) {
	factory := NewKafkaSourceFactory()

	knSourceParams := factory.KnSourceParams()
	assert.Assert(t, knSourceParams != nil)

	knSourceParams = factory.CreateKnSourceParams()
	assert.Equal(t, factory.KnSourceParams(), knSourceParams)
}

func TestCreateKafkaSourceParams(t *testing.T) {
	factory := NewKafkaSourceFactory()

	sourceParams := factory.CreateKafkaSourceParams()
	assert.Assert(t, sourceParams != nil)
	assert.Equal(t, factory.KafkaSourceParams(), sourceParams)
}

func TestCreateKafkaSourceClient(t *testing.T) {
	factory := NewKafkaSourceFactory()
	client := factory.CreateKafkaSourceClient("fake-namespace")

	assert.Assert(t, client != nil)
	assert.Equal(t, factory.KafkaSourceClient(), client)
	assert.Equal(t, factory.CreateKnSourceClient("fake-namespace"), client)
	assert.Equal(t, client.Namespace(), "fake-namespace")
}
