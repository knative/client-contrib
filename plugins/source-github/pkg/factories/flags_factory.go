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

package factories

import (
	sourcefactories "github.com/maximilien/kn-source-pkg/pkg/factories"
	sourcetypes "github.com/maximilien/kn-source-pkg/pkg/types"

	"github.com/maximilien/kn-source-github/pkg/types"

	"github.com/spf13/pflag"
)

type ghFlagsFactory struct {
	defaultFlagsFactory sourcetypes.FlagsFactory
	ghSourceFactory     types.GHSourceFactory
}

func NewGHFlagsFactory(ghSourceFactory types.GHSourceFactory) types.GHFlagsFactory {
	return &ghFlagsFactory{
		defaultFlagsFactory: sourcefactories.NewDefaultFlagsFactory(ghSourceFactory),
		ghSourceFactory:     ghSourceFactory,
	}
}

func (f *ghFlagsFactory) CreateFlags() *pflag.FlagSet {
	flagSet := f.defaultFlagsFactory.CreateFlags()
	flagSet.StringVar(&f.GHSourceParams().Org, "org", "", "The GitHub organization or username")
	flagSet.StringVar(&f.GHSourceParams().Repo, "repo", "", "Repository name to consume messages from")
	flagSet.StringVar(&f.GHSourceParams().APIURL, "api-url", "https://api.github.com", "The GitHub API URL to use")
	flagSet.StringVar(&f.GHSourceParams().SecretToken, "secret-token", "", "The GitHub secret-token to use")
	flagSet.StringVar(&f.GHSourceParams().AccessToken, "access-token", "", "The GitHub access-token to use")
	return flagSet
}

func (f *ghFlagsFactory) DeleteFlags() *pflag.FlagSet {
	flagSet := f.defaultFlagsFactory.DeleteFlags()
	return flagSet
}

func (f *ghFlagsFactory) UpdateFlags() *pflag.FlagSet {
	flagSet := f.defaultFlagsFactory.UpdateFlags()
	flagSet.StringVar(&f.GHSourceParams().Org, "org", "", "The GitHub organization or username")
	flagSet.StringVar(&f.GHSourceParams().Repo, "repo", "", "Repository name to consume messages from")
	flagSet.StringVar(&f.GHSourceParams().APIURL, "api-url", "https://api.github.com", "The GitHub API URL to use")
	flagSet.StringVar(&f.GHSourceParams().SecretToken, "secret-token", "", "The GitHub secret-token to use")
	flagSet.StringVar(&f.GHSourceParams().AccessToken, "access-token", "", "The GitHub access-token to use")
	return flagSet
}

func (f *ghFlagsFactory) DescribeFlags() *pflag.FlagSet {
	flagSet := f.defaultFlagsFactory.DescribeFlags()
	return flagSet
}
