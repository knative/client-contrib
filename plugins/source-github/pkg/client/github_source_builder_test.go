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
	"fmt"
	"testing"

	"gotest.tools/assert"

	duckv1 "knative.dev/pkg/apis/duck/v1"
)

var builder *GitHubSourceBuilder

func TestNewGitHubSourceBuilder(t *testing.T) {
	setup(t)
	assert.Assert(t, builder.ghSource.ObjectMeta.Name == "fake-builder")
}

func TestOrgRepo(t *testing.T) {
	setup(t)
	builder = builder.OrgRepo("fake-org", "fake-repo")
	assert.Assert(t, builder.ghSource.Spec.OwnerAndRepository == fmt.Sprintf("fake-org/fake-repo"))
}

func TestAPIURL(t *testing.T) {
	setup(t)
	builder = builder.APIURL("https://fake-api-url")
	assert.Assert(t, builder.ghSource.Spec.GitHubAPIURL == fmt.Sprintf("https://fake-api-url"))
}

func TestSink(t *testing.T) {
	setup(t)
	fakeSink := &duckv1.Destination{}
	builder := builder.Sink(fakeSink)
	assert.Assert(t, builder.ghSource.Spec.Sink == fakeSink)
}

func TestBuild(t *testing.T) {
	setup(t)
	build := builder.Build()
	assert.Assert(t, builder.ghSource == build)
}

// Private

func setup(t *testing.T) {
	builder = NewGitHubSourceBuilder("fake-builder")
	assert.Assert(t, builder != nil)
}
