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
	sourcetypes "github.com/maximilien/kn-source-pkg/pkg/types"
	"k8s.io/client-go/rest"
	"knative.dev/client-contrib/plugins/source-kafka/pkg/types"
)

// source_factory

func (f *kafkaClientFactory) KafkaSourceClient() types.KafkaSourceClient {
	return f.kafkaSourceClient
}

func (f *kafkaClientFactory) KnSourceParams() *sourcetypes.KnSourceParams {
	return f.CreateKnSourceParams()
}

func (f *kafkaClientFactory) CreateKnSourceParams() *sourcetypes.KnSourceParams {
	if f.kafkaSourceParams == nil {
		f.initKafkaSourceParams()
	}
	return f.kafkaSourceParams.KnSourceParams
}

func (f *kafkaClientFactory) CreateKnSourceClient(restConfig *rest.Config, namespace string) sourcetypes.KnSourceClient {
	f.CreateKafkaSourceClient(restConfig, namespace)
	return f.KafkaSourceClient()
}

// rune_factory

func (f *kafkaSourceRunEFactory) KnSourceParams() *sourcetypes.KnSourceParams {
	return f.KafkaSourceFactory().KnSourceParams()
}

func (f *kafkaSourceRunEFactory) KnSourceClient(restConfig *rest.Config, namespace string) sourcetypes.KnSourceClient {
	return f.KafkaSourceFactory().CreateKnSourceClient(restConfig, namespace)
}

func (f *kafkaSourceRunEFactory) KnSourceFactory() sourcetypes.KnSourceFactory {
	return f.kafkaSourceFactory
}

// flags_factory

func (f *kafkaSourceFlagsFactory) KnSourceFactory() sourcetypes.KnSourceFactory {
	return f.kafkaSourceFactory
}

func (f *kafkaSourceFlagsFactory) KnSourceParams() *sourcetypes.KnSourceParams {
	return f.kafkaSourceFactory.KnSourceParams()
}

// command_factory
func (f *kafkaSourceCommandFactory) KnSourceFactory() sourcetypes.KnSourceFactory {
	return f.kafkaSourceFactory
}

func (f *kafkaSourceCommandFactory) KnSourceParams() *sourcetypes.KnSourceParams {
	return f.kafkaSourceFactory.KnSourceParams()
}
