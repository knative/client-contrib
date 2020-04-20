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
	"knative.dev/client-contrib/plugins/source-kafka/pkg/types"
)

func TestNewKafkaSourceRunEFactory(t *testing.T) {
	runEFactory := createKafkaSourceRunEFactory()

	assert.Assert(t, runEFactory != nil)
}

func TestRunEFactory_KafkaSourceParams(t *testing.T) {
	runEFactory := createKafkaSourceRunEFactory()

	assert.Assert(t, runEFactory.KafkaSourceFactory().KnSourceParams() != nil)
	assert.Assert(t, runEFactory.KafkaSourceFactory().KafkaSourceParams() != nil)
}

func TestRunEFactory_KafkaSourceFactory(t *testing.T) {
	runEFactory := createKafkaSourceRunEFactory()

	assert.Assert(t, runEFactory.KnSourceFactory() != nil)
	assert.Assert(t, runEFactory.KafkaSourceFactory() != nil)
}

func TestRunEFactory_KafkaSourceClient(t *testing.T) {
	runEFactory := createKafkaSourceRunEFactory()
	knSourceClient := runEFactory.KafkaSourceClient("fake_namespace")
	assert.Assert(t, knSourceClient != nil)
}

func TestCreateRunE(t *testing.T) {
	runEFactory := createKafkaSourceRunEFactory()
	function := runEFactory.CreateRunE()
	assert.Assert(t, function != nil)
}

func TestDeleteRunE(t *testing.T) {
	runEFactory := createKafkaSourceRunEFactory()
	function := runEFactory.DeleteRunE()
	assert.Assert(t, function != nil)
}

// Private

func createKafkaSourceRunEFactory() types.KafkaSourceRunEFactory {
	factory := NewKafkaSourceFactory()
	return NewKafkaSourceRunEFactory(factory)
}
