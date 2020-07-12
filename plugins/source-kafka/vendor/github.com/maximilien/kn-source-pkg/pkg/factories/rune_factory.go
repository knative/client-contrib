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
	"fmt"

	"github.com/maximilien/kn-source-pkg/pkg/types"

	"github.com/spf13/cobra"
)

type DefautRunEFactory struct {
	knSourceFactory types.KnSourceFactory
}

func NewDefaultRunEFactory(knSourceFactory types.KnSourceFactory) types.RunEFactory {
	return &DefautRunEFactory{
		knSourceFactory: knSourceFactory,
	}
}

func (f *DefautRunEFactory) CreateRunE() types.RunE {
	return func(cmd *cobra.Command, args []string) error {
		namespace, err := f.KnSourceFactory().KnSourceParams().GetNamespace(cmd)
		if err != nil {
			return err
		}

		restConfig, err := f.KnSourceFactory().KnSourceParams().KnParams.RestConfig()
		if err != nil {
			return err
		}

		knSourceClient := f.KnSourceClient(restConfig, namespace)

		fmt.Printf("%s RunE called: args: %#v, client: %#v, sink: %s\n", cmd.Name(), args, knSourceClient, knSourceClient.KnSourceParams().SinkFlag)

		return nil
	}
}

func (f *DefautRunEFactory) DeleteRunE() types.RunE {
	return func(cmd *cobra.Command, args []string) error {
		namespace, err := f.KnSourceFactory().KnSourceParams().GetNamespace(cmd)
		if err != nil {
			return err
		}

		restConfig, err := f.KnSourceFactory().KnSourceParams().KnParams.RestConfig()
		if err != nil {
			return err
		}

		knSourceClient := f.KnSourceClient(restConfig, namespace)

		fmt.Printf("%s RunE called: args: %#v, client: %#v, sink: %s\n", cmd.Name(), args, knSourceClient, knSourceClient.KnSourceParams().SinkFlag)

		return nil
	}
}

func (f *DefautRunEFactory) UpdateRunE() types.RunE {
	return func(cmd *cobra.Command, args []string) error {
		namespace, err := f.KnSourceFactory().KnSourceParams().GetNamespace(cmd)
		if err != nil {
			return err
		}

		restConfig, err := f.KnSourceFactory().KnSourceParams().KnParams.RestConfig()
		if err != nil {
			return err
		}

		knSourceClient := f.KnSourceClient(restConfig, namespace)

		fmt.Printf("%s RunE called: args: %#v, client: %#v, sink: %s\n", cmd.Name(), args, knSourceClient, knSourceClient.KnSourceParams().SinkFlag)

		return nil
	}
}

func (f *DefautRunEFactory) DescribeRunE() types.RunE {
	return func(cmd *cobra.Command, args []string) error {
		namespace, err := f.KnSourceFactory().KnSourceParams().GetNamespace(cmd)
		if err != nil {
			return err
		}

		restConfig, err := f.KnSourceFactory().KnSourceParams().KnParams.RestConfig()
		if err != nil {
			return err
		}

		knSourceClient := f.KnSourceClient(restConfig, namespace)

		fmt.Printf("%s RunE called: args: %#v, client: %#v, sink: %s\n", cmd.Name(), args, knSourceClient, knSourceClient.KnSourceParams().SinkFlag)

		return nil
	}
}
