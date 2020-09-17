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

	v1alpha1 "knative.dev/eventing-contrib/github/pkg/apis/sources/v1alpha1"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	duckv1 "knative.dev/pkg/apis/duck/v1"
)

// GitHubSourceBuilder is for building the source
type GitHubSourceBuilder struct {
	ghSource *v1alpha1.GitHubSource
}

// NewGitHubSourceBuilder for building ApiServer source object
func NewGitHubSourceBuilder(name string) *GitHubSourceBuilder {
	return &GitHubSourceBuilder{
		ghSource: &v1alpha1.GitHubSource{
			ObjectMeta: metav1.ObjectMeta{
				Name: name,
			},
		},
	}
}

// OrgRepo sets the org/owner and repository
func (b *GitHubSourceBuilder) OrgRepo(org, repo string) *GitHubSourceBuilder {
	b.ghSource.Spec.OwnerAndRepository = fmt.Sprintf("%s/%s", org, repo)
	return b
}

// APIURL to set the value of the GitHub APIURL
func (b *GitHubSourceBuilder) APIURL(apiURL string) *GitHubSourceBuilder {
	b.ghSource.Spec.GitHubAPIURL = apiURL
	return b
}

// AccessToken the access-token to use for this GitHub source
func (b *GitHubSourceBuilder) AccessToken(accessToken string) *GitHubSourceBuilder {
	b.ghSource.Spec.AccessToken = v1alpha1.SecretValueFromSource{
		SecretKeyRef: &corev1.SecretKeySelector{
			Key: accessToken,
		},
	}
	return b
}

// SecretToken the secret-token to use for this GitHub source
func (b *GitHubSourceBuilder) SecretToken(secretToken string) *GitHubSourceBuilder {
	b.ghSource.Spec.SecretToken = v1alpha1.SecretValueFromSource{
		SecretKeyRef: &corev1.SecretKeySelector{
			Key: secretToken,
		},
	}
	return b
}

// Sink or destination of the source
func (b *GitHubSourceBuilder) Sink(sink *duckv1.Destination) *GitHubSourceBuilder {
	b.ghSource.Spec.Sink = sink
	return b
}

// Build the GitHubSource object
func (b *GitHubSourceBuilder) Build() *v1alpha1.GitHubSource {
	return b.ghSource
}
