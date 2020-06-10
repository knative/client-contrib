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
	"github.com/maximilien/kn-source-github/pkg/types"

	sourcetypes "github.com/maximilien/kn-source-pkg/pkg/types"
)

// GHSourceFactory

func (f *ghSourceFactory) KnSourceParams() *sourcetypes.KnSourceParams {
	if f.ghSourceParams == nil {
		f.initGHSourceParams()
	}
	return f.ghSourceParams.KnSourceParams
}

func (f *ghSourceFactory) GHSourceParams() *types.GHSourceParams {
	if f.ghSourceParams == nil {
		f.initGHSourceParams()
	}
	return f.ghSourceParams
}

func (f *ghSourceFactory) GHSourceClient() types.GHSourceClient {
	return f.ghSourceClient
}

// CommandFactory

func (f *ghCommandFactory) KnSourceFactory() sourcetypes.KnSourceFactory {
	return f.ghSourceFactory
}

func (f *ghCommandFactory) GHSourceFactory() types.GHSourceFactory {
	return f.ghSourceFactory
}

func (f *ghCommandFactory) GHSourceParams() *types.GHSourceParams {
	return f.ghSourceFactory.GHSourceParams()
}

func (f *ghCommandFactory) KnSourceParams() *sourcetypes.KnSourceParams {
	return f.ghSourceFactory.KnSourceParams()
}

// FlagsFactory

func (f *ghFlagsFactory) KnSourceFactory() sourcetypes.KnSourceFactory {
	return f.ghSourceFactory
}

func (f *ghFlagsFactory) KnSourceParams() *sourcetypes.KnSourceParams {
	return f.ghSourceFactory.KnSourceParams()
}

func (f *ghFlagsFactory) GHSourceParams() *types.GHSourceParams {
	return f.ghSourceFactory.GHSourceParams()
}

func (f *ghFlagsFactory) GHSourceFactory() types.GHSourceFactory {
	return f.ghSourceFactory
}

// RunEFactory

func (f *ghRunEFactory) KnSourceParams() *sourcetypes.KnSourceParams {
	return f.GHSourceFactory().KnSourceParams()
}

func (f *ghRunEFactory) KnSourceClient(namespace string) sourcetypes.KnSourceClient {
	return f.GHSourceFactory().CreateGHSourceClient(namespace)
}

func (f *ghRunEFactory) GHSourceClient(namespace string) types.GHSourceClient {
	return f.GHSourceFactory().CreateGHSourceClient(namespace)
}

func (f *ghRunEFactory) KnSourceFactory() sourcetypes.KnSourceFactory {
	return f.GHSourceFactory()
}

func (f *ghRunEFactory) GHSourceFactory() types.GHSourceFactory {
	return f.ghSourceFactory
}
