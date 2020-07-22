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
	v1alpha1 "knative.dev/eventing-contrib/kafka/source/pkg/apis/sources/v1alpha1"

	sourcetypes "github.com/maximilien/kn-source-pkg/pkg/types"

	"github.com/spf13/cobra"
	"k8s.io/client-go/rest"
	"knative.dev/client/pkg/kn/commands"
	"knative.dev/client/pkg/printers"
	duckv1 "knative.dev/pkg/apis/duck/v1"
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

func NewFakeKafkaSourceRunEFactory(ns string) types.KafkaSourceRunEFactory {
	kafkaFactory := NewFakeKafkaSourceFactory(ns)
	return &kafkaSourceRunEFactory{
		kafkaSourceFactory: kafkaFactory,
		kafkaSourceClient:  kafkaFactory.KafkaSourceClient(),
	}
}

func (f *kafkaSourceRunEFactory) KafkaSourceClient(restConfig *rest.Config, namespace string) (types.KafkaSourceClient, error) {
	var err error
	f.kafkaSourceClient, err = f.KafkaSourceFactory().CreateKafkaSourceClient(restConfig, namespace)
	return f.kafkaSourceClient, err
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

		restConfig, err := f.KnSourceParams().KnParams.RestConfig()
		if err != nil {
			return err
		}

		f.kafkaSourceClient, err = f.KafkaSourceClient(restConfig, namespace)
		if err != nil {
			return err
		}

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
				"cannot create kafka '%s' in namespace '%s' "+
					"because: %s", name, f.kafkaSourceClient.Namespace(), err)
		}

		b := client.NewKafkaSourceBuilder(name).
			BootstrapServers(f.kafkaSourceFactory.KafkaSourceParams().BootstrapServers).
			Topics(f.kafkaSourceFactory.KafkaSourceParams().Topics).
			ConsumerGroup(f.kafkaSourceFactory.KafkaSourceParams().ConsumerGroup).
			Sink(objectRef)

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

		restConfig, err := f.KnSourceParams().KnParams.RestConfig()
		if err != nil {
			return err
		}

		f.kafkaSourceClient, err = f.KafkaSourceClient(restConfig, namespace)
		if err != nil {
			return err
		}

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
		fmt.Printf("Kafka source update is not supported because kafka source spec is immutable.\n")
		return nil
	}
}

func (f *kafkaSourceRunEFactory) DescribeRunE() sourcetypes.RunE {
	return func(cmd *cobra.Command, args []string) error {
		var err error
		namespace, err := f.KnSourceParams().GetNamespace(cmd)
		if err != nil {
			return err
		}

		restConfig, err := f.KnSourceParams().KnParams.RestConfig()
		if err != nil {
			return err
		}

		f.kafkaSourceClient, err = f.KafkaSourceClient(restConfig, namespace)
		if err != nil {
			return err
		}

		if len(args) != 1 {
			return errors.New("requires the name of the source to create as single argument")
		}
		name := args[0]

		kafkaSource, err := f.kafkaSourceClient.GetKafkaSource(name)

		if err != nil {
			return fmt.Errorf(
				"cannot describe KafkaSource '%s' in namespace '%s' "+
					"because: %s", name, f.kafkaSourceClient.Namespace(), err)
		}

		out := cmd.OutOrStdout()
		dw := printers.NewPrefixWriter(out)

		writeKafkaSource(dw, kafkaSource)
		dw.WriteLine()
		if err := dw.Flush(); err != nil {
			return err
		}

		writeSink(dw, kafkaSource.Spec.Sink)
		dw.WriteLine()
		if err := dw.Flush(); err != nil {
			return err
		}

		commands.WriteConditions(dw, kafkaSource.Status.Conditions, true)
		if err := dw.Flush(); err != nil {
			return err
		}
		return nil
	}
}

func writeSink(dw printers.PrefixWriter, sink *duckv1.Destination) {
	subWriter := dw.WriteAttribute("Sink", "")
	subWriter.WriteAttribute("Name", sink.Ref.Name)
	subWriter.WriteAttribute("Namespace", sink.Ref.Namespace)
	ref := sink.Ref
	if ref != nil {
		subWriter.WriteAttribute("Kind", fmt.Sprintf("%s (%s)", sink.Ref.Kind, sink.Ref.APIVersion))
	}
	uri := sink.URI
	if uri != nil {
		subWriter.WriteAttribute("URI", uri.String())
	}
}

func writeKafkaSource(dw printers.PrefixWriter, source *v1alpha1.KafkaSource) {
	commands.WriteMetadata(dw, &source.ObjectMeta, true)
	dw.WriteAttribute("BootstrapServers", source.Spec.BootstrapServers)
	dw.WriteAttribute("Topics", source.Spec.Topics)
	dw.WriteAttribute("ConsumerGroup", source.Spec.ConsumerGroup)
}
