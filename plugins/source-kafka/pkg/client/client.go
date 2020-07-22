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
	sourceclient "github.com/maximilien/kn-source-pkg/pkg/client"
	sourcetypes "github.com/maximilien/kn-source-pkg/pkg/types"
	"knative.dev/client-contrib/plugins/source-kafka/pkg/types"
	knerrors "knative.dev/client/pkg/errors"
	v1alpha1 "knative.dev/eventing-contrib/kafka/source/pkg/apis/sources/v1alpha1"
	clientv1alpha1 "knative.dev/eventing-contrib/kafka/source/pkg/client/clientset/versioned/typed/sources/v1alpha1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	duckv1 "knative.dev/pkg/apis/duck/v1"
)

type kafkaSourceClient struct {
	namespace         string
	kafkaSourceParams *types.KafkaSourceParams
	knSourceClient    sourcetypes.KnSourceClient
	client            clientv1alpha1.SourcesV1alpha1Interface
}

// NewKafkaSourceClient is to create a KafkaSourceClient
func NewKafkaSourceClient(kafkaParams *types.KafkaSourceParams, restConfig *rest.Config, ns string) (types.KafkaSourceClient, error) {
	kafkaClient, err := kafkaParams.NewSourcesClient()
	if err != nil {
		return nil, knerrors.GetError(err)
	}
	return &kafkaSourceClient{
		kafkaSourceParams: kafkaParams,
		namespace:         ns,
		knSourceClient:    sourceclient.NewKnSourceClient(kafkaParams.KnSourceParams, restConfig, ns),
		client:            kafkaClient,
	}, nil
}

// RestConfig the REST cconfig
func (c *kafkaSourceClient) RestConfig() *rest.Config {
	return c.knSourceClient.RestConfig()
}

// KnSourceParams for common Kn source parameters
func (c *kafkaSourceClient) KnSourceParams() *sourcetypes.KnSourceParams {
	return c.kafkaSourceParams.KnSourceParams
}

// KafkaSourceParams for kafka source specific parameters
func (c *kafkaSourceClient) KafkaSourceParams() *types.KafkaSourceParams {
	return c.kafkaSourceParams
}

//CreateKafkaSource is used to create an instance of KafkaSource
func (c *kafkaSourceClient) CreateKafkaSource(kafkaSource *v1alpha1.KafkaSource) error {
	_, err := c.client.KafkaSources(c.namespace).Create(kafkaSource)
	if err != nil {
		return knerrors.GetError(err)
	}

	return nil
}

//DeleteKafkaSource is used to create an instance of KafkaSource
func (c *kafkaSourceClient) DeleteKafkaSource(name string) error {
	err := c.client.KafkaSources(c.namespace).Delete(name, &metav1.DeleteOptions{})
	if err != nil {
		return knerrors.GetError(err)
	}

	return nil
}

//GetKafkaSource is used to create an instance of KafkaSource
func (c *kafkaSourceClient) GetKafkaSource(name string) (*v1alpha1.KafkaSource, error) {
	kafkaSource, err := c.client.KafkaSources(c.namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		return nil, knerrors.GetError(err)
	}

	return kafkaSource, nil
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

// NewKafkaSourceBuilderFromExisting for building the object from existing KafkaSource object
func NewKafkaSourceBuilderFromExisting(kSource *v1alpha1.KafkaSource) *KafkaSourceBuilder {
	return &KafkaSourceBuilder{kafkaSource: kSource.DeepCopy()}
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
func (b *KafkaSourceBuilder) Sink(sink *duckv1.Destination) *KafkaSourceBuilder {
	b.kafkaSource.Spec.Sink = sink
	return b
}

// Build the KafkaSource object
func (b *KafkaSourceBuilder) Build() *v1alpha1.KafkaSource {
	return b.kafkaSource
}
