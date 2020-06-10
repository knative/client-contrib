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

	"github.com/maximilien/kn-source-github/pkg/client"
	"github.com/maximilien/kn-source-github/pkg/types"

	"knative.dev/eventing-contrib/github/pkg/apis/sources/v1alpha1"

	"gotest.tools/assert"

	sourcetypes "github.com/maximilien/kn-source-pkg/pkg/types"
)

func TestNewGHSourceFactory(t *testing.T) {
	ghSourceFactory := createFakeGHSourceFactory()

	assert.Assert(t, ghSourceFactory != nil)
}

func TestCreateKnSourceParams(t *testing.T) {
	ghSourceFactory := createFakeGHSourceFactory()

	ghSourceParams := ghSourceFactory.CreateKnSourceParams()
	assert.Assert(t, ghSourceParams != nil)
}

func TestCreateGHSourceParams(t *testing.T) {
	ghSourceFactory := createFakeGHSourceFactory()

	ghSourceParams := ghSourceFactory.CreateGHSourceParams()
	assert.Assert(t, ghSourceParams != nil)
}

func TestCreateKnSourceClient(t *testing.T) {
	ghSourceFactory := createFakeGHSourceFactory()
	client := ghSourceFactory.CreateKnSourceClient("fake-namespace")

	assert.Assert(t, client != nil)
	assert.Equal(t, client.Namespace(), "fake-namespace")
}

func TestCreateGHSourceClient(t *testing.T) {
	ghSourceFactory := createFakeGHSourceFactory()
	client := ghSourceFactory.CreateGHSourceClient("fake-namespace")

	assert.Assert(t, client != nil)
	assert.Equal(t, client.Namespace(), "fake-namespace")
}

// Private

func createFakeGHSourceFactory() types.GHSourceFactory {
	fakeGHSourceClient := client.NewFakeGHSourceClient(client.NewFakeGHSourceParams(), &v1alpha1.GitHubSource{})

	return &ghSourceFactory{
		ghSourceParams: &types.GHSourceParams{
			KnSourceParams: &sourcetypes.KnSourceParams{},
		},
		ghSourceClient: fakeGHSourceClient,
	}
}
