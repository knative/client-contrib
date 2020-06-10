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

	"knative.dev/eventing-contrib/github/pkg/apis/sources/v1alpha1"
)

func TestKnSourceParams(t *testing.T) {
	ghSourceParams := NewFakeGHSourceParams()
	ghSourceClient := NewFakeGHSourceClient(ghSourceParams, &v1alpha1.GitHubSource{})

	assert.Equal(t, ghSourceClient.KnSourceParams(), ghSourceParams.KnSourceParams)
}

func TestGHSourceParams(t *testing.T) {
	ghSourceParams := NewFakeGHSourceParams()
	ghSourceClient := NewFakeGHSourceClient(ghSourceParams, &v1alpha1.GitHubSource{})

	assert.Equal(t, ghSourceClient.GHSourceParams(), ghSourceParams)
}

func TestNamespace(t *testing.T) {
	ghSourceClient := NewFakeGHSourceClient(NewFakeGHSourceParams(), &v1alpha1.GitHubSource{})
	assert.Equal(t, ghSourceClient.Namespace(), "fake-namespace")
}

func TestRestConfig(t *testing.T) {
	ghSourceClient := NewFakeGHSourceClient(NewFakeGHSourceParams(), &v1alpha1.GitHubSource{})
	assert.Assert(t, ghSourceClient.RestConfig() != nil)
}

func TestGetGHSource(t *testing.T) {
	fakeGHSource := &v1alpha1.GitHubSource{}
	ghSourceClient := NewFakeGHSourceClient(NewFakeGHSourceParams(), fakeGHSource)
	ghSource, err := ghSourceClient.GetGHSource("fake-name")
	assert.Assert(t, err == nil)
	assert.Assert(t, ghSource == fakeGHSource)
}

func TestCreateGHSource(t *testing.T) {
	fakeGHSource := &v1alpha1.GitHubSource{}
	ghSourceClient := NewFakeGHSourceClient(NewFakeGHSourceParams(), fakeGHSource)
	ghSource, err := ghSourceClient.CreateGHSource(fakeGHSource)
	assert.Assert(t, err == nil)
	assert.Assert(t, ghSource == fakeGHSource)
}

func TestUpdateGHSource(t *testing.T) {
	fakeGHSource := &v1alpha1.GitHubSource{}
	ghSourceClient := NewFakeGHSourceClient(NewFakeGHSourceParams(), fakeGHSource)
	ghSource, err := ghSourceClient.UpdateGHSource(fakeGHSource)
	assert.Assert(t, err == nil)
	assert.Assert(t, ghSource == fakeGHSource)
}

func TestDeleteGHSource(t *testing.T) {
	fakeGHSource := &v1alpha1.GitHubSource{}
	ghSourceClient := NewFakeGHSourceClient(NewFakeGHSourceParams(), fakeGHSource)
	err := ghSourceClient.DeleteGHSource("fake-name")
	assert.Assert(t, err == nil)
}
