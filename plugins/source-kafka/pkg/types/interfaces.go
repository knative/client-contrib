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

package types

import (
	v1alpha1 "knative.dev/eventing-contrib/kafka/source/pkg/apis/sources/v1alpha1"

	sourcetypes "github.com/maximilien/kn-source-pkg/pkg/types"
)

type KafkaSourceClient interface {
	sourcetypes.KnSourceClient
	CreateKafkaSource(kafkaSource *v1alpha1.KafkaSource) error
}

type KafkaSourceFactory interface {
	sourcetypes.KnSourceFactory

	KafkaSourceParams() *KafkaSourceParams
	KafkaSourceClient() KafkaSourceClient

	CreateKafkaSourceClient(namespace string) KafkaSourceClient
	CreateKafkaSourceParams() *KafkaSourceParams
}

type KafkaSourceCommandFactory interface {
	sourcetypes.CommandFactory

	KafkaSourceFactory() KafkaSourceFactory
}

type KafkaSourceFlagsFactory interface {
	sourcetypes.FlagsFactory

	KafkaSourceFactory() KafkaSourceFactory
}

type KafkaSourceRunEFactory interface {
	sourcetypes.RunEFactory

	KafkaSourceFactory() KafkaSourceFactory
	KafkaSourceClient(namespace string) KafkaSourceClient
}
