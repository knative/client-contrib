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
	"errors"
	"fmt"

	"knative.dev/client-contrib/plugins/source-kafka/pkg/client"
	"knative.dev/client-contrib/plugins/source-kafka/pkg/types"

	sourcetypes "github.com/maximilien/kn-source-pkg/pkg/types"
	"github.com/maximilien/kn-source-pkg/pkg/util"

	"github.com/spf13/cobra"
)

type kafkaSourceRunEFactory struct {
	kafkaSourceClient  types.KafkaSourceClient
	kafkaSourceFactory types.KafkaSourceFactory
}

func NewKafkaSourceRunEFactory(kafkaFactory types.KafkaSourceFactory) types.KafkaSourceRunEFactory {
	return &kafkaSourceRunEFactory{
		kafkaSourceFactory: kafkaFactory,
		kafkaSourceClient:  kafkaFactory.KafkaSourceClient(),
	}
}

func (f *kafkaSourceRunEFactory) KnSourceParams() *sourcetypes.KnSourceParams {
	return f.KafkaSourceFactory().KnSourceParams()
}

func (f *kafkaSourceRunEFactory) KnSourceClient(namespace string) sourcetypes.KnSourceClient {
	return f.KafkaSourceFactory().KafkaSourceClient()
}

func (f *kafkaSourceRunEFactory) KafkaSourceClient(namespace string) types.KafkaSourceClient {
	f.kafkaSourceClient = f.kafkaSourceFactory.CreateKafkaSourceClient(namespace)
	return f.kafkaSourceClient
}

func (f *kafkaSourceRunEFactory) KnSourceFactory() sourcetypes.KnSourceFactory {
	return f.kafkaSourceFactory
}

func (f *kafkaSourceRunEFactory) KafkaSourceFactory() types.KafkaSourceFactory {
	return f.kafkaSourceFactory
}

func (f *kafkaSourceRunEFactory) CreateRunE() sourcetypes.RunE {
	return func(cmd *cobra.Command, args []string) error {
		var err error
		namespace, err := f.KnSourceParams().GetNamespace(cmd)
		if err != nil {
			return err
		}
		f.kafkaSourceClient = f.KafkaSourceClient(namespace)

		if len(args) != 1 {
			return errors.New("requires the name of the source to create as single argument")
		}
		name := args[0]

		dynamicClient, err := f.KnSourceParams().KnParams.NewDynamicClient(f.kafkaSourceClient.Namespace())
		if err != nil {
			return err
		}
		objectRef, err := f.KnSourceParams().SinkFlag.ResolveSink(dynamicClient, f.kafkaSourceClient.Namespace())
		if err != nil {
			return fmt.Errorf(
				"cannot create Kafka '%s' in namespace '%s' "+
					"because: %s", name, f.kafkaSourceClient.Namespace(), err)
		}

		b := client.NewKafkaSourceBuilder(name).
			BootstrapServers(f.kafkaSourceFactory.KafkaSourceParams().BootstrapServers).
			Topics(f.kafkaSourceFactory.KafkaSourceParams().Topics).
			ConsumerGroup(f.kafkaSourceFactory.KafkaSourceParams().ConsumerGroup).
			Sink(util.SinkToDuckV1Beta1(objectRef))

		err = f.kafkaSourceClient.CreateKafkaSource(b.Build())

		if err != nil {
			return fmt.Errorf(
				"cannot create KafkaSource '%s' in namespace '%s' "+
					"because: %s", name, f.kafkaSourceClient.Namespace(), err)
		}

		if err == nil {
			fmt.Fprintf(cmd.OutOrStdout(), "Kafka source '%s' created in namespace '%s'.\n", args[0], f.kafkaSourceClient.Namespace())
		}

		return err
	}
}

func (f *kafkaSourceRunEFactory) DeleteRunE() sourcetypes.RunE {
	return func(cmd *cobra.Command, args []string) error {
		var err error
		namespace, err := f.KnSourceParams().GetNamespace(cmd)
		if err != nil {
			return err
		}
		f.kafkaSourceClient = f.KafkaSourceClient(namespace)

		if len(args) != 1 {
			return errors.New("requires the name of the source to create as single argument")
		}
		name := args[0]

		err = f.kafkaSourceClient.DeleteKafkaSource(name)

		if err != nil {
			return fmt.Errorf(
				"cannot delete KafkaSource '%s' in namespace '%s' "+
					"because: %s", name, f.kafkaSourceClient.Namespace(), err)
		}

		if err == nil {
			fmt.Fprintf(cmd.OutOrStdout(), "Kafka source '%s' deleted in namespace '%s'.\n", args[0], f.kafkaSourceClient.Namespace())
		}

		return err
	}
}

func (f *kafkaSourceRunEFactory) UpdateRunE() sourcetypes.RunE {
	return func(cmd *cobra.Command, args []string) error {
		fmt.Printf("%s RunE function called for Kafka source: args: %#v, client: %#v\n", cmd.Name(), args, f.kafkaSourceClient)
		return nil
	}
}

func (f *kafkaSourceRunEFactory) DescribeRunE() sourcetypes.RunE {
	return func(cmd *cobra.Command, args []string) error {
		fmt.Printf("%s RunE function called for Kafka source: args: %#v, client: %#v\n", cmd.Name(), args, f.kafkaSourceClient)
		return nil
	}
}
