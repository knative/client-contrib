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

func TestKafkaCreateKnSourceParams(t *testing.T) {
	factory := NewFakeKafkaSourceFactory("fake-namespace")

	knSourceParams := factory.KnSourceParams()
	assert.Assert(t, knSourceParams != nil)

	knSourceParams = factory.CreateKnSourceParams()
	assert.Equal(t, factory.KnSourceParams(), knSourceParams)
}

func TestRunEFactory_KnSourceParams(t *testing.T) {
	runEFactory := NewFakeKafkaSourceRunEFactory("fake_namespace")
	assert.Assert(t, runEFactory.KafkaSourceFactory().KnSourceParams() != nil)
}

func TestRunEFactory_KnSourceFactory(t *testing.T) {
	runEFactory := NewFakeKafkaSourceRunEFactory("fake_namespace")
	assert.Assert(t, runEFactory.KnSourceFactory() != nil)
}

func TestFlagsFactory_KnSourceFactory(t *testing.T) {
	flagsFactory := createKafkaSourceFlagsFactory()
	assert.Assert(t, flagsFactory.KnSourceFactory() != nil)
}

func TestCommand_KnSourceFactory(t *testing.T) {
	commandFactory := createKafkaSourceCommandFactory()
	assert.Assert(t, commandFactory.KnSourceFactory() != nil)
}
