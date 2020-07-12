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

func TestNewKafkaSourceRunEFactory(t *testing.T) {
	runEFactory := NewFakeKafkaSourceRunEFactory("fake_namespace")
	assert.Assert(t, runEFactory != nil)
}

func TestRunEFactory_KafkaSourceParams(t *testing.T) {
	runEFactory := NewFakeKafkaSourceRunEFactory("fake_namespace")
	assert.Assert(t, runEFactory.KafkaSourceFactory().KafkaSourceParams() != nil)
}

func TestRunEFactory_KafkaSourceFactory(t *testing.T) {
	runEFactory := NewFakeKafkaSourceRunEFactory("fake_namespace")
	assert.Assert(t, runEFactory.KafkaSourceFactory() != nil)
}

func TestRunEFactory_KafkaSourceClient(t *testing.T) {
	runEFactory := NewFakeKafkaSourceRunEFactory("fake_namespace")
	knSourceClient, err := runEFactory.KafkaSourceClient(&rest.Config{}, "fake_namespace")
	assert.Assert(t, knSourceClient != nil)
	assert.Assert(t, err == nil)
}

func TestCreateRunE(t *testing.T) {
	runEFactory := NewFakeKafkaSourceRunEFactory("fake_namespace")
	function := runEFactory.CreateRunE()
	assert.Assert(t, function != nil)
}

func TestDeleteRunE(t *testing.T) {
	runEFactory := NewFakeKafkaSourceRunEFactory("fake_namespace")
	function := runEFactory.DeleteRunE()
	assert.Assert(t, function != nil)
}

func TestUpdateRunE(t *testing.T) {
	runEFactory := NewFakeKafkaSourceRunEFactory("fake_namespace")
	function := runEFactory.UpdateRunE()
	assert.Assert(t, function != nil)
}

func TestDescribeRunE(t *testing.T) {
	runEFactory := NewFakeKafkaSourceRunEFactory("fake_namespace")
	function := runEFactory.DescribeRunE()
	assert.Assert(t, function != nil)
}
