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
	sourcetypes "github.com/maximilien/kn-source-pkg/pkg/types"
	"knative.dev/client-contrib/plugins/source-kafka/pkg/types"
	knerrors "knative.dev/client/pkg/errors"
	v1alpha1 "knative.dev/eventing-contrib/kafka/source/pkg/apis/sources/v1alpha1"
	clientv1alpha1 "knative.dev/eventing-contrib/kafka/source/pkg/client/clientset/versioned/typed/sources/v1alpha1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	duckv1beta1 "knative.dev/pkg/apis/duck/v1beta1"
)

type kafkaSourceClient struct {
	kafkaSourceParams *types.KafkaSourceParams
	client            clientv1alpha1.SourcesV1alpha1Interface
	namespace         string
}

func NewKafkaSourceClient(kafkaParams *types.KafkaSourceParams, c clientv1alpha1.SourcesV1alpha1Interface, ns string) types.KafkaSourceClient {
	if c == nil {
		c, _ = kafkaParams.NewSourcesClient()
	}
	return &kafkaSourceClient{
		kafkaSourceParams: kafkaParams,
		client:            c,
		namespace:         ns,
	}
}

func (client *kafkaSourceClient) KnSourceParams() *sourcetypes.KnSourceParams {
	return client.kafkaSourceParams.KnSourceParams
}

func (client *kafkaSourceClient) KafkaSourceParams() *types.KafkaSourceParams {
	return client.kafkaSourceParams
}

//CreateKafkaSource is used to create an instance of ApiServerSource
func (c *kafkaSourceClient) CreateKafkaSource(kafkaSource *v1alpha1.KafkaSource) error {
	_, err := c.client.KafkaSources(c.namespace).Create(kafkaSource)
	if err != nil {
		return knerrors.GetError(err)
	}

	return nil
}

// Return the client's namespace
func (c *kafkaSourceClient) Namespace() string {
	return c.namespace
}

// KafkaSourceBuilder is for building the source
type KafkaSourceBuilder struct {
	kafkaSource *v1alpha1.KafkaSource
}

// NewKafkaSourceBuilder for building ApiServer source object
func NewKafkaSourceBuilder(name string) *KafkaSourceBuilder {
	return &KafkaSourceBuilder{kafkaSource: &v1alpha1.KafkaSource{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}}
}

// BootstrapServers to set the value of BootstrapServers
func (b *KafkaSourceBuilder) BootstrapServers(server string) *KafkaSourceBuilder {
	b.kafkaSource.Spec.BootstrapServers = server
	return b
}

// Topics to set the value of Topics
func (b *KafkaSourceBuilder) Topics(topics string) *KafkaSourceBuilder {
	b.kafkaSource.Spec.Topics = topics
	return b
}

// ConsumerGroup to set the value of ConsumerGroup
func (b *KafkaSourceBuilder) ConsumerGroup(consumerGroup string) *KafkaSourceBuilder {
	b.kafkaSource.Spec.ConsumerGroup = consumerGroup
	return b
}

// Sink or destination of the source
func (b *KafkaSourceBuilder) Sink(sink *duckv1beta1.Destination) *KafkaSourceBuilder {
	b.kafkaSource.Spec.Sink = sink
	return b
}

// Build the KafkaSource object
func (b *KafkaSourceBuilder) Build() *v1alpha1.KafkaSource {
	return b.kafkaSource
}
