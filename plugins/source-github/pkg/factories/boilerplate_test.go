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

// KnSourceFactory

func TestGHSourceFactory_KnSourceParams(t *testing.T) {
	ghSourceFactory := NewGHSourceFactory()

	knSourceParams := ghSourceFactory.KnSourceParams()
	assert.Assert(t, knSourceParams != nil)

	knSourceParams = ghSourceFactory.CreateKnSourceParams()
	assert.Equal(t, ghSourceFactory.KnSourceParams(), knSourceParams)
}

func TestGHSourceFactory_GHSourceParams(t *testing.T) {
	ghSourceFactory := NewGHSourceFactory()

	ghSourceParams := ghSourceFactory.GHSourceParams()
	assert.Assert(t, ghSourceParams != nil)

	ghSourceParams = ghSourceFactory.CreateGHSourceParams()
	assert.Equal(t, ghSourceFactory.GHSourceParams(), ghSourceParams)
}

// CommandFactory

func TestCommandFactory_GHSourceFactory(t *testing.T) {
	ghSourceFactory := NewGHSourceFactory()
	commandFactory := NewGHCommandFactory(ghSourceFactory)

	assert.Equal(t, commandFactory.GHSourceFactory(), ghSourceFactory)
}

// FlagsFactory

func TestFlagsFactory_KnSourceFactory(t *testing.T) {
	ghSourceFactory := NewGHSourceFactory()
	flagsFactory := NewGHFlagsFactory(ghSourceFactory)

	assert.Equal(t, flagsFactory.KnSourceFactory(), ghSourceFactory)
}

func TestFlagsFactory_GHSourceFactory(t *testing.T) {
	ghSourceFactory := NewGHSourceFactory()
	flagsFactory := NewGHFlagsFactory(ghSourceFactory)

	assert.Equal(t, flagsFactory.GHSourceFactory(), ghSourceFactory)
}

// RunEFactory

func TestRunEFactory_GHSourceClient(t *testing.T) {
	runEFactory := NewGHRunEFactory(createFakeGHSourceFactory())

	ghSourceClient := runEFactory.GHSourceClient("fake_namespace")
	assert.Assert(t, ghSourceClient != nil)
}

func TestRunEFactory_KnSourceParams(t *testing.T) {
	runEFactory := NewGHRunEFactory(createFakeGHSourceFactory())

	assert.Assert(t, runEFactory.GHSourceFactory().KnSourceParams() != nil)
}

func TestRunEFactory_GHSourceParams(t *testing.T) {
	runEFactory := NewGHRunEFactory(createFakeGHSourceFactory())

	assert.Assert(t, runEFactory.GHSourceFactory().GHSourceParams() != nil)
}

func TestRunEFactory_GHSourceClientFactory(t *testing.T) {
	runEFactory := NewGHRunEFactory(createFakeGHSourceFactory())

	assert.Assert(t, runEFactory.GHSourceFactory() != nil)
}
