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
	"knative.dev/client-contrib/plugins/source-kafka/pkg/types"

	sourcefactories "github.com/maximilien/kn-source-pkg/pkg/factories"
	sourcetypes "github.com/maximilien/kn-source-pkg/pkg/types"

	"github.com/spf13/cobra"
)

type kafkaSourceCommandFactory struct {
	kafkaSourceFactory    types.KafkaSourceFactory
	defaultCommandFactory sourcetypes.CommandFactory
}

func NewKafkaSourceCommandFactory(kafkaFactory types.KafkaSourceFactory) types.KafkaSourceCommandFactory {
	return &kafkaSourceCommandFactory{
		kafkaSourceFactory:    kafkaFactory,
		defaultCommandFactory: sourcefactories.NewDefaultCommandFactory(kafkaFactory),
	}
}

func (f *kafkaSourceCommandFactory) KafkaSourceFactory() types.KafkaSourceFactory {
	return f.kafkaSourceFactory
}

func (f *kafkaSourceCommandFactory) KafkaSourceParams() *types.KafkaSourceParams {
	return f.kafkaSourceFactory.KafkaSourceParams()
}

func (f *kafkaSourceCommandFactory) SourceCommand() *cobra.Command {
	sourceCmd := f.defaultCommandFactory.SourceCommand()
	sourceCmd.Use = "kafka"
	sourceCmd.Short = "Knative eventing kafka source plugin"
	sourceCmd.Long = "Manage Knative kafka eventing sources"
	return sourceCmd
}

func (f *kafkaSourceCommandFactory) CreateCommand() *cobra.Command {
	createCmd := f.defaultCommandFactory.CreateCommand()
	createCmd.Use = "create NAME --servers SERVERS --topics TOPICS --consumergroup GROUP --sink SINK"
	createCmd.Short = "Create a kafka source"
	createCmd.Example = `# Create a new kafka source 'mykafkasrc' which subscribes a kafka server 'my-cluster-kafka-bootstrap.kafka.svc:9092' at topic 'test-topic' using the consumer group ID 'test-consumer-group' and sends the events to service 'event-display'
kn source kafka create mykafkasrc --servers my-cluster-kafka-bootstrap.kafka.svc:9092 --topics test-topic --consumergroup test-consumer-group --sink svc:event-display`
	return createCmd
}

func (f *kafkaSourceCommandFactory) DeleteCommand() *cobra.Command {
	deleteCmd := f.defaultCommandFactory.DeleteCommand()
	deleteCmd.Use = "delete NAME"
	deleteCmd.Short = "Delete a kafka source"
	deleteCmd.Example = `# Delete a kafka source with name 'mykafkasrc'
kn source kafka delete mykafkasrc`
	return deleteCmd
}

func (f *kafkaSourceCommandFactory) UpdateCommand() *cobra.Command {
	return nil
}

func (f *kafkaSourceCommandFactory) DescribeCommand() *cobra.Command {
	describeCmd := f.defaultCommandFactory.DescribeCommand()
	describeCmd.Use = "describe NAME"
	describeCmd.Short = "Describe a kafka source"
	describeCmd.Example = `# Describe a kafka source with NAME
kn source kafka describe kafka-name`
	return describeCmd
}
