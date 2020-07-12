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
	"github.com/maximilien/kn-source-pkg/pkg/types"
	"k8s.io/client-go/rest"
)

// KnSourceFactory

func (f *DefautKnSourceFactory) KnSourceParams() *types.KnSourceParams {
	if f.knSourceParams == nil {
		f.knSourceParams = f.CreateKnSourceParams()
	}

	return f.knSourceParams
}

// CommandFactory

func (f *DefautCommandFactory) KnSourceFactory() types.KnSourceFactory {
	return f.knSourceFactory
}

// FlagsFactory

func (f *DefautFlagsFactory) KnSourceFactory() types.KnSourceFactory {
	return f.knSourceFactory
}

// RunEFactory

func (f *DefautRunEFactory) KnSourceFactory() types.KnSourceFactory {
	return f.knSourceFactory
}

func (f *DefautRunEFactory) KnSourceClient(restConfig *rest.Config, namespace string) types.KnSourceClient {
	return f.knSourceFactory.CreateKnSourceClient(restConfig, namespace)
}
