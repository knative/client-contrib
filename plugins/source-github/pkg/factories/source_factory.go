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
	"github.com/maximilien/kn-source-github/pkg/client"
	"github.com/maximilien/kn-source-github/pkg/types"

	"knative.dev/client/pkg/kn/commands/flags"

	sourcefactories "github.com/maximilien/kn-source-pkg/pkg/factories"
	sourcetypes "github.com/maximilien/kn-source-pkg/pkg/types"
)

type ghSourceFactory struct {
	ghSourceParams *types.GHSourceParams
	ghSourceClient types.GHSourceClient

	knSourceFactory sourcetypes.KnSourceFactory
}

func NewGHSourceFactory() types.GHSourceFactory {
	return &ghSourceFactory{
		ghSourceParams:  nil,
		ghSourceClient:  nil,
		knSourceFactory: sourcefactories.NewDefaultKnSourceFactory(),
	}
}

func (f *ghSourceFactory) CreateKnSourceParams() *sourcetypes.KnSourceParams {
	if f.ghSourceParams == nil {
		f.initGHSourceParams()
	}
	return f.ghSourceParams.KnSourceParams
}

func (f *ghSourceFactory) CreateGHSourceParams() *types.GHSourceParams {
	if f.ghSourceParams == nil {
		f.initGHSourceParams()
	}
	return f.ghSourceParams
}

func (f *ghSourceFactory) CreateKnSourceClient(namespace string) sourcetypes.KnSourceClient {
	return f.CreateGHSourceClient(namespace)
}

func (f *ghSourceFactory) CreateGHSourceClient(namespace string) types.GHSourceClient {
	if f.ghSourceClient == nil {
		f.initGHSourceClient(namespace)
	}
	return f.ghSourceClient
}

// Private

func (f *ghSourceFactory) initGHSourceClient(namespace string) {
	if f.ghSourceClient == nil {
		f.ghSourceClient = client.NewGHSourceClient(f.GHSourceParams(), namespace)
	}
}

// Private

func (f *ghSourceFactory) initGHSourceParams() {
	f.ghSourceParams = &types.GHSourceParams{
		KnSourceParams: &sourcetypes.KnSourceParams{
			SinkFlag: flags.SinkFlags{},
		},
	}
	f.ghSourceParams.KnSourceParams.Initialize()
}
