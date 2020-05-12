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
	"github.com/maximilien/kn-source-pkg/pkg/types"
)

type knSourceClient struct {
	knSourceParams *types.KnSourceParams
	namespace      string
}

// NewKnSourceClient creates a new KnSourceClient with parameters and namespace
func NewKnSourceClient(knSourceParams *types.KnSourceParams, namespace string) types.KnSourceClient {
	return &knSourceClient{
		knSourceParams: knSourceParams,
		namespace:      namespace,
	}
}

// KnSourceParams returns the client's KnSourceParams
func (client *knSourceClient) KnSourceParams() *types.KnSourceParams {
	return client.knSourceParams
}

// Namespace returns the client's namespace
func (client *knSourceClient) Namespace() string {
	return client.namespace
}
