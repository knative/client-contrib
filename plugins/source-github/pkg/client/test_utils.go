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
	"k8s.io/client-go/rest"

	"knative.dev/eventing-contrib/github/pkg/apis/sources/v1alpha1"

	"github.com/maximilien/kn-source-github/pkg/types"

	sourcestypes "github.com/maximilien/kn-source-pkg/pkg/types"

	fakes "github.com/maximilien/kn-source-github/pkg/fakes"
	sourcestypesfakes "github.com/maximilien/kn-source-pkg/pkg/types/typesfakes"
)

func NewFakeGHSourceParams() *types.GHSourceParams {
	return &types.GHSourceParams{
		KnSourceParams: &sourcestypes.KnSourceParams{},

		Org:  "fake-org",
		Repo: "fake-repo",

		APIURL:      "https://fake-api-url",
		SecretToken: "fake-secret-token",
		AccessToken: "fake-access-token",
	}
}

func NewFakeGHSourceClient(ghSourceParams *types.GHSourceParams, ghSource *v1alpha1.GitHubSource) types.GHSourceClient {
	fakeGitHubSources := &fakes.FakeGitHubSourceInterface{}
	fakeGitHubSources.GetReturns(ghSource, nil)
	fakeGitHubSources.CreateReturns(ghSource, nil)
	fakeGitHubSources.UpdateReturns(ghSource, nil)
	fakeGitHubSources.DeleteReturns(nil)

	fakeGHSourcesV1 := &fakes.FakeSourcesV1alpha1Interface{}
	fakeGHSourcesV1.GitHubSourcesReturns(fakeGitHubSources)

	knSourceClient := &sourcestypesfakes.FakeKnSourceClient{}
	knSourceClient.KnSourceParamsReturns(ghSourceParams.KnSourceParams)
	knSourceClient.NamespaceReturns("fake-namespace")
	knSourceClient.RestConfigReturns(&rest.Config{})

	return &ghSourceClient{
		namespace:      "fake-namespace",
		ghSourceParams: ghSourceParams,
		knSourceClient: knSourceClient,
		ghSourcesV1:    fakeGHSourcesV1,
	}
}
