// Copyright © 2018 The Knative Authors
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
//
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

package types

import (
	sourcetypes "github.com/maximilien/kn-source-pkg/pkg/types"

	v1alpha1 "knative.dev/eventing-contrib/github/pkg/apis/sources/v1alpha1"
)

// GHSourceClient the GitHub source client interface
//counterfeiter:generate . GHSourceClient
type GHSourceClient interface {
	sourcetypes.KnSourceClient
	GHSourceParams() *GHSourceParams

	GetGHSource(name string) (*v1alpha1.GitHubSource, error)
	CreateGHSource(ghSource *v1alpha1.GitHubSource) (*v1alpha1.GitHubSource, error)
	UpdateGHSource(ghSource *v1alpha1.GitHubSource) (*v1alpha1.GitHubSource, error)
	DeleteGHSource(name string) error
}

// GHSourceFactory the GitHub source factory interface
//counterfeiter:generate . GHSourceFactory
type GHSourceFactory interface {
	sourcetypes.KnSourceFactory

	GHSourceParams() *GHSourceParams
	GHSourceClient() GHSourceClient

	CreateGHSourceParams() *GHSourceParams
	CreateGHSourceClient(namespace string) GHSourceClient
}

// GHCommandFactory the GitHub source command factory interface
//counterfeiter:generate . GHCommandFactory
type GHCommandFactory interface {
	sourcetypes.CommandFactory

	GHSourceFactory() GHSourceFactory
}

// GHFlagsFactory the GitHub source flags factory interface
//counterfeiter:generate . GHFlagsFactory
type GHFlagsFactory interface {
	sourcetypes.FlagsFactory

	GHSourceFactory() GHSourceFactory
}

// GHRunEFactory the GitHub source RunE factory interface
//counterfeiter:generate . GHRunEFactory
type GHRunEFactory interface {
	sourcetypes.RunEFactory

	GHSourceFactory() GHSourceFactory
	GHSourceClient(namespace string) GHSourceClient
}
